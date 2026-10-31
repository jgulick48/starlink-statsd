# starlink-statsd

Polls a Starlink dish over gRPC (read-only `get_status` and `get_history`
requests only) and emits DogStatsD metrics, intended for Telegraf's `statsd`
input (`datadog_extensions = true`, `metric_separator = "_"`) into InfluxDB.

Every metric is prefixed `starlink.` and tagged `deploymentId:<id>`. Non-finite
values (NaN/Inf) are never sent. RPC failures are logged (repeats rate-limited
to once a minute) and never stop the exporter.

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-addr` | `192.168.100.1:9200` | Dish gRPC address |
| `-statsDHost` | `127.0.0.1:8125` | DogStatsD address |
| `-deploymentId` | `Dishy` | Value of the `deploymentId` tag |
| `-scanInterval` | `30` | Poll interval in seconds |
| `-history` | `true` | Also fetch `get_history` each poll |

Version and commit are set at build time with
`-ldflags "-X main.version=... -X main.commit=..."` and logged at startup.

## Metrics

Collector: `collector.poll` (counter, tag `result:success|error`),
`collector.poll_duration_ms` (timing).

Status (existing names, unchanged): `deviceinfo` (tags id, hardwareVersion,
softwareVersion, manufacturedVersion, countryCode), `generationnumber`,
`utcoffsets`, `bootcount`, `antirollbackversion`, `uptimes`, `alerts`
(tag `alertType`), `gpsStats.gpsvalid`, `gpsStats.gpssats`,
`secondstofirstnonemptyslot`, `poppingdroprate`, `downlinkthroughputbps`,
`uplinkthroughputbps`, `poppinglatencyms`,
`obstructionstats.{currentlyobstructed,fractionobstructed,valids,avgprolongedobstructionvalid}`,
`obstructionstats.avgprolongedobstructiondurrations` and
`...intervals` (only when valid), `stowrequested`, `boresightazimuthdeg`,
`boresightelevationdeg`, `ethspeedmbps`, `mobilityclass`,
`issnrabovenoisefloor`, `readystates.{cady,scp,l1l2,xphy,aap,rf}`,
`issnrpersistentlylow`, `classofservice`, `softwareupdatestate`,
`hasactuators`, `disablementcode`. Enum metrics carry the enum number as the
value and a `key:NAME` tag (`mobility_class`, `class_of_service`,
`software_update_state`, `has_actuators`, `disablement_code`).
The `Roaming` alert was removed from the dish API and is no longer sent.

Status (new): more `alerts` types (DishWaterDetected, RouterWaterDetected,
LowerSignalThanPredicted, NoEthernetLink, ObstructionMapReset,
SlowEthernetSpeeds100, UpsuRouterPortSlow, LowMotorCurrent, DbfTelemStale,
IsPowerSaveIdle), `signalquality`, `obstructionstats.timeobstructed`,
`obstructionstats.patchesvalid`,
`alignment.{tiltangledeg,boresightazimuthdeg,boresightelevationdeg,desiredboresightazimuthdeg,desiredboresightelevationdeg,attitudeuncertaintydeg}`,
`alignment.actuatorstate` / `alignment.attitudeestimationstate` (enum number,
tags `actuator_state` / `attitude_estimation_state`), `rebootreason`,
`dlbandwidthrestrictedreason`, `ulbandwidthrestrictedreason`,
`ismovingfastpersisted`, `swupdaterebootready`, `softwareupdate.progress`
(0..1), `upsu.dishpowerw`, `upsu.routerpowerw`, `initduration.<stage>`
(seconds, unreached stages skipped), `outage.active` (1 with `cause:<NAME>`
during an outage, else 0 with `cause:NONE`).

History (with `-history`): the dish keeps 1 sample/second ring buffers. Each
poll summarizes only samples written since the previous poll (first poll or a
dish reboot: the last scan interval). `history.samples`, `history.seconds_down`,
`history.droprate.mean`, `history.latency_ms.{mean,max,p95}` (drop-rate 1
samples excluded), `history.downlink_bps.{mean,max}`,
`history.uplink_bps.{mean,max}`, `history.power_w.{mean,max,min}`,
`history.energy_wh`.

Outages from history are reported once each (deduplicated by start time;
outages already present at startup are not replayed): a DogStatsD event
`Starlink outage: <CAUSE>` (Warning if longer than 30s, else Info), counter
`outage.count` and gauge `outage.duration_s`, all tagged `cause:<CAUSE>`.

## Development

```
make build       # build bin/starlink-statsd with version/commit ldflags
make test        # run tests
make vet         # go vet
make protoset    # refresh proto/dish.protoset from a live dish
make generate    # regenerate pkg/ from proto/dish.protoset
make help        # list all targets
```
