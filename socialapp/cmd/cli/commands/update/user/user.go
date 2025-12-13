package updateuser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "user",
		Usage:     "update user by username",
		ArgsUsage: "<username>",
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
			&cli.StringFlag{
				Name:    cliflags.EmailFlag,
				Aliases: []string{"e"},
			},
			&cli.StringFlag{
				Name:    cliflags.FirstNameFlag,
				Aliases: []string{"f"},
			},
			&cli.StringFlag{
				Name:    cliflags.LastNameFlag,
				Aliases: []string{"l"},
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

			email := c.String(cliflags.EmailFlag)
			firstname := c.String(cliflags.FirstNameFlag)
			lastname := c.String(cliflags.LastNameFlag)

			if c.Args().Len() < 1 {
				return fmt.Errorf("username argument is required")
			}
			targetUsername := c.Args().Get(0)

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.users.update", "socialapp.users.read"})
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

			// First, get the existing user to preserve fields not being updated
			existingUser, _, err := clnt.UserAPI.GetUserByUsername(ctx, targetUsername).Execute()
			if err != nil {
				return fmt.Errorf("error getting existing user: %v", err)
			}

			// Build the user update request with provided fields, keeping existing values for others
			userUpdate := *existingUser
			if email != "" {
				userUpdate.Email = email
			}
			if firstname != "" {
				userUpdate.FirstName = firstname
			}
			if lastname != "" {
				userUpdate.LastName = lastname
			}

			user, res, err := clnt.UserAPI.UpdateUser(ctx, targetUsername).User(userUpdate).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v\n", res)
				return fmt.Errorf("error when calling `UserAPI.UpdateUser`: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("expected status code 200, got %d", res.StatusCode)
			}

			b, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}

			fmt.Println(string(b))
			return nil
		},
	}
}
