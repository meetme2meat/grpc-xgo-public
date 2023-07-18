package grpc

import (
	"bytes"
	"context"
	"crypto/rsa"
	"fmt"
	"time"
	"xgo/auth/src/gen"
	"xgo/auth/src/internal/handler/repository"
	"xgo/auth/src/model"

	jwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedAuthServer
	jwtPrivatekey *rsa.PrivateKey
	repo          *repository.Repository
}

func (h *Handler) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	fmt.Println("received login request")
	user := model.User{
		Username: req.Username,
	}

	ok, err := h.repo.Login(ctx, user, req.Password)
	if err != nil && !ok {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	fmt.Println("user exist", ok)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
		"sub": user.Username,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(h.jwtPrivatekey)
	fmt.Println("token set")

	if err != nil {
		fmt.Println("we got error err", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.LoginResponse{Token: tokenString}, nil
}

func (h *Handler) Signup(ctx context.Context, req *gen.SignupRequest) (*gen.SignupResponse, error) {
	fmt.Println("received signup request")

	password := []byte(req.Password)
	passwordConfirm := []byte(req.PasswordConfirm)
	if result := bytes.Compare(password, passwordConfirm); result != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password data")
	}

	user := model.User{Username: req.Username}
	err := h.repo.Create(ctx, user, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.SignupResponse{Message: ""}, nil
}

func New(key *rsa.PrivateKey, repo *repository.Repository) *Handler {
	return &Handler{
		jwtPrivatekey: key,
		repo:          repo,
	}
}
