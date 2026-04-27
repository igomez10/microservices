package scopefromrole

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "scope-from-role",
		Usage:     "Remove a scope from a role",
		ArgsUsage: "<role-id> <scope-id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     cliflags.EnvFlag,
				Usage:    "Environment to use (live, local)",
				Value:    auth.DefaultEnv,
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.UsernameFlag,
				Usage:    "Username for authentication (or set SOCIALAPP_CLI_USERNAME)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.PasswordFlag,
				Usage:    "Password for authentication (or set SOCIALAPP_CLI_PASSWORD)",
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

			envName := auth.ResolveEnvironment(cmd.String(cliflags.EnvFlag))

			username, password, err := auth.ResolveCredentials(
				cmd.String(cliflags.UsernameFlag),
				cmd.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{scopes.SocialappRolesUpdate.String()})
			if err != nil {
				return fmt.Errorf("failed to get authenticated client: %w", err)
			}

			host, scheme, err := auth.GetAPIClientConfig(envName)
			if err != nil {
				return err
			}

			cfg := client.NewConfiguration()
			cfg.Host = host
			cfg.Scheme = scheme
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			httpResponse, err := apiClient.RoleAPI.RemoveScopeFromRole(ctx, roleIDStr, scopeIDStr).Execute()
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
