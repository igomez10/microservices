package update

import (
	updaterole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/update/role"
	updatescope "github.com/igomez10/microservices/socialapp/cmd/cli/commands/update/scope"
	updateuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/update/user"
	updateuserroles "github.com/igomez10/microservices/socialapp/cmd/cli/commands/update/user-roles"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "update commands",
		Commands: []*cli.Command{
			updateuser.GetCmd(),
			updaterole.GetCmd(),
			updatescope.GetCmd(),
			updateuserroles.GetCmd(),
		},
	}
}
