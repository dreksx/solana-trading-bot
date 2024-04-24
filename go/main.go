package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func main() {
	wsURL := os.Getenv("WS_URL")
	if wsURL == "" {
		panic("WS_URL env var is required")
	}

	fmt.Println(wsURL)
	client, err := ws.Connect(context.Background(), wsURL)
	if err != nil {
		panic(err)
	}
	//OPENBOOK_MARKET := solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")

	program := solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8") // serum
	{

		filters := []rpc.RPCFilter{}
		sub, err := client.ProgramSubscribeWithOpts(program, "processed", "", filters)

		//sub, err := client.AccountSubscribe(
		//	program,
		//	"processed",
		//)
		if err != nil {
			panic(err)
		}
		defer sub.Unsubscribe()

		for {
			got, err := sub.Recv()
			if err != nil {
				panic(err)
			}

			spew.Dump(got)
		}
	}
}
