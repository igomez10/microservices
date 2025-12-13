package scopefromrole

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "scope-from-role",
		Usage:     "Remove a scope from a role",
		ArgsUsage: "<role-id> <scope-id>",
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
			if cmd.Args().Len() != 2 {
				return fmt.Errorf("role ID and scope ID are required")
			}

			roleIDStr := cmd.Args().Get(0)
			roleID, err := strconv.ParseInt(roleIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid role ID: %v", err)
			}

			scopeIDStr := cmd.Args().Get(1)
			scopeID, err := strconv.ParseInt(scopeIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid scope ID: %v", err)
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
			cfg.Host = host
			cfg.Scheme = "http"
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			httpResponse, err := apiClient.RoleAPI.RemoveScopeFromRole(ctx, roleID, scopeID).Execute()
			if err != nil {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				fmt.Printf("Error: %v\nResponse: %s\n", err, string(bodyBytes))
				return err
			}

			if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
				bodyBytes, _ := io.ReadAll(httpResponse.Body)
				return fmt.Errorf("unexpected status code: %d, body: %s", httpResponse.StatusCode, string(bodyBytes))
			}

			fmt.Printf("Removed scope %d from role %d successfully\n", scopeID, roleID)
			return nil
		},
	}
}
