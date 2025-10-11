package list

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/users"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list users",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
			},

			&cli.StringFlag{
				Name:  "host",
				Value: "http://localhost:8086",
			},
			&cli.StringFlag{
				Name:  "tokenendpoint",
				Value: "http://localhost:8086/v1/oauth/token",
			},
			&cli.IntFlag{
				Name:    "pagesize",
				Aliases: []string{"ps"},
				Value:   10,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Value:   0,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			username1 := c.String("username")
			password := c.String("password")
			tokenendpoint := c.String("tokenendpoint")
			pagesize := c.Int("pagesize")
			offset := c.Int("offset")
			host := c.String("host")

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
