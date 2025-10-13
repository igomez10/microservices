package getcomment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/igomez10/microservices/socialapp/client"
	"github.com/igomez10/microservices/socialapp/cmd/cli/cliflags"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2/clientcredentials"
)

const defaultHost = "http://localhost:8086"

func GetCmd() *cli.Command {
	return &cli.Command{
		Name:      "comment",
		Usage:     "get comment by id",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cliflags.HostFlag,
				Value: defaultHost,
			},
			&cli.StringFlag{
				Name:  cliflags.TokenEndpointFlag,
				Value: fmt.Sprintf("%s/v1/oauth/token", defaultHost),
			},
			&cli.StringFlag{
				Name:    cliflags.UsernameFlag,
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    cliflags.PasswordFlag,
				Aliases: []string{"p"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			tokenendpoint := c.String(cliflags.TokenEndpointFlag)
			host := c.String(cliflags.HostFlag)
			if host != defaultHost && c.String(cliflags.TokenEndpointFlag) == fmt.Sprintf("%s/v1/oauth/token", defaultHost) {
				tokenendpoint = fmt.Sprintf("%s/v1/oauth/token", host)
			}
			username := c.String(cliflags.UsernameFlag)
			password := c.String(cliflags.PasswordFlag)

			if c.Args().Len() < 1 {
				return fmt.Errorf("comment id argument is required")
			}
			commentIDStr := c.Args().Get(0)
			commentID, err := strconv.ParseInt(commentIDStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid comment id: %v", err)
			}

			oauth2Config := clientcredentials.Config{
				ClientID:     username,
				ClientSecret: password,
				TokenURL:     tokenendpoint,
				Scopes:       []string{"socialapp.comments.read"},
			}
			httpClient := oauth2Config.Client(ctx)

			// Parse the host URL to extract scheme and host separately
			parsedURL, err := url.Parse(host)
			if err != nil {
				return fmt.Errorf("error parsing host URL: %v", err)
			}

			configuration := client.NewConfiguration()
			configuration.Host = parsedURL.Host
			configuration.Scheme = parsedURL.Scheme
			configuration.HTTPClient = httpClient

			clnt := client.NewAPIClient(configuration)

			comment, res, err := clnt.CommentAPI.GetComment(ctx, int32(commentID)).Execute()
			if err != nil {
				fmt.Printf("Full HTTP response: %v\n", res)
				return fmt.Errorf("error when calling `CommentAPI.GetComment`: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("expected status code 200, got %d", res.StatusCode)
			}

			b, err := json.MarshalIndent(comment, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling response: %v", err)
			}

			fmt.Println(string(b))
			return nil
		},
	}
}
