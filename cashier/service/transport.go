package service

import (
	"context"
	pb "pos-microservices/cashier/contract"
	"pos-microservices/cashier/model"

	gt "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	create gt.Handler
	login  gt.Handler
	get    gt.Handler
	pb.UnimplementedCashierServiceServer
}

func NewGRPCServer(_ context.Context, svc Service) pb.CashierServiceServer {
	set := ImplEndpoints(svc)
	createHandler := gt.NewServer(set.CreateEndpoint, decodeCreateCashierRequest, encodeCreateCashierResponse)
	signinHandler := gt.NewServer(set.SigninEndpoint, decodeSigninRequest, encodeSigninResponse)
	getHandler := gt.NewServer(set.GetByIDEndpoint, decodeGetRequest, encodeGetResponse)

	return &grpcServer{
		create: createHandler,
		login:  signinHandler,
		get:    getHandler,
	}
}

func (s *grpcServer) Signup(ctx context.Context, req *pb.Cashier) (*pb.GetCashierRequest, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetCashierRequest), nil
}

func (s *grpcServer) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.SigninResponse), nil
}

func (s *grpcServer) GetCashier(ctx context.Context, req *pb.GetCashierRequest) (*pb.Cashier, error) {
	_, rep, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.Cashier), nil
}

func (s *grpcServer) ListCashiers(context.Context, *pb.ListCashiersRequest) (*pb.ListCashiersResponse, error) {
	panic("unimplemented")
}

func (s *grpcServer) UpdateCashier(context.Context, *pb.Cashier) (*pb.Cashier, error) {
	panic("unimplemented")
}

func (s *grpcServer) DeleteCashier(context.Context, *pb.GetCashierRequest) (*pb.DeleteCashierResponse, error) {
	panic("unimplemented")
}

func decodeGetRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*model.Cashier)

	return &pb.Cashier{
		Id:       r.ID.Hex(),
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Created:  r.Created,
		Updated:  r.Updated,
	}, nil
}

func encodeSigninResponse(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(string)

	return &pb.SigninResponse{
		Token: r,
	}, nil
}

func decodeSigninRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeCreateCashierResponse(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(string)

	return &pb.GetCashierRequest{
		Id: r,
	}, nil
}

func decodeCreateCashierRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.Cashier)

	return &model.Cashier{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}, nil
}
