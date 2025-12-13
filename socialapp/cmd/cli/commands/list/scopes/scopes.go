package scopes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "scopes",
		Usage: "List all scopes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     cliflags.EnvFlag,
				Usage:    "Environment to use (live, local)",
				Value:    auth.DefaultEnv,
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.UsernameFlag,
				Aliases:  []string{"u"},
				Usage:    "Username for authentication (or set SOCIALAPP_CLI_USERNAME)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     cliflags.PasswordFlag,
				Aliases:  []string{"p"},
				Usage:    "Password for authentication (or set SOCIALAPP_CLI_PASSWORD)",
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
			envName := auth.ResolveEnvironment(cmd.String(cliflags.EnvFlag))

			username, password, err := auth.ResolveCredentials(
				cmd.String(cliflags.UsernameFlag),
				cmd.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			limit := int32(cmd.Int("limit"))
			offset := int32(cmd.Int("offset"))

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.scopes.list"})
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
