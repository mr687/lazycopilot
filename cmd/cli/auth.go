package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cli/oauth/device"
	"github.com/mr687/lazycopilot/pkg/config"
	"github.com/mr687/lazycopilot/pkg/utils"
	"github.com/spf13/cobra"
)

func newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Github",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("a valid subcommand is required. Use 'auth login' or 'auth logout'")
		},
	}

	cmd.AddCommand(loginCommand())
	cmd.AddCommand(logoutCommand())
	return cmd
}

func loginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Github",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := utils.GetConfigPath()
			filePath := utils.GetFilePath(configDir, "github-copilot", "apps.json")

			if utils.IsFileExists(filePath) {
				if user, err := utils.CheckCachedToken(filePath, config.ClientID); err == nil {
					fmt.Printf("Already authenticated as: %s\n", user.Login)
					return nil
				}
			}

			token, err := connectGithubDeviceFlow()
			if err != nil {
				return fmt.Errorf("failed to authenticate with Github: %v. Please try again.", err)
			}

			user, err := utils.FetchGithubUserInfo(token)
			if err != nil {
				return fmt.Errorf("failed to fetch user info: %v. Ensure your token is valid and try again.", err)
			}

			err = utils.SaveAuthToken(filePath, token, user.Login, config.ClientID)
			if err != nil {
				return fmt.Errorf("failed to save token to file: %v. Check your file permissions and try again.", err)
			}

			fmt.Println("Authentication successful.")
			return nil
		},
	}
	return cmd
}

func logoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Github",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := utils.GetConfigPath()
			filePath := utils.GetFilePath(configDir, "github-copilot", "apps.json")

			if !utils.IsFileExists(filePath) {
				return fmt.Errorf("no authentication found. You are not logged in.")
			}

			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove authentication file: %v. Ensure you have the necessary permissions and try again.", err)
			}

			fmt.Println("Logged out successfully. Please revoke the token manually from your GitHub account settings.")
			return nil
		},
	}
	return cmd
}

func connectGithubDeviceFlow() (string, error) {
	clientId := config.ClientID
	scopes := []string{"read:user"}
	httpClient := utils.GetHttpClient()
	code, err := device.RequestCode(httpClient, "https://github.com/login/device/code", clientId, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to request device code: %v. Check your internet connection and try again.", err)
	}

	fmt.Println("Code: ", code.UserCode)
	fmt.Println("Navigate to the URL and paste the code: ", code.VerificationURI)

	accessToken, err := device.Wait(
		context.TODO(),
		httpClient,
		"https://github.com/login/oauth/access_token",
		device.WaitOptions{
			ClientID:   clientId,
			DeviceCode: code,
			GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to wait for access token: %v. Ensure you completed the device flow and try again.", err)
	}
	return accessToken.Token, nil
}
