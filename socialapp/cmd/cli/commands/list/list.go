package list

import (
	listusers "github.com/igomez10/microservices/socialapp/cmd/cli/commands/list/users"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list commands",
		Commands: []*cli.Command{
			listusers.GetCmd(),
		},
	}
}
