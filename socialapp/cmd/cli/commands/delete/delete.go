package delete

import (
	deleterole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/delete/role"
	deletescope "github.com/igomez10/microservices/socialapp/cmd/cli/commands/delete/scope"
	deleteurl "github.com/igomez10/microservices/socialapp/cmd/cli/commands/delete/url"
	deleteuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/delete/user"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete commands",
		Commands: []*cli.Command{
			deleteuser.GetCmd(),
			deleterole.GetCmd(),
			deletescope.GetCmd(),
			deleteurl.GetCmd(),
		},
	}
}
