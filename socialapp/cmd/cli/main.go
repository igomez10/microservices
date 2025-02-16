package main

import (
	"context"
	"fmt"
	"os"

	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/users"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:                  "socialapp-cli",
		Usage:                 "A CLI for the socialapp service",
		UsageText:             "socialapp-cli [global options] command [command options] [arguments...]",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			users.GetCmd(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
