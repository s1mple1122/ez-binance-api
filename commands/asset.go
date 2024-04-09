package commands

import (
	"binance-cli/config"
	"context"
	"errors"
	"fmt"
	bc "github.com/binance/binance-connector-go"
	"github.com/urfave/cli/v3"
	"strings"
	"time"
)

var AssetQuery = &cli.Command{
	Name:    "query",
	Aliases: []string{"q"},
	Usage:   "query api-key Asset and symbol network",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "symbol",
			Aliases:  []string{"s"},
			Usage:    "symbol(required)",
			Required: true,
		},
	},
	Action: query,
}
var AssetWithdraw = &cli.Command{
	Name:    "withdraw",
	Aliases: []string{"w"},
	Usage:   "withdraw",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "symbol",
			Aliases:  []string{"s"},
			Usage:    "symbol(required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network",
			Aliases:  []string{"n"},
			Usage:    "network(required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "path",
			Aliases:  []string{"p"},
			Usage:    "filepath,enable csv(required)",
			Required: true,
		},
	},
	Action: withdraw,
}

func query(ctx context.Context, cmd *cli.Command) error {
	var cfg = new(config.Config)
	err := cfg.ReadConfigToml()
	if err != nil {
		return err
	}
	symbol := cmd.String("symbol")
	client := bc.NewClient(cfg.Private.ApiKey, cfg.Private.SecretKey, cfg.Private.BaseURL)
	// FundingWalletService - /sapi/v1/asset/get-funding-asset
	fundingWallet, err := client.NewFundingWalletService().Asset(symbol).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("balance:")
	fmt.Println(bc.PrettyPrint(fundingWallet))
	allCoinsInfo, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		return err
	}
	var network = make(map[string]string)
	for _, v := range allCoinsInfo {
		if strings.EqualFold(v.Coin, symbol) {
			for _, net := range v.NetworkList {
				network[net.Network] = net.Name

			}
		}
	}
	fmt.Println("enable network:")
	fmt.Println(bc.PrettyPrint(network))
	return nil
}
func withdraw(ctx context.Context, cmd *cli.Command) error {
	var cfg = new(config.Config)
	err := cfg.ReadConfigToml()
	if err != nil {
		return err
	}
	network := cmd.String("network")
	symbol := cmd.String("symbol")
	path := cmd.String("path")
	list, err := readCsv(path)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New("无转账明细")
	}
	fmt.Println("send message:")
	client := bc.NewClient(cfg.Private.ApiKey, cfg.Private.SecretKey, cfg.Private.BaseURL)
	for i, v := range list {
		withdraw, err := client.NewWithdrawService().Coin(strings.ToUpper(symbol)).Address(v.Address).
			Amount(v.Amount).Network(network).Do(context.Background())
		if err != nil {
			return err
		}
		fmt.Println(bc.PrettyPrint(withdraw))
		if i < len(list)-1 {
			time.Sleep(10 * time.Second)
		} else {
			return nil
		}
	}
	return nil
}
