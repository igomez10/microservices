package unfollowuser

import (
	"context"
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
		Usage:     "unfollow a user",
		ArgsUsage: "<followed_username> <follower_username>",
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
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			tokenendpoint := c.String(cliflags.TokenEndpointFlag)
			host := c.String(cliflags.HostFlag)
			if host != defaultHost && c.String(cliflags.TokenEndpointFlag) == fmt.Sprintf("%s/v1/oauth/token", defaultHost) {
				tokenendpoint = fmt.Sprintf("%s/v1/oauth/token", host)
			}
			username := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)

			if c.Args().Len() < 2 {
				return fmt.Errorf("both followed_username and follower_username arguments are required")
			}
			followedUsername := c.Args().Get(0)
			followerUsername := c.Args().Get(1)

			oauth2Config := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenendpoint,
				Scopes:       []string{"socialapp.follower.delete"},
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

			res, err := clnt.UserAPI.UnfollowUser(ctx, followedUsername, followerUsername).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v\n", res)
				return fmt.Errorf("error when calling `UserAPI.UnfollowUser`: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("expected status code 200, got %d", res.StatusCode)
			}

			fmt.Printf("%s has unfollowed %s\n", followerUsername, followedUsername)
			return nil
		},
	}
}
