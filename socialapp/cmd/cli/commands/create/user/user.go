package createuser

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
		Name:  "user",
		Usage: "create user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.EnvFlag,
				Usage: "Environment to use (live, local)",
				Value: auth.DefaultEnv,
			},
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Usage:   "Username for the new user",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
				Usage:   "Password for the new user",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    cliflags.EmailFlag,
				Usage:   "Email for the new user",
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

			username := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)
			email := c.String(cliflags.EmailFlag)
			firstname := c.String(cliflags.FirstNameFlag)
			lastname := c.String(cliflags.LastNameFlag)

			if username == "" {
				return fmt.Errorf("username is required")
			}
			if password == "" {
				return fmt.Errorf("password is required")
			}
			if email == "" {
				return fmt.Errorf("email is required")
			}

			host, scheme, err := auth.GetAPIClientConfig(envName)
			if err != nil {
				return err
			}

			configuration := client.NewConfiguration()
			configuration.Host = host
			configuration.Scheme = scheme
			configuration.HTTPClient = http.DefaultClient

			clnt := client.NewAPIClient(configuration)

			req := client.CreateUserRequest{
				Username: username,
				Password: password,
				Email:    email,
				FirstName: func() string {
					if len(firstname) > 0 {
						return firstname
					}
					return "FirstName"
				}(),
				LastName: func() string {
					if len(lastname) > 0 {
						return lastname
					}
					return "LastName"
				}(),
			}

			userres, res, err := clnt.UserAPI.CreateUser(ctx).CreateUserRequest(req).Execute()
			if err != nil {
				return fmt.Errorf("error when calling `UserAPI.CreateUser`: %v", err)
			}

			if res.StatusCode != 200 {
				return fmt.Errorf("expected status code 200, got %d", res.StatusCode)
			}

			b, err := json.MarshalIndent(userres, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}

			fmt.Println(string(b))
			return nil
		},
	}
}
