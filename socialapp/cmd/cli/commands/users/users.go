package users

import (
	"github.com/igomez10/microservices/socialapp/cmd/cli/commands/users/list"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Users commands",
		Commands: []*cli.Command{
			list.GetCmd(),
		},
	}
}
