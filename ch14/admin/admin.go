package main

import (
	"ch14/grpcapi"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.AdminClient
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 4445), opts...); err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client = grpcapi.NewAdminClient(conn)
	for {
		var cmd = new(grpcapi.Command)
		fmt.Printf("[+] Please enter a command to be sent to the implant.\n")
		if _, err := fmt.Scanln(&cmd.In); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("[+] Sending command down the pipeline [%s]\n", cmd.In)
		ctx := context.Background()
		cmd, err = client.RunCommand(ctx, cmd)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("[+] Recv: [%s]\n", cmd.Out)
	}

	//fmt.Println(cmd.Out)
}
