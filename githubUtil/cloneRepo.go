package githubUtil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	mod "task/models"
)

const gitHubAPIURL = "https://api.github.com/user/repos"

func CreateRepo(token string, repo mod.Repo) error {
	jsonData, err := json.Marshal(repo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", gitHubAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error: %s, response: %s", resp.Status, body)
	}

	fmt.Println("Successfully created repository:", repo.Name)
	return nil
}

func GitClone(username string, repo mod.Repo, dir string) error {
	url := fmt.Sprintf("https://github.com/%s/%s.git", username, repo.Name)
	cmd := exec.Command("git", "clone", url, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
