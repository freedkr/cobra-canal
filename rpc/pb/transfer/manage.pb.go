// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transfer/manage.proto

package transfer

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

func init() { proto.RegisterFile("transfer/manage.proto", fileDescriptor_8d652cc5f2f6026d) }

var fileDescriptor_8d652cc5f2f6026d = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8d, 0x31, 0x0b, 0x02, 0x21,
	0x18, 0x40, 0x9b, 0x24, 0x84, 0x20, 0x84, 0x96, 0x1b, 0x5b, 0xe3, 0xfc, 0xa0, 0x38, 0xda, 0xbb,
	0xa1, 0xa9, 0xa5, 0xb1, 0xed, 0x53, 0xd4, 0x84, 0x3b, 0x35, 0xfd, 0x8e, 0xa3, 0x7f, 0x1f, 0x48,
	0x8e, 0xef, 0xf1, 0xe0, 0xf1, 0x03, 0x65, 0x0c, 0xc5, 0x9a, 0x0c, 0x33, 0x06, 0x74, 0x46, 0xa6,
	0x1c, 0x29, 0x8a, 0x6d, 0xd3, 0x1d, 0x57, 0x58, 0xfe, 0xf6, 0x3c, 0x70, 0xf6, 0xa8, 0x95, 0x38,
	0x71, 0x76, 0x37, 0x34, 0x5a, 0x27, 0x76, 0xb2, 0x06, 0x4f, 0xf3, 0x59, 0x4c, 0xa1, 0x6e, 0xdf,
	0xb0, 0xa4, 0x31, 0x06, 0xeb, 0xdd, 0x71, 0x73, 0xbb, 0xbe, 0x06, 0xe7, 0xe9, 0xbd, 0x28, 0xa9,
	0xe3, 0x0c, 0x38, 0xad, 0xf8, 0x2d, 0xfd, 0x8a, 0x9e, 0x7c, 0x70, 0xa0, 0xa3, 0xca, 0xd8, 0x6b,
	0x0c, 0x38, 0x41, 0x4e, 0x1a, 0x92, 0x82, 0xf6, 0x56, 0xac, 0x6e, 0x2f, 0xbf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x25, 0x7f, 0x78, 0xf7, 0xa5, 0x00, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/transfer.Manage/GetCfg", in, out, opts...)
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
		FullMethod: "/transfer.Manage/GetCfg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServer).GetCfg(ctx, req.(*pb.Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Manage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "transfer.Manage",
	HandlerType: (*ManageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCfg",
			Handler:    _Manage_GetCfg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transfer/manage.proto",
}
