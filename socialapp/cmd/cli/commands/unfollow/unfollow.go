package unfollow

import (
	unfollowuser "github.com/igomez10/microservices/socialapp/cmd/cli/commands/unfollow/user"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "unfollow",
		Usage: "unfollow commands",
		Commands: []*cli.Command{
			unfollowuser.GetCmd(),
		},
	}
}
