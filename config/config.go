package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultConfigPath = "config.toml"
	DefaultGasLimit   = uint64(500000)
	DefaultFees       = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))
)

// Config defines all necessary configuration parameters.
type Config struct {
	RPC          RPCConfig    `toml:"rpc"`
	GRPC         GRPCConfig   `toml:"grpc"`
	WalletConfig WalletConfig `toml:"wallet"`
	TxConfig     TxConfig     `toml:"tx"`
}

// RPCConfig contains configuration of the RPC endpoint.
type RPCConfig struct {
	Address string `toml:"address"`
}

// GRPCConfig contains configuration of the gRPC endpoint.
type GRPCConfig struct {
	Address string `toml:"address"`
	UseTLS  bool   `toml:"use_tls"`
}

// WalletConfig contains wallet configuration that is used to sign transaction.
type WalletConfig struct {
	Mnemonic string `yaml:"mnemonic"`
	Password string `yaml:"password"`
}

// TxConfig contains configuration for transaction related parameters.
type TxConfig struct {
	GasLimit uint64 `yaml:"gas_limit"`
	Fees     string `yaml:"fees"`
}

// NewConfig builds a new Config instance.
func NewConfig(rpcCfg RPCConfig, gRPCCfg GRPCConfig) *Config {
	return &Config{
		RPC:  rpcCfg,
		GRPC: gRPCCfg,
	}
}

// SetupConfig takes the path to a configuration file and returns the properly parsed configuration.
func Read(configPath string) (*Config, error) {
	log.Debug().Msg("reading config file...")

	// Use default config path with a file name "config.toml" if it is empty
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return ParseString(data)
}

// ParseString attempts to read and parse  config from the given string bytes.
// An error reading or parsing the config results in a panic.
func ParseString(configData []byte) (*Config, error) {
	log.Debug().Msg("parsing config data...")

	var cfg Config
	err := toml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.TxConfig.GasLimit == 0 {
		cfg.TxConfig.GasLimit = DefaultGasLimit
	}
	if cfg.TxConfig.Fees == "" {
		cfg.TxConfig.Fees = DefaultFees.String()
	}

	return &cfg, nil
}
