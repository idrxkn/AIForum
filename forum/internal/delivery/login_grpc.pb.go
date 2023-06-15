package delivery

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// LoginServiceClientImpl is a client for interacting with the LoginService service.
type LoginServiceClientImpl struct {
	cc *grpc.ClientConn
}

// NewLoginServiceClientImpl creates a new LoginServiceClientImpl.
func NewLoginServiceClientImpl(cc *grpc.ClientConn) *LoginServiceClientImpl {
	return &LoginServiceClientImpl{cc}
}

// Login sends a login request.
func (c *LoginServiceClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := grpc.Invoke(ctx, "/delivery.LoginService/Login", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginServiceServerImpl is the server API for LoginService service.
type LoginServiceServerImpl struct {
}

// NewLoginServiceServerImpl creates a new LoginServiceServerImpl.
func NewLoginServiceServerImpl() *LoginServiceServerImpl {
	return &LoginServiceServerImpl{}
}

// Login handles login requests.
func (s *LoginServiceServerImpl) Login(ctx context.Context, in *LoginRequest) (*LoginResponse, error) {
	return nil, grpc.Errorf(grpc.Code(grpc.Unimplemented), "method Login not implemented")
}

func registerLoginServiceServer(s *grpc.Server, srv LoginServiceServer) {
	s.RegisterService(&_LoginService_serviceDesc, srv)
}

var _LoginService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "delivery.LoginService",
	HandlerType: (*LoginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _LoginService_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "delivery/login.proto",
}

func _LoginService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/delivery.LoginService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}
