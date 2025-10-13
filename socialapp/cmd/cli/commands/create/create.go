package create

import (
	createcomment "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/comment"
	createrole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/role"
	createscope "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/scope"
	createurl "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/url"
	createuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/create/user"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create commands",
		Commands: []*cli.Command{
			createuser.GetCmd(),
			createrole.GetCmd(),
			createscope.GetCmd(),
			createcomment.GetCmd(),
			createurl.GetCmd(),
		},
	}
}
