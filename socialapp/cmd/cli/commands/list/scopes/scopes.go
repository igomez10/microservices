package scopes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "scopes",
		Usage: "List all scopes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     cliflags.HostFlag,
				Usage:    "Host of the socialapp API",
				Value:    "http://localhost:8086",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.TokenEndpointFlag,
				Usage:    "Token endpoint of the socialapp API",
				Value:    "http://localhost:8086/v1/oauth/token",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.UsernameFlag,
				Aliases:  []string{"u"},
				Usage:    "Username of the socialapp API",
				Value:    "admin",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.PasswordFlag,
				Aliases:  []string{"p"},
				Usage:    "Password of the socialapp API",
				Value:    "admin",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "limit",
				Usage:    "Maximum number of scopes to return",
				Value:    100,
				Required: false,
			},
			&cli.IntFlag{
				Name:     "offset",
				Usage:    "Number of scopes to skip",
				Value:    0,
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)
			limit := int32(cmd.Int("limit"))
			offset := int32(cmd.Int("offset"))

			// Parse the host URL to extract scheme and host separately
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("error parsing host URL: %v", err)
			}

			oauthConfig := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenEndpoint,
				Scopes:       []string{"socialapp.scopes.list"},
			}
			source := oauthConfig.TokenSource(ctx)
			ctx = context.WithValue(ctx, client.ContextOAuth2, source)
			httpClient := oauthConfig.Client(ctx)

			cfg := client.NewConfiguration()
			cfg.Host = parsedURL.Host
			cfg.Scheme = parsedURL.Scheme
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			scopes, httpResponse, err := apiClient.ScopeAPI.ListScopes(ctx).Limit(limit).Offset(offset).Execute()
			if err != nil {
				if httpResponse != nil && httpResponse.Body != nil {
					bodyBytes, _ := io.ReadAll(httpResponse.Body)
					fmt.Printf("Error: %v\nResponse: %s\n", err, string(bodyBytes))
				} else {
					fmt.Printf("Error: %v\n", err)
				}
				return err
			}

			if httpResponse.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				return fmt.Errorf("unexpected status code: %d, body: %s", httpResponse.StatusCode, string(bodyBytes))
			}

			scopesJSON, err := json.MarshalIndent(scopes, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(scopesJSON))
			return nil
		},
	}
}
