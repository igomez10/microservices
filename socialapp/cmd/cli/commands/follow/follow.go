package follow

import (
	followuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/follow/user"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "follow",
		Usage: "follow commands",
		Commands: []*cli.Command{
			followuser.GetCmd(),
		},
	}
}
