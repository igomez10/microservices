package listusers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "list users",
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
			&cli.IntFlag{
				Name:    cliflags.PageSizeFlag,
				Aliases: []string{"ps"},
				Value:   10,
			},
			&cli.IntFlag{
				Name:    cliflags.OffsetFlag,
				Aliases: []string{"o"},
				Value:   0,
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

			pagesize := c.Int(cliflags.PageSizeFlag)
			offset := c.Int(cliflags.OffsetFlag)

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.users.list"})
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

			apiClient := client.NewAPIClient(configuration)

			us, r, err := apiClient.UserAPI.ListUsers(ctx).
				Limit(int32(pagesize)).
				Offset(int32(offset)).
				Execute()
			if err != nil {
				return fmt.Errorf("Error when calling `UserAPI.ListUsers`: %v\n %+v\n", err, r)
			}
			if r.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", r.StatusCode)
			}

			for i := range us {
				fmt.Println(us[i].Username, us[i].Email, us[i].FirstName, us[i].LastName)
			}
			return nil
		},
	}
}
