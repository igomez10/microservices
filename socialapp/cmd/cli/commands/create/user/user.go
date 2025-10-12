package createuser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
)

const defaultHost = "http://localhost:8086"

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Usage: "create user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:  cliflags.HostFlag,
				Value: defaultHost,
			},
			&cli.StringFlag{
				Name:  cliflags.TokenEndpointFlag,
				Value: fmt.Sprintf("%s/v1/oauth/token", defaultHost),
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
			// tokenendpoint := c.String(cliflags.TokenEndpointFlag)
			host := c.String(cliflags.HostFlag)
			if host != defaultHost && c.String(cliflags.TokenEndpointFlag) == fmt.Sprintf("%s/v1/oauth/token", defaultHost) {
				// tokenendpoint = fmt.Sprintf("%s/v1/oauth/token", host)
			}
			username1 := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)
			email := c.String(cliflags.EmailFlag)
			firstname := c.String(cliflags.FirstNameFlag)
			lastname := c.String(cliflags.LastNameFlag)

			if username1 == "" {
				return fmt.Errorf("username is required")
			}
			if password == "" {
				return fmt.Errorf("password is required")
			}
			if email == "" {
				return fmt.Errorf("email is required")
			}

			// Parse the host URL to extract scheme and host separately
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("error parsing host URL: %v", err)
			}

			configuration := client.NewConfiguration()
			configuration.Host = parsedURL.Host
			configuration.Scheme = parsedURL.Scheme
			configuration.HTTPClient = http.DefaultClient

			clnt := client.NewAPIClient(configuration)

			req := client.CreateUserRequest{
				Username: username1,
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

			if res.StatusCode != 201 {
				return fmt.Errorf("expected status code 201, got %d", res.StatusCode)
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
