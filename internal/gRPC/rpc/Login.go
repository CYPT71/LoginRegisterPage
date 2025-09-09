package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/go-webauthn/webauthn/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginStartRequest carries the username initiating a login.
type LoginStartRequest struct {
	Username string `json:"username"`
}

// LoginStartResponse returns the options for the WebAuthn login ceremony encoded as JSON.
type LoginStartResponse struct {
	OptionsJSON string `json:"options_json"`
}

// LoginEndRequest carries the credential response returned by the client.
type LoginEndRequest struct {
	Username   string `json:"username"`
	Credential []byte `json:"credential"`
}

// TokenResponse wraps a JWT token for successful authentications.
type TokenResponse struct {
	Token string `json:"token"`
}

// AuthService exposes authentication related RPCs.
type AuthService struct{}

// LoginStart begins the WebAuthn login flow for a user and returns the options.
func (a *AuthService) LoginStart(ctx context.Context, req *LoginStartRequest) (*LoginStartResponse, error) {
	user := &domain.UserModel{Username: req.Username}
	if !user.Find() {
		return nil, status.Error(codes.Unauthenticated, "no user with this username")
	}
	user.ParseCredentials()

	options, sessionData, err := utils.Web.BeginLogin(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	session := &domain.UserSessions{
		SessionData: sessionData,
		DisplayName: user.Username,
		Expiration:  time.Minute * 5,
	}
	go session.DeleteAfter(utils.DeleteSession)
	utils.SaveSession(session)

	b, _ := json.Marshal(options)
	return &LoginStartResponse{OptionsJSON: string(b)}, nil
}

// LoginEnd completes the WebAuthn login flow and returns a JWT token.
func (a *AuthService) LoginEnd(ctx context.Context, req *LoginEndRequest) (*TokenResponse, error) {
	user := &domain.UserModel{Username: req.Username}
	if !user.Find() {
		return nil, status.Error(codes.Unauthenticated, "no user with this username")
	}
	user.ParseCredentials()

	session, ok := utils.GetSession(user.Username)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "session Not exist")
	}

	parsedCredential, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(req.Credential))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	creds, err := utils.Web.ValidateLogin(user, *session.SessionData, parsedCredential)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	session.SessionCred = creds
	session.Expiration = time.Minute * 48 * 60
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	session.Jwt = token
	go session.DeleteAfter(utils.DeleteSession)
	utils.SaveSession(session)
	user.Credentials = append(user.Credentials, *creds)
	user.SaveCredentials()

	return &TokenResponse{Token: token}, nil
}

// RegisterAuthServiceServer registers the AuthService on a gRPC server.
func RegisterAuthServiceServer(s *grpc.Server, srv *AuthService) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginStart",
			Handler:    _AuthService_LoginStart_Handler,
		},
		{
			MethodName: "LoginEnd",
			Handler:    _AuthService_LoginEnd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func _AuthService_LoginStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*AuthService).LoginStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/LoginStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*AuthService).LoginStart(ctx, req.(*LoginStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_LoginEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginEndRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*AuthService).LoginEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/LoginEnd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*AuthService).LoginEnd(ctx, req.(*LoginEndRequest))
	}
	return interceptor(ctx, in, info, handler)
}
