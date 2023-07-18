package grpc

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"
	"xgo/main/src/gen"
	model "xgo/main/src/models"

	"xgo/main/src/internal/controller/company"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedCompanyServiceServer
	jwtPublicKey *rsa.PublicKey
	ctrl         *company.Controller
}

func New(ctrl *company.Controller, jwtPublicKey *rsa.PublicKey) *Handler {
	return &Handler{ctrl: ctrl, jwtPublicKey: jwtPublicKey}
}

func (h *Handler) verifyJWT(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "valid token required.")
	}

	jwtToken, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "valid token required.")
	}

	return validateToken(jwtToken[0], h.jwtPublicKey)
}

func validateToken(token string, publicKey *rsa.PublicKey) error {

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return errors.New("token invalid")
	}

	expiredTime, err := jwtToken.Claims.GetExpirationTime()
	if err != nil {
		return errors.New("invalid expiration time")
	}

	if expiredTime.Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func (h *Handler) GetCompany(ctx context.Context, req *gen.GetCompanyRequest) (*gen.GetCompanyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	c, err := h.ctrl.GetCompany(ctx, req.Id)
	if err != nil && errors.Is(err, company.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetCompanyResponse{
		Company: &gen.Company{
			Id:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Employee:    c.EmployeeCount,
			Registered:  c.Registered,
			Type:        gen.Type(gen.Type_value[c.Type]),
		},
	}, nil
}

func (h *Handler) CreateCompany(ctx context.Context, req *gen.CreateCompanyRequest) (*gen.CreateCompanyResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid data")
	}
	// // Please go through README
	// err := h.verifyJWT(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	company := model.Company{
		ID:            req.Company.Id,
		Name:          req.Company.Name,
		Description:   req.Company.Description,
		EmployeeCount: req.Company.Employee,
		Registered:    req.Company.Registered,
		Type:          gen.Type_name[int32(req.Company.Type)],
	}

	id, err := h.ctrl.CreateCompany(ctx, &company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.CreateCompanyResponse{Id: id}, nil
}

func (h *Handler) PatchCompany(ctx context.Context, req *gen.PatchCompanyRequest) (*gen.PatchCompanyResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	oldCompany := model.Company{
		ID:            req.Company.Id,
		Name:          req.Company.Name,
		Description:   req.Company.Description,
		EmployeeCount: req.Company.Employee,
		Registered:    req.Company.Registered,
		Type:          gen.Type_name[int32(req.Company.Type)],
	}

	newCompany, err := h.ctrl.PatchCompany(ctx, &oldCompany)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PatchCompanyResponse{
		Company: &gen.Company{
			Id:          newCompany.ID,
			Name:        newCompany.Name,
			Description: newCompany.Description,
			Employee:    newCompany.EmployeeCount,
			Registered:  newCompany.Registered,
			Type:        gen.Type(gen.Type_value[newCompany.Type]),
		},
	}, nil
}

func (h *Handler) DeleteCompany(ctx context.Context, req *gen.DeleteCompanyRequest) (*gen.DeleteCompanyResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	id, err := h.ctrl.DeleteCompany(ctx, req.Id)
	if err != nil && errors.Is(err, company.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.DeleteCompanyResponse{Id: id}, nil
}
