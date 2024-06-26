package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"time"

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
	SwapBaseInAmount       [16]byte
	SwapQuoteOutAmount     [16]byte
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      [16]byte
	SwapBaseOutAmount      [16]byte
	SwapQuote2BaseFee      uint64
	BaseVault              solana.PublicKey
	QuoteVault             solana.PublicKey
	BaseMint               solana.PublicKey // 32 bytes
	QuoteMint              solana.PublicKey // 32 bytes
}

func main() {
	//ctx := context.Background()
	wsURL := os.Getenv("WS_URL")
	if wsURL == "" {
		panic("WS_URL env var is required")
	}

	client, err := ws.Connect(context.Background(), wsURL)
	if err != nil {
		panic(err)
	}
	//httpURL := strings.Replace(wsURL, "wss://", "https://", 1)
	//r := rpc.New(httpURL)
	//str := "2iD1qz74YRzWoBYSBctpdDdgRaWGS1d7E3FTtWMMWE6G3WKNhcDdjnNEcvG4EMrnwYdawtVkA1oeYmbQrKKue23g"
	//s := solana.MustSignatureFromBase58(str)
	//version := uint64(0)
	//n := time.Now()
	//_, err = r.GetTransaction(ctx, s, &rpc.GetTransactionOpts{
	//	MaxSupportedTransactionVersion: &version,
	//	Encoding:                       solana.EncodingBase64,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(time.Since(n))
	//os.Exit(1)
	//OPENBOOK_MARKET := solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")

	quoteMint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	openBook := solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
	cached := make(map[solana.PublicKey]struct{}, 100000)

	now := uint64(time.Now().Unix())

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

		openTimeOffset := 224
		mintOffset := 400
		for {
			got, err := sub.Recv()
			if err != nil {
				panic(err)
			}
			t := time.Now()
			bytes := got.Value.Account.Data.GetBinary()
			poolOpenTime := binary.LittleEndian.Uint64(bytes[openTimeOffset : openTimeOffset+8])
			mintValue := solana.PublicKeyFromBytes(bytes[mintOffset : mintOffset+32])
			f := time.Since(t)
			if poolOpenTime < now {
				continue
			}
			if _, exist := cached[mintValue]; exist {
				continue
			}
			fmt.Println(time.Now().Format("15:04:05.999999999"), "\t", mintValue.String(), "\t", f)

			cached[mintValue] = struct{}{}
		}
	}
}
