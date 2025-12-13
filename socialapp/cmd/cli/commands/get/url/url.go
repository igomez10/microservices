package geturl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/igomez10/microservices/socialapp/cmd/cli/pkg/auth"
	"github.com/igomez10/microservices/socialapp/pkg/scopes"
	"github.com/urfave/cli/v3"
)

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:  "url",
		Usage: "get url",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.EnvFlag,
				Usage: "Environment to use (live, local)",
				Value: auth.DefaultEnv,
			},
			&cli.StringFlag{
				Name:    cliflags.AliasFlag,
				Aliases: []string{"a"},
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
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			envName := auth.ResolveEnvironment(c.String(cliflags.EnvFlag))

			username, password, err := auth.ResolveCredentials(
				c.String(cliflags.UsernameFlag),
				c.String(cliflags.PasswordFlag),
			)
			if err != nil {
				return err
			}

			alias := c.String(cliflags.AliasFlag)
			if alias == "" {
				return fmt.Errorf("alias is required")
			}

			httpClient, err := auth.GetHTTPClient(ctx, envName, username, password, []string{scopes.ShortlyUrlCreate.String(), scopes.ShortlyUrlDelete.String()})
			if err != nil {
				return fmt.Errorf("failed to get authenticated client: %w", err)
			}

			host, scheme, err := auth.GetAPIClientConfig(envName)
			if err != nil {
				return err
			}

			configuration := client.NewConfiguration()
			configuration.Host = host
			configuration.Scheme = scheme
			configuration.HTTPClient = httpClient

			clnt := client.NewAPIClient(configuration)

			res, err := clnt.URLAPI.GetUrl(ctx, alias).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v \n", res)
				return fmt.Errorf("error when calling `URLAPI.GetUrl`: %v", err)
			}

			b, err := json.MarshalIndent(res.Body, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}
			fmt.Println(string(b))
			return nil
		},
	}
}
