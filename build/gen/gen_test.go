package main

import (
	"testing"

	"github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestGeneratedFields(t *testing.T) {
	cases := []struct {
		msg    protoreflect.ProtoMessage
		fields []string
	}{
		{&device.DishGetStatusResponse{}, []string{"alignment_stats", "signal_quality", "reboot_reason", "upsu_stats", "software_update_stats", "initialization_duration_seconds"}},
		{&device.DishGetHistoryResponse{}, []string{"power_in", "outages", "current"}},
		{&device.DishAlerts{}, []string{"dish_water_detected", "lower_signal_than_predicted", "no_ethernet_link"}},
		{&device.Request{}, []string{"get_status", "get_history"}},
	}
	for _, c := range cases {
		md := c.msg.ProtoReflect().Descriptor()
		for _, f := range c.fields {
			if md.Fields().ByName(protoreflect.Name(f)) == nil {
				t.Errorf("%s missing field %s", md.Name(), f)
			}
		}
	}
	od := (&device.Request{}).ProtoReflect().Descriptor().Oneofs().ByName("request")
	if od == nil {
		t.Fatal("Request has no oneof named request")
	}
	for _, f := range []string{"get_status", "get_history"} {
		fd := od.Fields().ByName(protoreflect.Name(f))
		if fd == nil {
			t.Errorf("oneof missing %s", f)
		}
	}
}
