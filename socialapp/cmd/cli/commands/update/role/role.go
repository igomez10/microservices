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
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "role",
		Usage:     "Update a role",
		ArgsUsage: "<role-id>",
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
			roleID, err := strconv.ParseInt(roleIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid role ID: %v", err)
			}

			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)
			name := cmd.String(cliflags.NameFlag)
			description := cmd.String(cliflags.DescriptionFlag)

			oauthConfig := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenEndpoint,
				Scopes:       []string{scopes.SocialappRolesUpdate.String(), scopes.SocialappRolesRead.String()},
			}

			httpClient := oauthConfig.Client(ctx)

			cfg := client.NewConfiguration()
			cfg.Host = host
			cfg.Scheme = "http"
			cfg.HTTPClient = httpClient

			apiClient := client.NewAPIClient(cfg)

			// Get existing role to preserve fields not being updated
			existingRole, _, err := apiClient.RoleAPI.GetRole(ctx, int32(roleID)).Execute()
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

			role, httpResponse, err := apiClient.RoleAPI.UpdateRole(ctx, int32(roleID)).Role(*updatedRole).Execute()
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
