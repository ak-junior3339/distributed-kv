// centralized server connecting gRPC client

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/iips-oss/distributed-kv/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKvstoreClient(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			fmt.Printf("QUITTING\n")
			return
		}
		cmd := strings.Split(line, " ")
		if len(cmd) < 2 {
			fmt.Printf("Error: invalid command syntax (too few arguments)\n")
			continue
		}
		method := cmd[0]
		key := cmd[1]

		switch method {
		case "GET":
			if len(cmd) != 2 {
				fmt.Printf("Error: GET command requires exactly 1 argument. Usage: GET <key>\n")
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			r, err := c.KvGet(ctx, &pb.OpKeyReq{Key: key})
			cancel()
			if err != nil {
				fmt.Printf("RPC Error (GET): %v\n", err)
				continue
			}
			value := r.GetValue()
			if value != "" {
				fmt.Printf("%s\n", r.GetValue())
			}
		case "SET":
			if len(cmd) < 3 {
				fmt.Printf("Error: SET command requires a value. Usage: SET <key> <value>\n")
				continue
			}
			value := strings.Join(cmd[2:], " ")
			// value := cmd[2]
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := c.KvSet(ctx, &pb.SetReq{Key: key, Value: value})
			cancel()
			if err != nil {
				fmt.Printf("RPC Error (SET): %v\n", err)
				continue
			}
		case "DEL":
			if len(cmd) != 2 {
				fmt.Printf("Error: DEL command requires exactly 1 argument. Usage: DEL <key>\n")
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := c.KvDel(ctx, &pb.OpKeyReq{Key: key})
			cancel()
			if err != nil {
				fmt.Printf("RPC Error (DEL): %v\n", err)
				continue
			}
		case "exit": // maybe handle ctrl+c/d singles too
			fmt.Printf("QUITING\n")
			return
		default:
			fmt.Printf("Error: unknown command %q. Use GET, SET, DEL, or exit.\n", method)
			continue
		}
	}
}
