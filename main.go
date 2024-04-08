package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	bc "github.com/binance/binance-connector-go"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
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
			},
			{
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
				Action: send,
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Private Private
}
type Private struct {
	ApiKey    string
	SecretKey string
	BaseURL   string
}
type SendMessage struct {
	Address string
	Amount  float64
}

func (config *Config) readConfigToml() error {
	if _, err := toml.DecodeFile("config.toml", config); err != nil {
		return err
	}
	return nil
}
func readCsv(path string) (list []SendMessage, err error) {
	file, err := os.Open(path)
	if err != nil {
		return []SendMessage{}, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) != 2 {
			return []SendMessage{}, err
		}
		f, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return []SendMessage{}, err
		}
		message := SendMessage{
			Address: record[0],
			Amount:  f,
		}
		list = append(list, message)
	}
	return
}
func query(ctx context.Context, cmd *cli.Command) error {
	var cfg = new(Config)
	err := cfg.readConfigToml()
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
func send(ctx context.Context, cmd *cli.Command) error {
	var cfg = new(Config)
	err := cfg.readConfigToml()
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
