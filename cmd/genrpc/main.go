package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vmkteam/rpcgen/v2"
	"github.com/vmkteam/rpcgen/v2/golang"
	"github.com/vmkteam/zenrpc/v2"

	"github.com/yaBliznyk/newsportal/internal/rpc"
)

func main() {
	rpcServer := zenrpc.NewServer(zenrpc.Options{})
	rpcServer.Register("news", rpc.NewsService{})

	generated, err := rpcgen.FromSMD(rpcServer.SMD()).
		GoClient(golang.Settings{
			Package:         "newsportal",
			CallerNamespace: "news",
		}).
		Generate()
	if err != nil {
		log.Fatal(err)
	}

	dst := "tests/newsportal/newsportal.go"
	if err := os.WriteFile(dst, generated, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Client generated: %s\n", dst)
}
