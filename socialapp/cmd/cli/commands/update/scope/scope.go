package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "scope",
		Usage:     "Update a scope",
		ArgsUsage: "<scope-id>",
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
			if _, err := strconv.ParseInt(scopeIDStr, 10, 32); err != nil {
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

			name := cmd.String("name")
			description := cmd.String("description")

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.scopes.update", "socialapp.scopes.read"})
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

			// Get existing scope to preserve fields not being updated
			existingScope, _, err := apiClient.ScopeAPI.GetScope(ctx, scopeIDStr).Execute()
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

			scope, httpResponse, err := apiClient.ScopeAPI.UpdateScope(ctx, scopeIDStr).Scope(*updatedScope).Execute()
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
