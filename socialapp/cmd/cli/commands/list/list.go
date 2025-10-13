package list

import (
	listroles "github.com/igomez10/microservices/socialapp/cmd/cli/commands/list/roles"
	listscopes "github.com/igomez10/microservices/socialapp/cmd/cli/commands/list/scopes"
	listusers "github.com/igomez10/microservices/socialapp/cmd/cli/commands/list/users"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list commands",
		Commands: []*cli.Command{
			listusers.GetCmd(),
			listroles.GetCmd(),
			listscopes.GetCmd(),
		},
	}
}
