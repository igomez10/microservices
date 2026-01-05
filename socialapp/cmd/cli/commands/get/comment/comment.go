package getcomment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "comment",
		Usage:     "get comment by id",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.EnvFlag,
				Usage: "Environment to use (live, local)",
				Value: auth.DefaultEnv,
			},
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Usage:   "Username for authentication (or set SOCIALAPP_CLI_USERNAME)",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
				Usage:   "Password for authentication (or set SOCIALAPP_CLI_PASSWORD)",
				Aliases: []string{"p"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			envName := auth.ResolveEnvironment(c.String(cliflags.EnvFlag))

			username, password, err := auth.ResolveCredentials(
				c.String(cliflags.UsernameFlag),
				c.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			if c.Args().Len() < 1 {
				return fmt.Errorf("comment id argument is required")
			}
			commentIDStr := c.Args().Get(0)
			commentID, err := strconv.ParseInt(commentIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid comment id: %v", err)
			}

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.comments.read"})
			if err != nil {
				return fmt.Errorf("failed to get authenticated client: %w", err)
			}

			host, scheme, err := auth.GetAPIClientConfig(envName)
			if err != nil {
				return err
			}

			configuration := client.NewConfiguration()
			configuration.Host = host
			configuration.Scheme = scheme
			configuration.HTTPClient = httpClient

			clnt := client.NewAPIClient(configuration)

			comment, res, err := clnt.CommentAPI.GetComment(ctx, commentID).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v\n", res)
				return fmt.Errorf("error when calling `CommentAPI.GetComment`: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("expected status code 200, got %d", res.StatusCode)
			}

			b, err := json.MarshalIndent(comment, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}

			fmt.Println(string(b))
			return nil
		},
	}
}
