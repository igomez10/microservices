package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "scope",
		Usage: "Create a new scope",
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
			&cli.StringFlag{
				Name:     cliflags.NameFlag,
				Usage:    "Scope name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     cliflags.DescriptionFlag,
				Usage:    "Scope description",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)
			name := cmd.String("name")
			description := cmd.String(cliflags.DescriptionFlag)

			oauthConfig := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenEndpoint,
				Scopes:       []string{scopes.SocialappScopesCreate.String()},
			}

			httpClient := oauthConfig.Client(ctx)

			cfg := client.NewConfiguration()
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("error parsing host URL: %v", err)
			}
			cfg.Host = parsedURL.Host
			cfg.Scheme = parsedURL.Scheme
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			scope := client.NewScope(name, description)

			createdScope, httpResponse, err := apiClient.ScopeAPI.CreateScope(ctx).Scope(*scope).Execute()
			if err != nil {
				if httpResponse != nil && httpResponse.Body != nil {
					bodyBytes, _ := io.ReadAll(httpResponse.Body)
					fmt.Printf("Error: %v\nResponse: %s\n", err, string(bodyBytes))
				} else {
					fmt.Printf("Error: %v\n", err)
				}
				return err
			}

			if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				return fmt.Errorf("unexpected status code: %d, body: %s", httpResponse.StatusCode, string(bodyBytes))
			}

			scopeJSON, err := json.MarshalIndent(createdScope, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(scopeJSON))
			return nil
		},
	}
}
