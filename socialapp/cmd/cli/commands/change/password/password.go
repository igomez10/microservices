package password

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "password",
		Usage: "Change the password of the current user",
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
			host := cmd.String(cliflags.HostFlag)
			tokenEndpoint := cmd.String(cliflags.TokenEndpointFlag)
			username := cmd.String(cliflags.UsernameFlag)
			password := cmd.String(cliflags.PasswordFlag)
			oldPassword := cmd.String("oldpassword")
			newPassword := cmd.String("newpassword")

			oauthConfig := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenEndpoint,
				Scopes:       []string{"socialapp.users.update"},
			}

			httpClient := oauthConfig.Client(ctx)

			cfg := client.NewConfiguration()
			cfg.Host = host
			cfg.Scheme = "http"
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
