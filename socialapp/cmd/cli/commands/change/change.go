package change

import (
	changepassword "github.com/igomez10/microservices/socialapp/cmd/cli/commands/change/password"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "change",
		Usage: "change commands",
		Commands: []*cli.Command{
			changepassword.GetCmd(),
		},
	}
}
