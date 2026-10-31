package collector

import (
	"github.com/jgulick48/starlink-statsd/internal/metrics"
	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

type enumLike interface{ String() string }

func enumTag(key string, e enumLike) string { return metrics.Tag(key, e.String()) }

// EmitStatus sends all metrics derived from a get_status response.
func EmitStatus(m *metrics.Sink, st *pb.DishGetStatusResponse) {
	if st == nil {
		return
	}
	g := m.Gauge
	b := metrics.Bool

	if d := st.DeviceInfo; d != nil {
		g("deviceinfo", 0,
			metrics.Tag("id", d.Id),
			metrics.Tag("hardwareVersion", d.HardwareVersion),
			metrics.Tag("softwareVersion", d.SoftwareVersion),
			metrics.Tag("manufacturedVersion", d.ManufacturedVersion),
			metrics.Tag("countryCode", d.CountryCode))
		g("generationnumber", float64(d.GenerationNumber))
		g("utcoffsets", float64(d.UtcOffsetS))
		g("bootcount", float64(d.Bootcount))
		g("antirollbackversion", float64(d.AntiRollbackVersion))
	}
	if d := st.DeviceState; d != nil {
		g("uptimes", float64(d.UptimeS))
	}
	if a := st.Alerts; a != nil {
		for _, al := range []struct {
			name string
			v    bool
		}{
			{"MotorsStuck", a.MotorsStuck},
			{"ThermalThrottle", a.ThermalThrottle},
			{"ThermalShutdown", a.ThermalShutdown},
			{"MastNotNearVertical", a.MastNotNearVertical},
			{"UnexpectedLocation", a.UnexpectedLocation},
			{"SlowEthernetSpeeds", a.SlowEthernetSpeeds},
			{"InstallPending", a.InstallPending},
			{"IsHeating", a.IsHeating},
			{"PowerSupplyThermalThrottle", a.PowerSupplyThermalThrottle},
			{"DishWaterDetected", a.DishWaterDetected},
			{"RouterWaterDetected", a.RouterWaterDetected},
			{"LowerSignalThanPredicted", a.LowerSignalThanPredicted},
			{"NoEthernetLink", a.NoEthernetLink},
			{"ObstructionMapReset", a.ObstructionMapReset},
			{"SlowEthernetSpeeds100", a.SlowEthernetSpeeds_100},
			{"UpsuRouterPortSlow", a.UpsuRouterPortSlow},
			{"LowMotorCurrent", a.LowMotorCurrent},
			{"DbfTelemStale", a.DbfTelemStale},
			{"IsPowerSaveIdle", a.IsPowerSaveIdle},
		} {
			g("alerts", b(al.v), metrics.Tag("alertType", al.name))
		}
	}
	if s := st.GpsStats; s != nil {
		g("gpsStats.gpsvalid", b(s.GpsValid))
		g("gpsStats.gpssats", float64(s.GpsSats))
	}
	g("secondstofirstnonemptyslot", float64(st.SecondsToFirstNonemptySlot))
	g("poppingdroprate", float64(st.PopPingDropRate))
	g("downlinkthroughputbps", float64(st.DownlinkThroughputBps))
	g("uplinkthroughputbps", float64(st.UplinkThroughputBps))
	g("poppinglatencyms", float64(st.PopPingLatencyMs))
	if o := st.ObstructionStats; o != nil {
		g("obstructionstats.currentlyobstructed", b(o.CurrentlyObstructed))
		g("obstructionstats.fractionobstructed", float64(o.FractionObstructed))
		g("obstructionstats.valids", float64(o.ValidS))
		g("obstructionstats.avgprolongedobstructionvalid", b(o.AvgProlongedObstructionValid))
		if o.AvgProlongedObstructionValid {
			g("obstructionstats.avgprolongedobstructiondurrations", float64(o.AvgProlongedObstructionDurationS))
			g("obstructionstats.avgprolongedobstructionintervals", float64(o.AvgProlongedObstructionIntervalS))
		}
		g("obstructionstats.timeobstructed", float64(o.TimeObstructed))
		g("obstructionstats.patchesvalid", float64(o.PatchesValid))
	}
	g("signalquality", float64(st.SignalQuality))
	g("stowrequested", b(st.StowRequested))
	g("boresightazimuthdeg", float64(st.BoresightAzimuthDeg))
	g("boresightelevationdeg", float64(st.BoresightElevationDeg))
	g("ethspeedmbps", float64(st.EthSpeedMbps))
	g("mobilityclass", float64(st.MobilityClass), enumTag("mobility_class", st.MobilityClass))
	g("issnrabovenoisefloor", b(st.IsSnrAboveNoiseFloor))
	if r := st.ReadyStates; r != nil {
		g("readystates.cady", b(r.Cady))
		g("readystates.scp", b(r.Scp))
		g("readystates.l1l2", b(r.L1L2))
		g("readystates.xphy", b(r.Xphy))
		g("readystates.aap", b(r.Aap))
		g("readystates.rf", b(r.Rf))
	}
	g("issnrpersistentlylow", b(st.IsSnrPersistentlyLow))
	g("classofservice", float64(st.ClassOfService), enumTag("class_of_service", st.ClassOfService))
	g("softwareupdatestate", float64(st.SoftwareUpdateState), enumTag("software_update_state", st.SoftwareUpdateState))
	g("hasactuators", float64(st.HasActuators), enumTag("has_actuators", st.HasActuators))
	g("disablementcode", float64(st.DisablementCode), enumTag("disablement_code", st.DisablementCode))

	if a := st.AlignmentStats; a != nil {
		g("alignment.tiltangledeg", float64(a.TiltAngleDeg))
		g("alignment.boresightazimuthdeg", float64(a.BoresightAzimuthDeg))
		g("alignment.boresightelevationdeg", float64(a.BoresightElevationDeg))
		g("alignment.desiredboresightazimuthdeg", float64(a.DesiredBoresightAzimuthDeg))
		g("alignment.desiredboresightelevationdeg", float64(a.DesiredBoresightElevationDeg))
		g("alignment.attitudeuncertaintydeg", float64(a.AttitudeUncertaintyDeg))
		g("alignment.actuatorstate", float64(a.ActuatorState), enumTag("actuator_state", a.ActuatorState))
		g("alignment.attitudeestimationstate", float64(a.AttitudeEstimationState), enumTag("attitude_estimation_state", a.AttitudeEstimationState))
	}
	g("rebootreason", float64(st.RebootReason), enumTag("reboot_reason", st.RebootReason))
	g("dlbandwidthrestrictedreason", float64(st.DlBandwidthRestrictedReason), enumTag("dl_bandwidth_restricted_reason", st.DlBandwidthRestrictedReason))
	g("ulbandwidthrestrictedreason", float64(st.UlBandwidthRestrictedReason), enumTag("ul_bandwidth_restricted_reason", st.UlBandwidthRestrictedReason))
	g("ismovingfastpersisted", b(st.IsMovingFastPersisted))
	g("swupdaterebootready", b(st.SwupdateRebootReady))
	if s := st.SoftwareUpdateStats; s != nil {
		g("softwareupdate.progress", float64(s.SoftwareUpdateProgress))
	}
	if u := st.UpsuStats; u != nil {
		g("upsu.dishpowerw", float64(u.DishPower))
		g("upsu.routerpowerw", float64(u.RouterPower))
	}
	if i := st.InitializationDurationSeconds; i != nil {
		for _, s := range []struct {
			name string
			v    int32
		}{
			{"attitudeinitialization", i.AttitudeInitialization},
			{"burstdetected", i.BurstDetected},
			{"ekfconverged", i.EkfConverged},
			{"firstcplane", i.FirstCplane},
			{"firstpopping", i.FirstPopPing},
			{"gpsvalid", i.GpsValid},
			{"initialnetworkentry", i.InitialNetworkEntry},
			{"networkschedule", i.NetworkSchedule},
			{"rfready", i.RfReady},
			{"stableconnection", i.StableConnection},
		} {
			if s.v >= 0 {
				g("initduration."+s.name, float64(s.v))
			}
		}
	}
	if st.Outage != nil {
		g("outage.active", 1, enumTag("cause", st.Outage.Cause))
	} else {
		g("outage.active", 0, metrics.Tag("cause", "NONE"))
	}
}

// EmitHistory sends window statistics.
func EmitHistory(m *metrics.Sink, s Summary) {
	g := m.Gauge
	g("history.samples", float64(s.Samples))
	if s.Samples == 0 {
		return
	}
	g("history.droprate.mean", s.DropMean)
	g("history.seconds_down", float64(s.SecondsDown))
	if s.HasLatency {
		g("history.latency_ms.mean", s.LatencyMean)
		g("history.latency_ms.max", s.LatencyMax)
		g("history.latency_ms.p95", s.LatP95)
	}
	g("history.downlink_bps.mean", s.DownMean)
	g("history.downlink_bps.max", s.DownMax)
	g("history.uplink_bps.mean", s.UpMean)
	g("history.uplink_bps.max", s.UpMax)
	if s.HasPower {
		g("history.power_w.mean", s.PowerMean)
		g("history.power_w.max", s.PowerMax)
		g("history.power_w.min", s.PowerMin)
		g("history.energy_wh", s.EnergyWh)
	}
}
