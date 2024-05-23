package main

import (
	"fmt"
	"os"
	gitutil "task/githubUtil"
	mod "task/models"
	"time"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable is required")
		return
	}

	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		fmt.Println("GITHUB_USERNAME environment variable is required")
		return
	}

	repos := []mod.Repo{
		{Name: "go-app", Private: false},
		{Name: "nextjs-app", Private: false},
		{Name: "wordpress-site", Private: false},
	}

	for _, repo := range repos {
		if !gitutil.RepoExists(token, username, repo.Name) {
			if err := gitutil.CreateRepo(token, repo); err != nil {
				fmt.Printf("Failed to create repository: %s\n", repo.Name)
				fmt.Println("Error:", err)
				continue
			}
		}
		// Sleep for a brief period to allow GitHub to register the repository
		time.Sleep(2 * time.Second)
		if err := gitutil.InitializeRepo(token, username, repo); err != nil {
			fmt.Printf("Failed to initialize repository: %s\n", repo.Name)
			fmt.Println("Error:", err)
		}
	}
}
