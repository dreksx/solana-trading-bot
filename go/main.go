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

	quoteMint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	openBook := solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
	program := solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8") // serum
	{
		var quoteMintOffset uint64 = 432
		var openBookOffset uint64 = 560
		var statusOffset uint64 = 0
		// Assuming 'marketProgramId' follows 'marketId'
		//marketProgramIDOffset := 864 + 32 + 32 // Adjust according to actual layout

		filters := []rpc.RPCFilter{

			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: quoteMintOffset,
					Bytes:  solana.Base58(quoteMint.Bytes()),
				},
			},
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: openBookOffset, // replace with correct offset for 'marketProgramId'
					Bytes:  solana.Base58(openBook.Bytes()),
				},
			},
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: statusOffset, // replace with correct offset for 'status'
					Bytes:  solana.Base58([]byte{6, 0, 0, 0, 0, 0, 0, 0}),
				},
			},
		}
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
