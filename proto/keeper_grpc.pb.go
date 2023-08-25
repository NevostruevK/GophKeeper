// For generation use:
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/keeper.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/keeper.proto

package proto

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
	Keeper_GetSpecs_FullMethodName       = "/keeper.Keeper/GetSpecs"
	Keeper_GetSpecsOfType_FullMethodName = "/keeper.Keeper/GetSpecsOfType"
	Keeper_GetData_FullMethodName        = "/keeper.Keeper/GetData"
	Keeper_GetDescription_FullMethodName = "/keeper.Keeper/GetDescription"
	Keeper_AddRecord_FullMethodName      = "/keeper.Keeper/AddRecord"
)

// KeeperClient is the client API for Keeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperClient interface {
	GetSpecs(ctx context.Context, in *GetSpecsRequest, opts ...grpc.CallOption) (*Specs, error)
	GetSpecsOfType(ctx context.Context, in *GetSpecsOfTypeRequest, opts ...grpc.CallOption) (*Specs, error)
	GetData(ctx context.Context, in *DataSpec, opts ...grpc.CallOption) (*Data, error)
	GetDescription(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*Data, error)
	AddRecord(ctx context.Context, in *Record, opts ...grpc.CallOption) (*RecordID, error)
}

type keeperClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperClient(cc grpc.ClientConnInterface) KeeperClient {
	return &keeperClient{cc}
}

func (c *keeperClient) GetSpecs(ctx context.Context, in *GetSpecsRequest, opts ...grpc.CallOption) (*Specs, error) {
	out := new(Specs)
	err := c.cc.Invoke(ctx, Keeper_GetSpecs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) GetSpecsOfType(ctx context.Context, in *GetSpecsOfTypeRequest, opts ...grpc.CallOption) (*Specs, error) {
	out := new(Specs)
	err := c.cc.Invoke(ctx, Keeper_GetSpecsOfType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) GetData(ctx context.Context, in *DataSpec, opts ...grpc.CallOption) (*Data, error) {
	out := new(Data)
	err := c.cc.Invoke(ctx, Keeper_GetData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) GetDescription(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*Data, error) {
	out := new(Data)
	err := c.cc.Invoke(ctx, Keeper_GetDescription_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) AddRecord(ctx context.Context, in *Record, opts ...grpc.CallOption) (*RecordID, error) {
	out := new(RecordID)
	err := c.cc.Invoke(ctx, Keeper_AddRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServer is the server API for Keeper service.
// All implementations must embed UnimplementedKeeperServer
// for forward compatibility
type KeeperServer interface {
	GetSpecs(context.Context, *GetSpecsRequest) (*Specs, error)
	GetSpecsOfType(context.Context, *GetSpecsOfTypeRequest) (*Specs, error)
	GetData(context.Context, *DataSpec) (*Data, error)
	GetDescription(context.Context, *RecordID) (*Data, error)
	AddRecord(context.Context, *Record) (*RecordID, error)
	mustEmbedUnimplementedKeeperServer()
}

// UnimplementedKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServer struct {
}

func (UnimplementedKeeperServer) GetSpecs(context.Context, *GetSpecsRequest) (*Specs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSpecs not implemented")
}
func (UnimplementedKeeperServer) GetSpecsOfType(context.Context, *GetSpecsOfTypeRequest) (*Specs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSpecsOfType not implemented")
}
func (UnimplementedKeeperServer) GetData(context.Context, *DataSpec) (*Data, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedKeeperServer) GetDescription(context.Context, *RecordID) (*Data, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDescription not implemented")
}
func (UnimplementedKeeperServer) AddRecord(context.Context, *Record) (*RecordID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRecord not implemented")
}
func (UnimplementedKeeperServer) mustEmbedUnimplementedKeeperServer() {}

// UnsafeKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperServer will
// result in compilation errors.
type UnsafeKeeperServer interface {
	mustEmbedUnimplementedKeeperServer()
}

func RegisterKeeperServer(s grpc.ServiceRegistrar, srv KeeperServer) {
	s.RegisterService(&Keeper_ServiceDesc, srv)
}

func _Keeper_GetSpecs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSpecsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetSpecs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_GetSpecs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetSpecs(ctx, req.(*GetSpecsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_GetSpecsOfType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSpecsOfTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetSpecsOfType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_GetSpecsOfType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetSpecsOfType(ctx, req.(*GetSpecsOfTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_GetData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetData(ctx, req.(*DataSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_GetDescription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetDescription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_GetDescription_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetDescription(ctx, req.(*RecordID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_AddRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Record)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).AddRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_AddRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).AddRecord(ctx, req.(*Record))
	}
	return interceptor(ctx, in, info, handler)
}

// Keeper_ServiceDesc is the grpc.ServiceDesc for Keeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Keeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keeper.Keeper",
	HandlerType: (*KeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSpecs",
			Handler:    _Keeper_GetSpecs_Handler,
		},
		{
			MethodName: "GetSpecsOfType",
			Handler:    _Keeper_GetSpecsOfType_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _Keeper_GetData_Handler,
		},
		{
			MethodName: "GetDescription",
			Handler:    _Keeper_GetDescription_Handler,
		},
		{
			MethodName: "AddRecord",
			Handler:    _Keeper_AddRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/keeper.proto",
}