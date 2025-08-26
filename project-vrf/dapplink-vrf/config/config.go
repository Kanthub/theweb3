package config

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/dapplink-vrf/flags"
)

const (
	defaultConfirmations = 64
	defaultLoopInterval  = 5000
)

type Config struct {
	Migrations     string
	Chain          ChainConfig
	MasterDB       DBConfig
	SlaveDB        DBConfig
	SlaveDbEnable  bool
	ApiCacheEnable bool
}

type ChainConfig struct {
	ChainRpcUrl                       string
	ChainId                           uint
	StartingHeight                    uint64
	Confirmations                     uint64
	BlockStep                         uint64
	Contracts                         []common.Address
	MainLoopInterval                  time.Duration
	EventInterval                     time.Duration
	CallInterval                      time.Duration
	PrivateKey                        string
	DappLinkVrfContractAddress        string
	DappLinkVrfFactoryContractAddress string
	CallerAddress                     string
	NumConfirmations                  uint64
	SafeAbortNonceTooLowCount         uint64
	Mnemonic                          string
	CallerHDPath                      string
	Passphrase                        string
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func LoadConfig(cliCtx *cli.Context) (Config, error) {
	var cfg Config
	cfg = NewConfig(cliCtx)

	if cfg.Chain.Confirmations == 0 {
		cfg.Chain.Confirmations = defaultConfirmations
	}

	if cfg.Chain.MainLoopInterval == 0 {
		cfg.Chain.MainLoopInterval = defaultLoopInterval
	}

	log.Info("loaded chain config", "config", cfg.Chain)
	return cfg, nil
}

func LoadContracts() []common.Address {
	var Contracts []common.Address
	Contracts = append(Contracts, DappLinkVrfAddr)
	return Contracts
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		Migrations: ctx.String(flags.MigrationsFlag.Name),
		Chain: ChainConfig{
			ChainId:                           ctx.Uint(flags.ChainIdFlag.Name),
			ChainRpcUrl:                       ctx.String(flags.ChainRpcFlag.Name),
			StartingHeight:                    ctx.Uint64(flags.StartingHeightFlag.Name),
			Confirmations:                     ctx.Uint64(flags.ConfirmationsFlag.Name),
			BlockStep:                         ctx.Uint64(flags.BlocksStepFlag.Name),
			Contracts:                         LoadContracts(),
			MainLoopInterval:                  ctx.Duration(flags.MainIntervalFlag.Name),
			EventInterval:                     ctx.Duration(flags.EventIntervalFlag.Name),
			CallInterval:                      ctx.Duration(flags.CallIntervalFlag.Name),
			PrivateKey:                        ctx.String(flags.PrivateKeyFlag.Name),
			DappLinkVrfContractAddress:        ctx.String(flags.DappLinkVrfContractAddressFlag.Name),
			DappLinkVrfFactoryContractAddress: ctx.String(flags.DappLinkVrfFactoryContractAddressFlag.Name),
			CallerAddress:                     ctx.String(flags.CallerAddressFlag.Name),
			NumConfirmations:                  ctx.Uint64(flags.NumConfirmationsFlag.Name),
			SafeAbortNonceTooLowCount:         ctx.Uint64(flags.SafeAbortNonceTooLowCountFlag.Name),
			Mnemonic:                          ctx.String(flags.MnemonicFlag.Name),
			CallerHDPath:                      ctx.String(flags.CallerHDPathFlag.Name),
			Passphrase:                        ctx.String(flags.PassphraseFlag.Name),
		},
		MasterDB: DBConfig{
			Host:     ctx.String(flags.MasterDbHostFlag.Name),
			Port:     ctx.Int(flags.MasterDbPortFlag.Name),
			Name:     ctx.String(flags.MasterDbNameFlag.Name),
			User:     ctx.String(flags.MasterDbUserFlag.Name),
			Password: ctx.String(flags.MasterDbPasswordFlag.Name),
		},
		SlaveDB: DBConfig{
			Host:     ctx.String(flags.SlaveDbHostFlag.Name),
			Port:     ctx.Int(flags.SlaveDbPortFlag.Name),
			Name:     ctx.String(flags.SlaveDbNameFlag.Name),
			User:     ctx.String(flags.SlaveDbUserFlag.Name),
			Password: ctx.String(flags.SlaveDbPasswordFlag.Name),
		},
		SlaveDbEnable: ctx.Bool(flags.SlaveDbEnableFlag.Name),
	}
}
