package geturl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

const defaultHost = "http://localhost:8086"

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "url",
		Usage: "get url",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    cliflags.AliasFlag,
				Aliases: []string{"a"},
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
			username1 := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)
			alias := c.String(cliflags.AliasFlag)

			if alias == "" {
				return fmt.Errorf("alias is required")
			}

		oauth2Config := clientcredentials.Config{
			ClientID:     username1,
			ClientSecret: password,
			TokenURL:     tokenendpoint,
			Scopes:       []string{scopes.ShortlyUrlCreate.String(), scopes.ShortlyUrlDelete.String()},
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

			res, err := clnt.URLAPI.GetUrl(ctx, alias).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v \n", res)
				return fmt.Errorf("error when calling `URLAPI.GetUrl`: %v", err)
			}

			b, err := json.MarshalIndent(res.Body, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}
			fmt.Println(string(b))
			return nil
		},
	}
}
