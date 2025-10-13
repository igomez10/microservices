package updateuser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

const defaultHost = "http://localhost:8086"

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "user",
		Usage:     "update user by username",
		ArgsUsage: "<username>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.HostFlag,
				Value: defaultHost,
			},
			&cli.StringFlag{
				Name:  cliflags.TokenEndpointFlag,
				Value: fmt.Sprintf("%s/v1/oauth/token", defaultHost),
			},
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
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
			tokenendpoint := c.String(cliflags.TokenEndpointFlag)
			host := c.String(cliflags.HostFlag)
			if host != defaultHost && c.String(cliflags.TokenEndpointFlag) == fmt.Sprintf("%s/v1/oauth/token", defaultHost) {
				tokenendpoint = fmt.Sprintf("%s/v1/oauth/token", host)
			}
			username := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)
			email := c.String(cliflags.EmailFlag)
			firstname := c.String(cliflags.FirstNameFlag)
			lastname := c.String(cliflags.LastNameFlag)

			if c.Args().Len() < 1 {
				return fmt.Errorf("username argument is required")
			}
			targetUsername := c.Args().Get(0)

			oauth2Config := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenendpoint,
				Scopes:       []string{"socialapp.users.update", "socialapp.users.read"},
			}
			httpClient := oauth2Config.Client(ctx)

			// Parse the host URL to extract scheme and host separately
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("error parsing host URL: %v", err)
			}

			configuration := client.NewConfiguration()
			configuration.Host = parsedURL.Host
			configuration.Scheme = parsedURL.Scheme
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
