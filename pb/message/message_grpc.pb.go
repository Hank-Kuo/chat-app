// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: pb/message/message.proto

package message

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	MessageReceived(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*ReceiveMessageResponse, error)
	ReplyReceived(ctx context.Context, in *ReceiveReplyRequest, opts ...grpc.CallOption) (*ReceiveReplyResponse, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) MessageReceived(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*ReceiveMessageResponse, error) {
	out := new(ReceiveMessageResponse)
	err := c.cc.Invoke(ctx, "/message.MessageService/MessageReceived", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) ReplyReceived(ctx context.Context, in *ReceiveReplyRequest, opts ...grpc.CallOption) (*ReceiveReplyResponse, error) {
	out := new(ReceiveReplyResponse)
	err := c.cc.Invoke(ctx, "/message.MessageService/ReplyReceived", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations should embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	MessageReceived(context.Context, *ReceiveMessageRequest) (*ReceiveMessageResponse, error)
	ReplyReceived(context.Context, *ReceiveReplyRequest) (*ReceiveReplyResponse, error)
}

// UnimplementedMessageServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) MessageReceived(context.Context, *ReceiveMessageRequest) (*ReceiveMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageReceived not implemented")
}
func (UnimplementedMessageServiceServer) ReplyReceived(context.Context, *ReceiveReplyRequest) (*ReceiveReplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplyReceived not implemented")
}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_MessageReceived_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).MessageReceived(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.MessageService/MessageReceived",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).MessageReceived(ctx, req.(*ReceiveMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_ReplyReceived_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveReplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).ReplyReceived(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.MessageService/ReplyReceived",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).ReplyReceived(ctx, req.(*ReceiveReplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MessageReceived",
			Handler:    _MessageService_MessageReceived_Handler,
		},
		{
			MethodName: "ReplyReceived",
			Handler:    _MessageService_ReplyReceived_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/message/message.proto",
}
