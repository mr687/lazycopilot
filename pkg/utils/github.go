package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type CachedToken struct {
	User        string `json:"user"`
	OauthToken  string `json:"oauth_token"`
	GithubAppID string `json:"githubAppId"`
}

type GithubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	// Add other fields as needed
}

func GetFilePath(configDir, subDir, fileName string) string {
	return filepath.Join(configDir, subDir, fileName)
}

func CheckCachedToken(filePath, clientId string) (*GithubUser, error) {
	var apps map[string]CachedToken
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open apps.json: %v", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&apps)
	if err != nil {
		return nil, fmt.Errorf("failed to decode apps.json: %v", err)
	}

	cachedToken, exists := apps["github.com:"+clientId]
	if !exists {
		return nil, errors.New("cached token not found")
	}

	user, err := FetchGithubUserInfo(cachedToken.OauthToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FetchGithubUserInfo(token string) (*GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %s", resp.Status)
	}

	var user GithubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func SaveAuthToken(filePath, token, username, clientId string) error {
	var apps map[string]CachedToken
	if IsFileExists(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open apps.json: %v", err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(&apps)
		if err != nil {
			return fmt.Errorf("failed to decode apps.json: %v", err)
		}
	} else {
		apps = make(map[string]CachedToken)
	}

	apps["github.com:"+clientId] = CachedToken{
		User:        username,
		OauthToken:  token,
		GithubAppID: clientId,
	}

	data, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal apps.json: %v", err)
	}

	err = SaveFile(filePath, data)
	if err != nil {
		return fmt.Errorf("failed to save apps.json: %v", err)
	}

	return nil
}

func GetHttpClient() *http.Client {
	return &http.Client{}
}
