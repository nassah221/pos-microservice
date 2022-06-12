package service

import (
	"context"
	pb "pos-microservices/cashier/contract"
	"pos-microservices/cashier/model"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Set struct {
	CreateEndpoint     endpoint.Endpoint
	GetByIDEndpoint    endpoint.Endpoint
	GetByEmailEndpoint endpoint.Endpoint
	GetAllEndpoint     endpoint.Endpoint
	UpdateEndpoint     endpoint.Endpoint
	DeleteEndpoint     endpoint.Endpoint
	SigninEndpoint     endpoint.Endpoint
}

func authenticatedEndpoint(svc Service, next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		headers, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.New(codes.Unauthenticated, "missing headers").Err()
		}

		tokens := headers["authorization"]

		if len(tokens) != 1 {
			return nil, status.New(codes.Unauthenticated, "missing token").Err()
		}

		tokenString := tokens[0]

		token, err := svc.VerifyToken(tokenString)
		if err != nil {
			return nil, status.New(codes.Unauthenticated, "invalid token").Err()
		}

		ctx = context.WithValue(ctx, "token", token)

		return next(ctx, request)
	}
}

// func issueNewToken(svc Service, next endpoint.Endpoint) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		req := request.(*pb.SigninRequest)

// 		token, err := svc.IssueToken(req., req.Id, req.Duration)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return pb.SigninResponse{Token: token}, nil
// 	}
// }

func NewEndpointSet(s Service) Set {
	return Set{
		CreateEndpoint:     MakeCreateEndpoint(s),
		GetByIDEndpoint:    authenticatedEndpoint(s, MakeGetByIDEndpoint(s)),
		GetByEmailEndpoint: MakeGetByEmailEndpoint(s),
		GetAllEndpoint:     MakeGetAllEndpoint(s),
		UpdateEndpoint:     MakeUpdateEndpoint(s),
		DeleteEndpoint:     MakeDeleteEndpoint(s),
		SigninEndpoint:     MakeSigninEndpoint(s),
	}
}

func ImplEndpoints(s Service) Set {
	return Set{
		CreateEndpoint:  MakeCreateEndpoint(s),
		GetByIDEndpoint: authenticatedEndpoint(s, MakeGetByIDEndpoint(s)),
		SigninEndpoint:  MakeSigninEndpoint(s),
	}
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*model.Cashier)

		id, err := s.Signup(ctx, req)
		if err != nil {
			return nil, err
		}

		return id, nil
	}
}

func MakeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetCashierRequest)

		cashier, err := s.GetByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return cashier, nil
	}
}

func MakeGetByEmailEndpoint(s Service) endpoint.Endpoint {
	panic("unimplemented")
}

func MakeGetAllEndpoint(s Service) endpoint.Endpoint {
	panic("unimplemented")
}

func MakeUpdateEndpoint(s Service) endpoint.Endpoint {
	panic("unimplemented")
}

func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	panic("unimplemented")
}

func MakeSigninEndpoint(s Service) endpoint.Endpoint {
	// TODO: validate password and email
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SigninRequest)

		tokenString, err := s.IssueToken(req.Email, 0)
		if err != nil {
			return nil, err
		}

		return tokenString, nil
	}
}
