package githubUtil

import (
	"fmt"
	"net/http"
)

func RepoExists(token, username, repoName string) bool {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
