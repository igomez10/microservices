package listusers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/users"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

const defaultHost = "http://localhost:8086"

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "list users",
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
			tokenendpoint := c.String(cliflags.TokenEndpointFlag)
			host := c.String(cliflags.HostFlag)
			if host != defaultHost && c.String(cliflags.TokenEndpointFlag) == fmt.Sprintf("%s/v1/oauth/token", defaultHost) {
				tokenendpoint = fmt.Sprintf("%s/v1/oauth/token", host)
			}
			username1 := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)
			pagesize := c.Int(cliflags.PageSizeFlag)
			offset := c.Int(cliflags.OffsetFlag)

			u, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("Error when parsing url: %v", err)
			}

			ctx, apiClient := users.GetApiClient(ctx, u)

			conf := clientcredentials.Config{
				ClientID:     username1,
				ClientSecret: password,
				Scopes:       []string{"socialapp.users.list"},
				TokenURL:     tokenendpoint,
			}
			ctx = context.WithValue(ctx, client.ContextOAuth2, conf.TokenSource(ctx))

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
