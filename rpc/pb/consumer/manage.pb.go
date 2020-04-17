// Code generated by protoc-gen-go. DO NOT EDIT.
// source: consumer/manage.proto

package consumer

import (
	context "context"
	fmt "fmt"
	pb "github.com/always-waiting/cobra-canal/rpc/pb"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("consumer/manage.proto", fileDescriptor_339f83969384741a) }

var fileDescriptor_339f83969384741a = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8d, 0x3d, 0x0b, 0xc2, 0x30,
	0x14, 0x00, 0x9d, 0x8a, 0x04, 0x04, 0x09, 0xb8, 0x74, 0x74, 0x95, 0xe6, 0x81, 0x52, 0xdc, 0xed,
	0xe0, 0xe4, 0xe2, 0xe8, 0xf6, 0x12, 0xd2, 0x18, 0x68, 0xf2, 0x62, 0x3e, 0x28, 0xfe, 0x7b, 0x21,
	0x98, 0xf1, 0x8e, 0x83, 0x63, 0x07, 0x45, 0x3e, 0x15, 0xa7, 0x23, 0x38, 0xf4, 0x68, 0xb4, 0x08,
	0x91, 0x32, 0xf1, 0x6d, 0xd3, 0x3d, 0x93, 0x98, 0xfe, 0xf6, 0x3c, 0xb2, 0xee, 0x51, 0x2b, 0x7e,
	0x62, 0xdd, 0x5d, 0xe7, 0x69, 0x36, 0x7c, 0x27, 0x6a, 0xf0, 0xd4, 0x9f, 0xa2, 0x53, 0xee, 0xf7,
	0x0d, 0x53, 0x98, 0xc8, 0xcf, 0xd6, 0x1c, 0x37, 0xb7, 0xeb, 0x6b, 0x34, 0x36, 0xbf, 0x8b, 0x14,
	0x8a, 0x1c, 0xe0, 0xb2, 0xe2, 0x37, 0x0d, 0x2b, 0xda, 0x6c, 0xbd, 0x01, 0x45, 0x32, 0xe2, 0xa0,
	0xd0, 0xe3, 0x02, 0x31, 0x28, 0x08, 0x12, 0xda, 0x5b, 0x76, 0x75, 0x7b, 0xf9, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x49, 0x3b, 0x88, 0xa8, 0xa5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ManageClient is the client API for Manage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ManageClient interface {
	GetCfg(ctx context.Context, in *pb.Request, opts ...grpc.CallOption) (*pb.RespConfig, error)
}

type manageClient struct {
	cc *grpc.ClientConn
}

func NewManageClient(cc *grpc.ClientConn) ManageClient {
	return &manageClient{cc}
}

func (c *manageClient) GetCfg(ctx context.Context, in *pb.Request, opts ...grpc.CallOption) (*pb.RespConfig, error) {
	out := new(pb.RespConfig)
	err := c.cc.Invoke(ctx, "/consumer.Manage/GetCfg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManageServer is the server API for Manage service.
type ManageServer interface {
	GetCfg(context.Context, *pb.Request) (*pb.RespConfig, error)
}

// UnimplementedManageServer can be embedded to have forward compatible implementations.
type UnimplementedManageServer struct {
}

func (*UnimplementedManageServer) GetCfg(ctx context.Context, req *pb.Request) (*pb.RespConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCfg not implemented")
}

func RegisterManageServer(s *grpc.Server, srv ManageServer) {
	s.RegisterService(&_Manage_serviceDesc, srv)
}

func _Manage_GetCfg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(pb.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServer).GetCfg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/consumer.Manage/GetCfg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServer).GetCfg(ctx, req.(*pb.Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Manage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "consumer.Manage",
	HandlerType: (*ManageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCfg",
			Handler:    _Manage_GetCfg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consumer/manage.proto",
}
