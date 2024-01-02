package main

import (
	"ch14/grpcapi"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type implantServer struct {
	work, output chan *grpcapi.Command
	grpcapi.UnimplementedImplantServer
}

type adminServer struct {
	work, output chan *grpcapi.Command
	grpcapi.UnimplementedAdminServer
}

func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
	s := new(implantServer)
	s.work = work
	s.output = output
	return s
}

func NewAdminServer(work, output chan *grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s
}

func (s *implantServer) FetchCommand(ctx context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	var cmd = new(grpcapi.Command)
	select {
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("channel closed")
	default:
		//no work
		return cmd, nil
	}
}

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		s.work <- cmd // this blocks since it's un unbuffered channel. We need to use a go routine to send work down pipeline then receive the output by continuing execution.
	}()
	res = <-s.output // continued execution here. this blocks as well. Both need to happen for execution to continue.
	return res, nil
}

func main() {
	var (
		implantListener, adminListener net.Listener
		err                            error
		opts                           []grpc.ServerOption
		work, output                   chan *grpcapi.Command
	)
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output)

	fmt.Printf("[+] Starting implant server on localhost:4444\n")
	if implantListener, err = net.Listen("tcp",
		fmt.Sprintf("localhost:%d", 4444)); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] Starting admin server on localhost:4445\n")
	if adminListener, err = net.Listen("tcp",
		fmt.Sprintf("localhost:%d", 4445)); err != nil {
		log.Fatal(err)
	}

	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)

	go func() {
		grpcImplantServer.Serve(implantListener)
	}()

	grpcAdminServer.Serve(adminListener)
}
