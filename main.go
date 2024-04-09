package main

import (
	"binance-cli/commands"
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	var Commands []*cli.Command
	Commands = append(Commands, commands.AssetQuery, commands.AssetWithdraw)
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Commands:              Commands,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
