package get

import (
	geturl "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/url"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get commands",
		Commands: []*cli.Command{
			geturl.GetCmd(),
		},
	}
}
