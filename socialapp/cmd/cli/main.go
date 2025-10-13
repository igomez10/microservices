package main

import (
	"context"
	"fmt"
	"os"

	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/add"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/change"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/create"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/delete"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/follow"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/get"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/list"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/remove"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/unfollow"
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/update"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:                  "socialapp-cli",
		Usage:                 "A CLI for the socialapp service",
		UsageText:             "socialapp-cli [global options] command [command options] [arguments...]",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			list.GetCmd(),
			get.GetCmd(),
			create.GetCmd(),
			update.GetCmd(),
			delete.GetCmd(),
			follow.GetCmd(),
			unfollow.GetCmd(),
			add.GetCmd(),
			remove.GetCmd(),
			change.GetCmd(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
