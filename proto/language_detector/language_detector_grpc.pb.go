// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: proto/language_detector/language_detector.proto

package languageDetector

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

const (
	LanguageDetector_DetectLanguage_FullMethodName = "/languageDetector.LanguageDetector/DetectLanguage"
	LanguageDetector_HealthCheck_FullMethodName    = "/languageDetector.LanguageDetector/HealthCheck"
)

// LanguageDetectorClient is the client API for LanguageDetector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LanguageDetectorClient interface {
	DetectLanguage(ctx context.Context, opts ...grpc.CallOption) (LanguageDetector_DetectLanguageClient, error)
	HealthCheck(ctx context.Context, in *PingPong, opts ...grpc.CallOption) (*PingPong, error)
}

type languageDetectorClient struct {
	cc grpc.ClientConnInterface
}

func NewLanguageDetectorClient(cc grpc.ClientConnInterface) LanguageDetectorClient {
	return &languageDetectorClient{cc}
}

func (c *languageDetectorClient) DetectLanguage(ctx context.Context, opts ...grpc.CallOption) (LanguageDetector_DetectLanguageClient, error) {
	stream, err := c.cc.NewStream(ctx, &LanguageDetector_ServiceDesc.Streams[0], LanguageDetector_DetectLanguage_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &languageDetectorDetectLanguageClient{stream}
	return x, nil
}

type LanguageDetector_DetectLanguageClient interface {
	Send(*LanguageDetectRequest) error
	Recv() (*LanguageDetectResponse, error)
	grpc.ClientStream
}

type languageDetectorDetectLanguageClient struct {
	grpc.ClientStream
}

func (x *languageDetectorDetectLanguageClient) Send(m *LanguageDetectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *languageDetectorDetectLanguageClient) Recv() (*LanguageDetectResponse, error) {
	m := new(LanguageDetectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *languageDetectorClient) HealthCheck(ctx context.Context, in *PingPong, opts ...grpc.CallOption) (*PingPong, error) {
	out := new(PingPong)
	err := c.cc.Invoke(ctx, LanguageDetector_HealthCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LanguageDetectorServer is the server API for LanguageDetector service.
// All implementations must embed UnimplementedLanguageDetectorServer
// for forward compatibility
type LanguageDetectorServer interface {
	DetectLanguage(LanguageDetector_DetectLanguageServer) error
	HealthCheck(context.Context, *PingPong) (*PingPong, error)
	mustEmbedUnimplementedLanguageDetectorServer()
}

// UnimplementedLanguageDetectorServer must be embedded to have forward compatible implementations.
type UnimplementedLanguageDetectorServer struct {
}

func (UnimplementedLanguageDetectorServer) DetectLanguage(LanguageDetector_DetectLanguageServer) error {
	return status.Errorf(codes.Unimplemented, "method DetectLanguage not implemented")
}
func (UnimplementedLanguageDetectorServer) HealthCheck(context.Context, *PingPong) (*PingPong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedLanguageDetectorServer) mustEmbedUnimplementedLanguageDetectorServer() {}

// UnsafeLanguageDetectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LanguageDetectorServer will
// result in compilation errors.
type UnsafeLanguageDetectorServer interface {
	mustEmbedUnimplementedLanguageDetectorServer()
}

func RegisterLanguageDetectorServer(s grpc.ServiceRegistrar, srv LanguageDetectorServer) {
	s.RegisterService(&LanguageDetector_ServiceDesc, srv)
}

func _LanguageDetector_DetectLanguage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LanguageDetectorServer).DetectLanguage(&languageDetectorDetectLanguageServer{stream})
}

type LanguageDetector_DetectLanguageServer interface {
	Send(*LanguageDetectResponse) error
	Recv() (*LanguageDetectRequest, error)
	grpc.ServerStream
}

type languageDetectorDetectLanguageServer struct {
	grpc.ServerStream
}

func (x *languageDetectorDetectLanguageServer) Send(m *LanguageDetectResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *languageDetectorDetectLanguageServer) Recv() (*LanguageDetectRequest, error) {
	m := new(LanguageDetectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _LanguageDetector_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingPong)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LanguageDetectorServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LanguageDetector_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LanguageDetectorServer).HealthCheck(ctx, req.(*PingPong))
	}
	return interceptor(ctx, in, info, handler)
}

// LanguageDetector_ServiceDesc is the grpc.ServiceDesc for LanguageDetector service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LanguageDetector_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "languageDetector.LanguageDetector",
	HandlerType: (*LanguageDetectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _LanguageDetector_HealthCheck_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DetectLanguage",
			Handler:       _LanguageDetector_DetectLanguage_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/language_detector/language_detector.proto",
}
