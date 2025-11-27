package wallet

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/evm/crypto/ethsecp256k1"
	// evmhd "github.com/cosmos/evm/crypto/hd"
	"github.com/evmos/evmos/v20/crypto/ethsecp256k1"
	evmhd "github.com/evmos/evmos/v20/crypto/hd"
)

func RecoverPrivKeyFromMnemonic(mnemonic, password string) (*ethsecp256k1.PrivKey, error) {
	// Use EthSecp256k1 algorithm to derive private key from mnemonic
	deriveFn := evmhd.EthSecp256k1.Derive()

	// Get the BIP44 path for Ethereum (coin type 60)
	path := sdk.GetConfig().GetFullBIP44Path()
	fmt.Println("path=", path)

	privKeyBytes, err := deriveFn(mnemonic, password, path)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	// Generate the private key from derived bytes
	generateFn := evmhd.EthSecp256k1.Generate()
	privKey := generateFn(privKeyBytes)

	return privKey.(*ethsecp256k1.PrivKey), nil
}

func Address(privKey *ethsecp256k1.PrivKey) sdk.AccAddress {
	return sdk.AccAddress(privKey.PubKey().Address())
}
