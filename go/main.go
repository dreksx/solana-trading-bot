package main

import (
	"context"
	"fmt"
	"os"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type LiquidityStateV4 struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       [128]byte
	SwapQuoteOutAmount     [128]byte
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      [128]byte
	SwapBaseOutAmount      [128]byte
	SwapQuote2BaseFee      uint64
	BaseVault              solana.PublicKey
	QuoteVault             solana.PublicKey
	BaseMint               solana.PublicKey // 32 bytes
	QuoteMint              solana.PublicKey // 32 bytes
	LpMint                 solana.PublicKey // 32 bytes
	OpenOrders             solana.PublicKey // 32 bytes
	MarketId               solana.PublicKey // 32 bytes
	MarketProgramId        solana.PublicKey // 32 bytes
	TargetOrders           solana.PublicKey // 32 bytes
	WithdrawQueue          solana.PublicKey // 32 bytes
	LpVault                solana.PublicKey // 32 bytes
	Owner                  solana.PublicKey // 32 bytes
	LpReserve              uint64
	Padding                [3]uint64 // To handle the seq(u64(), 3, 'padding')
}

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

			var mint LiquidityStateV4
			err = bin.NewBorshDecoder(got.Value.Account.Data.GetBinary()).Decode(&mint)
			if err != nil {
				panic(err)
			}

			fmt.Println(mint.PoolOpenTime, mint.BaseMint.String(), mint.QuoteMint.String())
			os.Exit(1)

			//mint := bytes[baseMintOffset:baseMintOffsetEnd]
			//fmt.Println(string(poolOpenTime), string(mint))
		}
	}
}
