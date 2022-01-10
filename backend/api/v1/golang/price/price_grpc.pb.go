// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package price

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

// PriceClient is the client API for Price service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PriceClient interface {
	// Get price per gram of the metal on LME
	Price(ctx context.Context, in *PriceRequest, opts ...grpc.CallOption) (*PriceResponse, error)
}

type priceClient struct {
	cc grpc.ClientConnInterface
}

func NewPriceClient(cc grpc.ClientConnInterface) PriceClient {
	return &priceClient{cc}
}

func (c *priceClient) Price(ctx context.Context, in *PriceRequest, opts ...grpc.CallOption) (*PriceResponse, error) {
	out := new(PriceResponse)
	err := c.cc.Invoke(ctx, "/core.backend.integration.api.v1.price.Price/Price", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PriceServer is the server API for Price service.
// All implementations must embed UnimplementedPriceServer
// for forward compatibility
type PriceServer interface {
	// Get price per gram of the metal on LME
	Price(context.Context, *PriceRequest) (*PriceResponse, error)
	mustEmbedUnimplementedPriceServer()
}

// UnimplementedPriceServer must be embedded to have forward compatible implementations.
type UnimplementedPriceServer struct {
}

func (UnimplementedPriceServer) Price(context.Context, *PriceRequest) (*PriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Price not implemented")
}
func (UnimplementedPriceServer) mustEmbedUnimplementedPriceServer() {}

// UnsafePriceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PriceServer will
// result in compilation errors.
type UnsafePriceServer interface {
	mustEmbedUnimplementedPriceServer()
}

func RegisterPriceServer(s grpc.ServiceRegistrar, srv PriceServer) {
	s.RegisterService(&Price_ServiceDesc, srv)
}

func _Price_Price_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PriceServer).Price(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.backend.integration.api.v1.price.Price/Price",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PriceServer).Price(ctx, req.(*PriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Price_ServiceDesc is the grpc.ServiceDesc for Price service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Price_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.backend.integration.api.v1.price.Price",
	HandlerType: (*PriceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Price",
			Handler:    _Price_Price_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "price.proto",
}
