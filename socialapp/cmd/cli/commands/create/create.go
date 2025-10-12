package create

import (
	createuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/user"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create commands",
		Commands: []*cli.Command{
			createuser.GetCmd(),
		},
	}
}
