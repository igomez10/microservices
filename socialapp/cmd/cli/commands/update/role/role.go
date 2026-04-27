package role

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
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "role",
		Usage:     "Update a role",
		ArgsUsage: "<role-id>",
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
				Usage:    "Role name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.DescriptionFlag,
				Usage:    "Role description",
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("role ID is required")
			}

			roleIDStr := cmd.Args().Get(0)
			if _, err := strconv.ParseInt(roleIDStr, 10, 32); err != nil {
				return fmt.Errorf("invalid role ID: %v", err)
			}

			envName := auth.ResolveEnvironment(cmd.String(cliflags.EnvFlag))

			username, password, err := auth.ResolveCredentials(
				cmd.String(cliflags.UsernameFlag),
				cmd.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			name := cmd.String(cliflags.NameFlag)
			description := cmd.String(cliflags.DescriptionFlag)

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{scopes.SocialappRolesUpdate.String(), scopes.SocialappRolesRead.String()})
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

			// Get existing role to preserve fields not being updated
			existingRole, _, err := apiClient.RoleAPI.GetRole(ctx, roleIDStr).Execute()
			if err != nil {
				return fmt.Errorf("failed to get existing role: %v", err)
			}

			// Use existing values if not provided
			if name == "" {
				name = existingRole.Name
			}
			if description == "" && existingRole.Description != nil {
				description = *existingRole.Description
			}

			updatedRole := client.NewRole(name)
			if description != "" {
				updatedRole.Description = &description
			}

			role, httpResponse, err := apiClient.RoleAPI.UpdateRole(ctx, roleIDStr).Role(*updatedRole).Execute()
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

			roleJSON, err := json.MarshalIndent(role, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(roleJSON))
			return nil
		},
	}
}
