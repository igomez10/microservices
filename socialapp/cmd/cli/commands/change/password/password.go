package password

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
		Name:  "password",
		Usage: "Change the password of the current user",
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
				Name:     "oldpassword",
				Usage:    "Old password",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "newpassword",
				Usage:    "New password",
				Required: true,
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

			oldPassword := cmd.String("oldpassword")
			newPassword := cmd.String("newpassword")

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{"socialapp.users.update"})
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

			changePasswordReq := *client.NewChangePasswordRequest(oldPassword, newPassword)

			user, httpResponse, err := apiClient.UserAPI.ChangePassword(ctx).ChangePasswordRequest(changePasswordReq).Execute()
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

			userJSON, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(userJSON))
			return nil
		},
	}
}
