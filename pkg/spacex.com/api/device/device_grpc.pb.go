// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package device

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DeviceClient is the client API for Device service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (Device_StreamClient, error)
	Handle(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type deviceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceClient(cc grpc.ClientConnInterface) DeviceClient {
	return &deviceClient{cc}
}

func (c *deviceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Device_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Device_ServiceDesc.Streams[0], "/SpaceX.API.Device.Device/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceStreamClient{stream}
	return x, nil
}

type Device_StreamClient interface {
	Send(*ToDevice) error
	Recv() (*FromDevice, error)
	grpc.ClientStream
}

type deviceStreamClient struct {
	grpc.ClientStream
}

func (x *deviceStreamClient) Send(m *ToDevice) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deviceStreamClient) Recv() (*FromDevice, error) {
	m := new(FromDevice)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deviceClient) Handle(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/SpaceX.API.Device.Device/Handle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServer is the server API for Device service.
// All implementations must embed UnimplementedDeviceServer
// for forward compatibility
type DeviceServer interface {
	Stream(Device_StreamServer) error
	Handle(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedDeviceServer()
}

// UnimplementedDeviceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServer struct {
}

func (UnimplementedDeviceServer) Stream(Device_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedDeviceServer) Handle(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Handle not implemented")
}
func (UnimplementedDeviceServer) mustEmbedUnimplementedDeviceServer() {}

// UnsafeDeviceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServer will
// result in compilation errors.
type UnsafeDeviceServer interface {
	mustEmbedUnimplementedDeviceServer()
}

func RegisterDeviceServer(s grpc.ServiceRegistrar, srv DeviceServer) {
	s.RegisterService(&Device_ServiceDesc, srv)
}

func _Device_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeviceServer).Stream(&deviceStreamServer{stream})
}

type Device_StreamServer interface {
	Send(*FromDevice) error
	Recv() (*ToDevice, error)
	grpc.ServerStream
}

type deviceStreamServer struct {
	grpc.ServerStream
}

func (x *deviceStreamServer) Send(m *FromDevice) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deviceStreamServer) Recv() (*ToDevice, error) {
	m := new(ToDevice)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Device_Handle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServer).Handle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SpaceX.API.Device.Device/Handle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServer).Handle(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Device_ServiceDesc is the grpc.ServiceDesc for Device service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Device_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SpaceX.API.Device.Device",
	HandlerType: (*DeviceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Handle",
			Handler:    _Device_Handle_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Device_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device/device.proto",
}
