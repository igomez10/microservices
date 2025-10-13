package get

import (
	getcomment "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/comment"
	getfeed "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/feed"
	getrole "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/role"
	getscope "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/scope"
	geturl "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/url"
	getuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/user"
	getusercomments "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/user-comments"
	getuserfollowers "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/user-followers"
	getuserfollowing "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/user-following"
	getuserroles "github.com/igomez10/microservices/socialapp/cmd/cli/commands/get/user-roles"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get commands",
		Commands: []*cli.Command{
			geturl.GetCmd(),
			getuser.GetCmd(),
			getusercomments.GetCmd(),
			getfeed.GetCmd(),
			getuserfollowers.GetCmd(),
			getuserfollowing.GetCmd(),
			getcomment.GetCmd(),
			getuserroles.GetCmd(),
			getscope.GetCmd(),
			getrole.GetCmd(),
		},
	}
}
