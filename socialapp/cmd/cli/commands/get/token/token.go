package gettoken

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2"
)

type tokenOutput struct {
	AccessToken     string   `json:"access_token"`
	TokenType       string   `json:"token_type"`
	ExpiresIn       int64    `json:"expires_in,omitempty"`
	RequestedScopes []string `json:"requested_scopes,omitempty"`
}

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "token",
		Usage: "get a new OAuth token",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.EnvFlag,
				Usage: "Environment to use (live, local)",
				Value: auth.DefaultEnv,
			},
			&cli.StringFlag{
				Name:  cliflags.HostFlag,
				Usage: "Base host URL to use for token endpoint (e.g., http://localhost:8086). Overrides --env",
			},
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Usage:   "Username for authentication (or set SOCIALAPP_CLI_USERNAME)",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
				Usage:   "Password for authentication (or set SOCIALAPP_CLI_PASSWORD)",
				Aliases: []string{"p"},
			},
			&cli.StringSliceFlag{
				Name:    cliflags.ScopesFlag,
				Usage:   "OAuth scopes to request (repeat flag or provide comma-separated values)",
				Aliases: []string{"s"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			envName := auth.ResolveEnvironment(c.String(cliflags.EnvFlag))
			hostURL := c.String(cliflags.HostFlag)

			username, password, err := auth.ResolveCredentials(
				c.String(cliflags.UsernameFlag),
				c.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			requestedScopes := c.StringSlice(cliflags.ScopesFlag)

			var token *oauth2.Token
			if hostURL != "" {
				host, scheme, err := auth.ParseHostURL(hostURL)
				if err != nil {
					return err
				}
				if host == "" || scheme == "" {
					return fmt.Errorf("invalid host URL %q (expected format: http(s)://host[:port])", hostURL)
				}

				normalizedHostURL := fmt.Sprintf("%s://%s", scheme, host)
				tokenEndpoint := normalizedHostURL + "/v1/oauth/token"
				cacheEnvironment := "host:" + normalizedHostURL

				token, err = auth.GetFreshTokenWithTokenEndpoint(
					ctx,
					cacheEnvironment,
					tokenEndpoint,
					username,
					password,
					requestedScopes,
				)
				if err != nil {
					return fmt.Errorf("failed to get token: %w", err)
				}
			} else {
				token, err = auth.GetFreshToken(ctx, envName, username, password, requestedScopes)
				if err != nil {
					return fmt.Errorf("failed to get token: %w", err)
				}
			}

			output := tokenOutput{
				AccessToken:     token.AccessToken,
				TokenType:       token.TokenType,
				RequestedScopes: requestedScopes,
			}
			if !token.Expiry.IsZero() {
				output.ExpiresIn = int64(time.Until(token.Expiry).Seconds())
			}

			tokenJSON, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal token response: %w", err)
			}

			fmt.Println(string(tokenJSON))
			return nil
		},
	}
}
