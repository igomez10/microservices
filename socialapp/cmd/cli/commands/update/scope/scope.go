package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "scope",
		Usage:     "Update a scope",
		ArgsUsage: "<scope-id>",
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
				Usage:    "Username of the socialapp API",
				Value:    "admin",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.PasswordFlag,
				Usage:    "Password of the socialapp API",
				Value:    "admin",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.NameFlag,
				Usage:    "Scope name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.DescriptionFlag,
				Usage:    "Scope description",
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("scope ID is required")
			}

			scopeIDStr := cmd.Args().Get(0)
			scopeID, err := strconv.ParseInt(scopeIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid scope ID: %v", err)
			}

			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)
			name := cmd.String("name")
			description := cmd.String("description")

			oauthConfig := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenEndpoint,
				Scopes:       []string{"socialapp.scopes.update", "socialapp.scopes.read"},
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

			// Get existing scope to preserve fields not being updated
			existingScope, _, err := apiClient.ScopeAPI.GetScope(ctx, int32(scopeID)).Execute()
			if err != nil {
				return fmt.Errorf("failed to get existing scope: %v", err)
			}

			// Use existing values if not provided
			if name == "" {
				name = existingScope.Name
			}
			if description == "" {
				description = existingScope.Description
			}

			updatedScope := client.NewScope(name, description)

			scope, httpResponse, err := apiClient.ScopeAPI.UpdateScope(ctx, int32(scopeID)).Scope(*updatedScope).Execute()
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

			scopeJSON, err := json.MarshalIndent(scope, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(scopeJSON))
			return nil
		},
	}
}
