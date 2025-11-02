package scopetorole

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "scope-to-role",
		Usage:     "Add scope(s) to a role",
		ArgsUsage: "<role-id> <scope-name> [scope-name...]",
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 2 {
				return fmt.Errorf("role ID and at least one scope name are required")
			}

			roleIDStr := cmd.Args().Get(0)
			roleID, err := strconv.ParseInt(roleIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid role ID: %v", err)
			}

			// Collect all scope names
			scopeNames := make([]string, 0, cmd.Args().Len()-1)
			for i := 1; i < cmd.Args().Len(); i++ {
				scopeNames = append(scopeNames, cmd.Args().Get(i))
			}

			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)

		oauthConfig := clientcredentials.Config{
			ClientID:     username,
			ClientSecret: password,
			TokenURL:     tokenEndpoint,
			Scopes:       []string{scopes.SocialappRolesUpdate.String()},
		}

			httpClient := oauthConfig.Client(ctx)

			cfg := client.NewConfiguration()
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("invalid host URL: %v", err)
			}
			cfg.Host = parsedURL.Host
			cfg.Scheme = parsedURL.Scheme
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			httpResponse, err := apiClient.RoleAPI.AddScopeToRole(ctx, int32(roleID)).RequestBody(scopeNames).Execute()
			if err != nil {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				fmt.Printf("Error: %v\nResponse: %s\n", err, string(bodyBytes))
				return err
			}

			if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				return fmt.Errorf("unexpected status code: %d, body: %s", httpResponse.StatusCode, string(bodyBytes))
			}

			fmt.Printf("Added scopes %v to role %d successfully\n", scopeNames, roleID)
			return nil
		},
	}
}
