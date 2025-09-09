package rpc

import (
	"context"
	"strings"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Empty is used for methods that do not require input.
type Empty struct{}

// MessageResponse wraps a simple message string.
type MessageResponse struct {
	Message string `json:"message"`
}

// EditUserRequest carries fields to update a user.
type EditUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UsersResponse returns a list of users.
type UsersResponse struct {
	Users []domain.UserModel `json:"users"`
}

// UserService exposes user-related RPCs.
type UserService struct{}

// sessionFromContext extracts the session based on Authorization metadata.
func sessionFromContext(ctx context.Context) (*domain.UserSessions, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization")
	}
	parts := strings.SplitN(vals[0], " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
	}
	sess, ok := utils.GetSessionByToken(parts[1])
	if !ok || !utils.CheckJWT(sess, parts[1]) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return sess, nil
}

// About returns information about the authenticated user.
func (u *UserService) About(ctx context.Context, _ *Empty) (*domain.UserModel, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &domain.UserModel{Username: sess.DisplayName}
	return user.Get(), nil
}

// All returns all users if the requester has owner permission.
func (u *UserService) All(ctx context.Context, _ *Empty) (*UsersResponse, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &domain.UserModel{Username: sess.DisplayName}
	if user == nil || user.Get().Permission&domain.Permissions["owner"] != 1 {
		return nil, status.Error(codes.PermissionDenied, "not owner")
	}
	users, err := domain.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &UsersResponse{Users: users}, nil
}

// Logout deletes the current session.
func (u *UserService) Logout(ctx context.Context, _ *Empty) (*MessageResponse, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	utils.DeleteSession(sess.DisplayName)
	return &MessageResponse{Message: "logout"}, nil
}

// Edit updates the authenticated user's details.
func (u *UserService) Edit(ctx context.Context, req *EditUserRequest) (*domain.UserModel, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &domain.UserModel{Username: sess.DisplayName}
	user = user.Get()
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	user.Update()
	return user, nil
}

// Delete removes the authenticated user's account.
func (u *UserService) Delete(ctx context.Context, _ *Empty) (*MessageResponse, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &domain.UserModel{Username: sess.DisplayName}
	user.Delete()
	utils.DeleteSession(user.Username)
	return &MessageResponse{Message: "deleted"}, nil
}

// DeleteCred removes the first credential from the user.
func (u *UserService) DeleteCred(ctx context.Context, _ *Empty) (*domain.UserModel, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &domain.UserModel{Username: sess.DisplayName}
	user = user.Get()
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}
	user.Incredentials = strings.Split(user.Incredentials, ";")[0]
	user.Update()
	return user, nil
}

// RegisterUserServiceServer registers the UserService on a gRPC server.
func RegisterUserServiceServer(s *grpc.Server, srv *UserService) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auth.UserService",
	HandlerType: (*UserService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "About", Handler: _UserService_About_Handler},
		{MethodName: "All", Handler: _UserService_All_Handler},
		{MethodName: "Logout", Handler: _UserService_Logout_Handler},
		{MethodName: "Edit", Handler: _UserService_Edit_Handler},
		{MethodName: "Delete", Handler: _UserService_Delete_Handler},
		{MethodName: "DeleteCred", Handler: _UserService_DeleteCred_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func _UserService_About_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).About(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/About"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).About(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_All_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).All(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/All"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).All(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/Logout"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).Logout(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/Edit"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).Edit(ctx, req.(*EditUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/Delete"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).Delete(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteCred_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*UserService).DeleteCred(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.UserService/DeleteCred"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*UserService).DeleteCred(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
