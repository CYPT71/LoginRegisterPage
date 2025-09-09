package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/go-webauthn/webauthn/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterStartRequest carries the username initiating registration.
type RegisterStartRequest struct {
	Username string `json:"username"`
}

// RegisterStartResponse returns the WebAuthn options encoded as JSON.
type RegisterStartResponse struct {
	OptionsJSON string `json:"options_json"`
}

// RegisterEndRequest carries the credential returned by the client.
type RegisterEndRequest struct {
	Username   string `json:"username"`
	Credential []byte `json:"credential"`
}

// RegisterPasswordRequest creates a user with a password instead of WebAuthn.
type RegisterPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterStart begins the WebAuthn registration ceremony for a user.
func (a *AuthService) RegisterStart(ctx context.Context, req *RegisterStartRequest) (*RegisterStartResponse, error) {
	user := &domain.UserModel{Username: req.Username}
	if _, ok := utils.GetSession(user.Username); ok {
		utils.DeleteSession(user.Username)
	}
	if user.Find() {
		if len(user.Credentials) > 0 && user.Password == "" {
			return nil, status.Error(codes.AlreadyExists, "user already registered")
		}
		user.Delete()
	}
	user.Create()
	options, sessionData, err := utils.Web.BeginRegistration(*user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	session := &domain.UserSessions{
		DisplayName: options.Response.User.Name,
		SessionData: sessionData,
		Expiration:  time.Hour,
	}
	go session.DeleteAfter(utils.DeleteSession)
	utils.SaveSession(session)
	b, _ := json.Marshal(options)
	return &RegisterStartResponse{OptionsJSON: string(b)}, nil
}

// RegisterEnd completes the registration ceremony and returns a JWT token.
func (a *AuthService) RegisterEnd(ctx context.Context, req *RegisterEndRequest) (*TokenResponse, error) {
	user := &domain.UserModel{Username: req.Username}
	credential, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(req.Credential))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !user.Find() {
		return nil, status.Error(codes.NotFound, "not found")
	}
	session, ok := utils.GetSession(user.Username)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "session Not exist")
	}
	creds, err := utils.Web.CreateCredential(user, *session.SessionData, credential)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	user.Credentials = append(user.Credentials, *creds)
	user.SaveCredentials()
	session.SessionCred = creds
	session.Expiration = time.Hour * 48
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	session.Jwt = token
	go session.DeleteAfter(utils.DeleteSession)
	utils.DeleteSession(user.Username)
	utils.SaveSession(session)
	return &TokenResponse{Token: token}, nil
}

// RegisterPassword registers a user using a traditional password and returns a JWT.
func (a *AuthService) RegisterPassword(ctx context.Context, req *RegisterPasswordRequest) (*TokenResponse, error) {
	user := &domain.UserModel{Username: req.Username}
	if user.Find() {
		return nil, status.Error(codes.Unauthenticated, "not authorise")
	}
	if len(req.Password) <= 2 {
		return nil, status.Error(codes.InvalidArgument, "password to short")
	}
	if err := user.SetPassword(req.Password); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	session := &domain.UserSessions{DisplayName: user.Username, Expiration: time.Hour * 48}
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	session.Jwt = token
	utils.SaveSession(session)
	go session.DeleteAfter(utils.DeleteSession)
	user.Create()
	return &TokenResponse{Token: token}, nil
}
