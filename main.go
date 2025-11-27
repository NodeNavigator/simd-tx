package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	// "github.com/cosmos/evm/crypto/ethsecp256k1"
	"github.com/evmos/evmos/v20/crypto/ethsecp256k1"
	"github.com/jaybxyz/simd-tx/client"
	"github.com/jaybxyz/simd-tx/codec"
	"github.com/jaybxyz/simd-tx/config"
	"github.com/jaybxyz/simd-tx/wallet"
)

/*
TODO
1. Remove config.toml as this project is just sample code and it increases line of code to write
2. Do we need rpcURL anymore? Can we use grpc to query network info?
3. Research Cosmos SDK's written interfaces to see if this code can be shortend in any way or create a simple library to make the process simpler

*/

var (
	timeout = 5 * time.Second
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // human-friendly output
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	config, err := config.Read(config.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config.toml file: %w", err))
	}

	// Configure SDK Bech32 prefixes and coin type for the network.
	// Use `investnet` as the address prefix instead of the default `cosmos`.
	// Place this before any address encoding/decoding or calls that depend on the SDK config.
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("investnet", "investnetpub")
	cfg.SetBech32PrefixForValidator("investnetvaloper", "investnetvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("investnetvalcons", "investnetvalconspub")
	// If you want to use the Ethereum coin type (BIP44 coin type 60), set it here.
	// Otherwise change to your network's coin type.
	cfg.SetCoinType(60)
	cfg.Seal()
	// Connect Tendermint RPC client
	rpcClient, err := client.ConnectRPCWithTimeout(config.RPC.Address, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to connect RPC client: %w", err))
	}

	// Connect gRPC client
	gRPCConn, err := client.ConnectGRPCWithTimeout(ctx, config.GRPC.Address, config.GRPC.UseTLS, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to connect gRPC client: %w", err))
	}
	defer gRPCConn.Close()

	// Recover private key from mnemonic phrases
	privKeyEth, err := wallet.RecoverPrivKeyFromMnemonic(config.WalletConfig.Mnemonic, config.WalletConfig.Password)
	if err != nil {
		panic(fmt.Errorf("recovering private key: %w", err))
	}

	// Convert ethsecp256k1 private key to secp256k1 private key
	privKey := &ethsecp256k1.PrivKey{Key: privKeyEth.Bytes()}

	fmt.Println("privKey:", privKey.String())
	chainID, _ := rpcClient.NetworkChainID(ctx)
	creator := wallet.Address(privKey)
	baseAccount, _ := gRPCConn.GetAccount(ctx, creator.String())
	accNum := baseAccount.GetAccountNumber()
	accSeq := baseAccount.GetSequence()
	gasLimit := config.TxConfig.GasLimit
	fees, err := sdk.ParseCoinsNormalized(config.TxConfig.Fees)
	fmt.Println("chainID:", chainID)
	fmt.Println("creator:", creator.String())
	fmt.Println("baseAccount:", baseAccount)
	fmt.Println("accNum:", accNum)
	fmt.Println("accSeq:", accSeq)
	fmt.Println("gasLimit:", gasLimit)
	fmt.Println("fees:", fees)
	if err != nil {
		panic(fmt.Errorf("failed to parse coins %w", err))
	}

	// Create new MsgSend for test
	msg := banktypes.MsgSend{
		FromAddress: creator.String(),
		ToAddress:   "investnet1fx944mzagwdhx0wz7k9tfztc8g3lkfk6w756rn",
		// Token has 18 decimals. To send 10 invst, multiply by 10^18.
		Amount: sdk.NewCoins(sdk.NewCoin("invst", sdkmath.NewIntWithDecimal(10, 18))),
	}
	fmt.Println("msg======", msg)
	msgs := []sdk.Msg{&msg}

	tx := client.NewTx(
		chainID,
		accNum,
		accSeq,
		gasLimit,
		fees,
		msgs...,
	)

	txCfg := codec.MakeEncodingConfig().TxConfig
	txBytes, err := client.SignTx(tx, txCfg, privKey)
	if err != nil {
		fmt.Printf("failed to sign transaction: %v", err)
		return
	}

	fmt.Println("txBytes====", txBytes)
	// Use BLOCK mode while debugging so this call waits for the tx to be included in a block
	// and returns the chain response. Switch back to SYNC or ASYNC for non-blocking behavior.
	resp, err := gRPCConn.BroadcastTx(ctx, txBytes, sdktx.BroadcastMode_BROADCAST_MODE_SYNC)
	if err != nil {
		fmt.Printf("failed to broadcast transaction: %v\n", err)
		if resp != nil && resp.TxResponse != nil {
			fmt.Printf("broadcast response (error): code=%d raw_log=%s txhash=%s\n", resp.TxResponse.Code, resp.TxResponse.RawLog, resp.TxResponse.TxHash)
		}
		return
	}

	if resp != nil && resp.TxResponse != nil {
		tr := resp.TxResponse
		fmt.Printf("broadcast result: code=%d txhash=%s gasWanted=%d gasUsed=%d raw_log=%s\n", tr.Code, tr.TxHash, tr.GasWanted, tr.GasUsed, tr.RawLog)
		if tr.Code != 0 {
			fmt.Printf("transaction failed on-chain with code=%d raw_log=%s\n", tr.Code, tr.RawLog)
			return
		}
	}

	log.Info().Msg("Go to the following link to see if transaction is successfully included in a block")
	if resp != nil && resp.TxResponse != nil {
		log.Info().Msg("http://localhost:1317/cosmos/tx/v1beta1/txs/" + resp.TxResponse.TxHash)
	}
}
