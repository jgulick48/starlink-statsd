#/bin/bash

/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset \
 --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd \
 --go_opt=Mspacex/api/telemetron/public/common/time.proto=github.com/jgulick48/starlink-statsd \
 --go_opt=Mspacex/api/device/services/unlock/service.proto=github.com/jgulick48/starlink-statsd \
 --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=spacex/api/satellites/network \
 spacex/api/device/device.proto
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/common/status/status.proto  --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/command.proto --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/common.proto --go_opt=Mproto/dish.protoset=github.com/jgulick48/starlink-statsd
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/dish.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/common/protobuf --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/satellites/network
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/satellites/network/ut_disablement_codes.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/common/protobuf --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=spacex.com/api/satellites/network
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/common/protobuf/internal.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=spacex.com/api/common/protobuf
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/dish_config.proto
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/telemetron/public/common/time.proto --go_opt=Mspacex/api/telemetron/public/common/time.proto=spacex.com/api/telemetron/public/common
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/wifi.proto --go_opt=Mspacex/api/telemetron/public/common/time.proto=spacex.com/api/telemetron/public/common
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/wifi_config.proto
/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset spacex/api/device/transceiver.proto --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device

find pkg/spacex.com -name "*.go" | xargs sed -i.bak 's|spacex/api|github.com/jgulick48/starlink-statsd/pkg/spacex.com/api|g'
find pkg/spacex.com -name "*.bak" | xargs rm



/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/ --go-grpc_out=pkg/ --descriptor_set_in=proto/dish.protoset spacex/api/satellites/network/ut_disablement_codes.proto \
 --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd \


/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=proto/dish.protoset \
 --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=github.com/jgulick48/starlink-statsd/pkg/spacex/satellites/network \
 spacex/api/device/dish.proto


/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/    --go-grpc_out=pkg/ \
 --descriptor_set_in=proto/dish.protoset spacex/api/satellites/network/ut_disablement_codes.proto \
 --go_opt=Mspacex/api/common/protobuf/internal.proto=github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/common/protobuf \
 --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=spacex.com/api/satellites/network


/opt/homebrew/opt/protobuf@3.20/bin/protoc --go_out=pkg/ --go-grpc_out=pkg/ --go-grpc_opt=paths=source_relative --descriptor_set_in=proto/dish.protoset \
 --go_opt=Mspacex/api/satellites/network/ut_disablement_codes.proto=spacex.com/api/satellites/network \
 spacex/api/satellites/network/ut_disablement_codes.proto