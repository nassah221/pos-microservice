package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pos-microservices/cashier/auth"
	pb "pos-microservices/cashier/contract"
	db "pos-microservices/cashier/mongo"
	"pos-microservices/cashier/repository"
	"pos-microservices/cashier/service"
	"syscall"

	"google.golang.org/grpc"
)

var (
	gRPCAddr = flag.String("grpc", ":8081",
		"gRPC listen address")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	errChan := make(chan error)

	cfg, err := db.NewConfig("../.env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s, err := db.NewStore(context.Background(), cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	r := repository.NewRepository(s)

	a := auth.NewAuthenticator("secret")

	svc := service.NewService(r, a)

	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		server := service.NewGRPCServer(ctx, svc)
		grpcServer := grpc.NewServer()
		pb.RegisterCashierServiceServer(grpcServer, server)
		errChan <- grpcServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println("Listening on: ", *gRPCAddr)
	fmt.Println(<-errChan)
}
