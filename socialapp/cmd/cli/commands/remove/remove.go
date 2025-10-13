package remove

import (
	removescopefromrole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/remove/scope-from-role"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove commands",
		Commands: []*cli.Command{
			removescopefromrole.GetCmd(),
		},
	}
}
