package add

import (
	addscopetorole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/add/scope-to-role"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add commands",
		Commands: []*cli.Command{
			addscopetorole.GetCmd(),
		},
	}
}
