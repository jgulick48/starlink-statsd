package starlink

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jgulick48/starlink-statsd/internal/metrics"
	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
	"google.golang.org/grpc"
)

const metricsNamespace = "starlink"

type Client interface {
	StartStatsLoop(ch chan bool, tick time.Duration)
}

type client struct {
	addr         string
	deploymentId string
}

func NewClient(address string, deploymentId string) Client {
	client := &client{
		addr:         address,
		deploymentId: deploymentId,
	}
	return client
}

func (c *client) StartStatsLoop(ch chan bool, tick time.Duration) {
	timer := time.NewTicker(tick)
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	if err != nil {
	}
	slc := pb.NewDeviceClient(conn)
	for {
		select {
		case <-ch:
			return
		case <-timer.C:
			in := new(pb.Request)
			in.Request = &pb.Request_GetStatus{}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// r is the gprc response
			r, err := slc.Handle(ctx, in)
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			status := r.GetDishGetStatus()
			c.emitStats(status)
		}
	}
}

func (c *client) emitStats(status *pb.DishGetStatusResponse) {
	if !metrics.StatsEnabled || status == nil {
		return
	}
	standardTags := []string{
		fmt.Sprintf("%s:%v", "deploymentId", c.deploymentId),
	}
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, ""), []string{}, 0)
	if status.DeviceInfo != nil {
		tags := []string{
			fmt.Sprintf("%s:%s", "id", status.DeviceInfo.Id),
			fmt.Sprintf("%s:%s", "hardwareVersion", status.DeviceInfo.HardwareVersion),
			fmt.Sprintf("%s:%s", "softwareVersion", status.DeviceInfo.SoftwareVersion),
			fmt.Sprintf("%s:%s", "manufacturedVersion", status.DeviceInfo.ManufacturedVersion),
			fmt.Sprintf("%s:%s", "countryCode", status.DeviceInfo.CountryCode),
		}
		tags = append(tags, standardTags...)
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "deviceinfo"), tags, float64(0))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "generationnumber"), standardTags, float64(status.DeviceInfo.GenerationNumber))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "utcoffsets"), standardTags, float64(status.DeviceInfo.UtcOffsetS))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "bootcount"), standardTags, float64(status.DeviceInfo.Bootcount))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "antirollbackversion"), standardTags, float64(status.DeviceInfo.AntiRollbackVersion))
	}
	if status.DeviceState != nil {
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "uptimes"), standardTags, float64(status.DeviceState.UptimeS))
	}
	if status.Alerts != nil {
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:MotorsStuck"}, standardTags...), boolToFloat64(status.Alerts.MotorsStuck))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:ThermalThrottle"}, standardTags...), boolToFloat64(status.Alerts.ThermalThrottle))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:ThermalShutdown"}, standardTags...), boolToFloat64(status.Alerts.ThermalShutdown))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:MastNotNearVertical"}, standardTags...), boolToFloat64(status.Alerts.MastNotNearVertical))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:UnexpectedLocation"}, standardTags...), boolToFloat64(status.Alerts.UnexpectedLocation))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:SlowEthernetSpeeds"}, standardTags...), boolToFloat64(status.Alerts.SlowEthernetSpeeds))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:Roaming"}, standardTags...), boolToFloat64(status.Alerts.Roaming))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:InstallPending"}, standardTags...), boolToFloat64(status.Alerts.InstallPending))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:IsHeating"}, standardTags...), boolToFloat64(status.Alerts.IsHeating))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "alerts"), append([]string{"alertType:PowerSupplyThermalThrottle"}, standardTags...), boolToFloat64(status.Alerts.PowerSupplyThermalThrottle))
	}
	if status.GpsStats != nil {
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "gpsStats", "gpsvalid"), standardTags, boolToFloat64(status.GpsStats.GpsValid))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "gpsStats", "gpssats"), standardTags, float64(status.GpsStats.GpsSats))
	}
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "secondstofirstnonemptyslot"), standardTags, float64(status.SecondsToFirstNonemptySlot))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "poppingdroprate"), standardTags, float64(status.PopPingDropRate))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "downlinkthroughputbps"), standardTags, float64(status.DownlinkThroughputBps))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "uplinkthroughputbps"), standardTags, float64(status.UplinkThroughputBps))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "poppinglatencyms"), standardTags, float64(status.PopPingLatencyMs))
	if status.ObstructionStats != nil {
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "currentlyobstructed"), standardTags, boolToFloat64(status.ObstructionStats.CurrentlyObstructed))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "fractionobstructed"), standardTags, float64(status.ObstructionStats.FractionObstructed))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "valids"), standardTags, float64(status.ObstructionStats.ValidS))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "avgprolongedobstructiondurrations"), standardTags, float64(status.ObstructionStats.AvgProlongedObstructionDurationS))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "avgprolongedobstructionintervals"), standardTags, float64(status.ObstructionStats.AvgProlongedObstructionIntervalS))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "obstructionstats", "avgprolongedobstructionvalid"), standardTags, boolToFloat64(status.ObstructionStats.AvgProlongedObstructionValid))
	}
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "stowrequested"), standardTags, boolToFloat64(status.StowRequested))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "boresightazimuthdeg"), standardTags, float64(status.BoresightAzimuthDeg))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "boresightelevationdeg"), standardTags, float64(status.BoresightElevationDeg))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "ethspeedmbps"), standardTags, float64(status.EthSpeedMbps))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "mobilityclass"), append([]string{status.MobilityClass.String()}, standardTags...), float64(status.MobilityClass))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "issnrabovenoisefloor"), standardTags, boolToFloat64(status.IsSnrAboveNoiseFloor))
	if status.ReadyStates != nil {
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "cady"), standardTags, boolToFloat64(status.ReadyStates.Cady))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "scp"), standardTags, boolToFloat64(status.ReadyStates.Scp))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "l1l2"), standardTags, boolToFloat64(status.ReadyStates.L1L2))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "xphy"), standardTags, boolToFloat64(status.ReadyStates.Xphy))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "aap"), standardTags, boolToFloat64(status.ReadyStates.Aap))
		metrics.SendGaugeMetric(fmt.Sprintf("%s.%s.%s", metricsNamespace, "readystates", "rf"), standardTags, boolToFloat64(status.ReadyStates.Rf))
	}
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "issnrpersistentlylow"), standardTags, boolToFloat64(status.IsSnrPersistentlyLow))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "classofservice"), append([]string{status.ClassOfService.String()}, standardTags...), float64(status.ClassOfService))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "softwareupdatestate"), append([]string{status.SoftwareUpdateState.String()}, standardTags...), float64(status.SoftwareUpdateState))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "hasactuators"), append([]string{status.SoftwareUpdateState.String()}, standardTags...), float64(status.SoftwareUpdateState))
	metrics.SendGaugeMetric(fmt.Sprintf("%s.%s", metricsNamespace, "disablementcode"), append([]string{status.DisablementCode.String()}, standardTags...), float64(status.DisablementCode))
}

func boolToFloat64(val bool) float64 {
	if val {
		return float64(1)
	}
	return float64(0)
}
