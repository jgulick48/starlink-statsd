#/bin/bash

protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/device.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd --go_opt=Mspacex/api/telemetron/public/common/time.proto=github.com/jgulick48/starlink-statsd
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/common/status/status.proto  --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/command.proto --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/common.proto --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/dish.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/common/protobuf --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/satellites/network
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/satellites/network/ut_disablement_codes.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/common/protobuf --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=spacex.com/api/satellites/network
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/common/protobuf/internal.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=spacex.com/api/common/protobuf
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/dish_config.proto
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/telemetron/public/common/time.proto --go_opt=Mspacex/api/telemetron/public/common/time.proto=spacex.com/api/telemetron/public/common
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/wifi.proto --go_opt=Mspacex/api/telemetron/public/common/time.proto=spacex.com/api/telemetron/public/common
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/wifi_config.proto
protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/transceiver.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device

find pkg/spacex.com -name "*.go" | xargs sed -i.bak 's|spacex/api|github.com/jgulick48/starlink-statsd/pkg/spacex.com/api|g'
find pkg/spacex.com -name "*.bak" | xargs rm