// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.7.1
// source: api/proto/tradepipe.proto

package grpc

import (
	context "context"
	login "github.com/Sannrox/tradepipe/pkg/grpc/login"
	timeline "github.com/Sannrox/tradepipe/pkg/grpc/timeline"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TradePipeClient is the client API for TradePipe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TradePipeClient interface {
	Login(ctx context.Context, in *login.Credentials, opts ...grpc.CallOption) (*login.ProcessId, error)
	Verify(ctx context.Context, in *login.TwoFA, opts ...grpc.CallOption) (*login.TwoFAReturn, error)
	DownloadAll(ctx context.Context, in *timeline.DownloadAll, opts ...grpc.CallOption) (*timeline.DownloadAllResponse, error)
}

type tradePipeClient struct {
	cc grpc.ClientConnInterface
}

func NewTradePipeClient(cc grpc.ClientConnInterface) TradePipeClient {
	return &tradePipeClient{cc}
}

func (c *tradePipeClient) Login(ctx context.Context, in *login.Credentials, opts ...grpc.CallOption) (*login.ProcessId, error) {
	out := new(login.ProcessId)
	err := c.cc.Invoke(ctx, "/grpc.TradePipe/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tradePipeClient) Verify(ctx context.Context, in *login.TwoFA, opts ...grpc.CallOption) (*login.TwoFAReturn, error) {
	out := new(login.TwoFAReturn)
	err := c.cc.Invoke(ctx, "/grpc.TradePipe/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tradePipeClient) DownloadAll(ctx context.Context, in *timeline.DownloadAll, opts ...grpc.CallOption) (*timeline.DownloadAllResponse, error) {
	out := new(timeline.DownloadAllResponse)
	err := c.cc.Invoke(ctx, "/grpc.TradePipe/DownloadAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TradePipeServer is the server API for TradePipe service.
// All implementations must embed UnimplementedTradePipeServer
// for forward compatibility
type TradePipeServer interface {
	Login(context.Context, *login.Credentials) (*login.ProcessId, error)
	Verify(context.Context, *login.TwoFA) (*login.TwoFAReturn, error)
	DownloadAll(context.Context, *timeline.DownloadAll) (*timeline.DownloadAllResponse, error)
	mustEmbedUnimplementedTradePipeServer()
}

// UnimplementedTradePipeServer must be embedded to have forward compatible implementations.
type UnimplementedTradePipeServer struct {
}

func (UnimplementedTradePipeServer) Login(context.Context, *login.Credentials) (*login.ProcessId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedTradePipeServer) Verify(context.Context, *login.TwoFA) (*login.TwoFAReturn, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedTradePipeServer) DownloadAll(context.Context, *timeline.DownloadAll) (*timeline.DownloadAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadAll not implemented")
}
func (UnimplementedTradePipeServer) mustEmbedUnimplementedTradePipeServer() {}

// UnsafeTradePipeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TradePipeServer will
// result in compilation errors.
type UnsafeTradePipeServer interface {
	mustEmbedUnimplementedTradePipeServer()
}

func RegisterTradePipeServer(s grpc.ServiceRegistrar, srv TradePipeServer) {
	s.RegisterService(&TradePipe_ServiceDesc, srv)
}

func _TradePipe_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(login.Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradePipeServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TradePipe/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradePipeServer).Login(ctx, req.(*login.Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _TradePipe_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(login.TwoFA)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradePipeServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TradePipe/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradePipeServer).Verify(ctx, req.(*login.TwoFA))
	}
	return interceptor(ctx, in, info, handler)
}

func _TradePipe_DownloadAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(timeline.DownloadAll)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradePipeServer).DownloadAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TradePipe/DownloadAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradePipeServer).DownloadAll(ctx, req.(*timeline.DownloadAll))
	}
	return interceptor(ctx, in, info, handler)
}

// TradePipe_ServiceDesc is the grpc.ServiceDesc for TradePipe service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TradePipe_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.TradePipe",
	HandlerType: (*TradePipeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _TradePipe_Login_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _TradePipe_Verify_Handler,
		},
		{
			MethodName: "DownloadAll",
			Handler:    _TradePipe_DownloadAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/tradepipe.proto",
}
