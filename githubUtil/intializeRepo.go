package githubUtil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	mod "task/models"
	util "task/utility"
	"time"
)

const (
	goTemplate = `# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /go-app

EXPOSE 8080

CMD [ "/go-app" ]
`
	nextjsTemplate = `# syntax=docker/dockerfile:1
FROM node:18-alpine

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install

COPY . .

RUN yarn build

EXPOSE 3000

CMD [ "yarn", "start" ]
`
	wordpressTemplate = `# syntax=docker/dockerfile:1
FROM wordpress:latest

COPY ./wp-content /var/www/html/wp-content
`
)

func createFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}

func InitializeRepo(token, username string, repo mod.Repo) error {
	dir, err := ioutil.TempDir("", repo.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if !RepoExists(token, username, repo.Name) {
		if err := CreateRepo(token, repo); err != nil {
			fmt.Printf("Failed to create repository: %s\n", repo.Name)
			fmt.Println("Error:", err)
			return err
		}
		// Sleep for a brief period to allow GitHub to register the repository
		time.Sleep(2 * time.Second)
	}

	if err := GitClone(username, repo, dir); err != nil {
		return err
	}

	// Add Dockerfile
	switch repo.Name {
	case "go-app":
		if err := createFile(filepath.Join(dir, "Dockerfile"), goTemplate); err != nil {
			return err
		}
		// Additional Go-specific setup (e.g., go.mod, go.sum, main.go)
		if err := createFile(filepath.Join(dir, "main.go"), "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"); err != nil {
			return err
		}
		if err := createFile(filepath.Join(dir, "go.mod"), "module github.com/"+username+"/"+repo.Name+"\ngo 1.16"); err != nil {
			return err
		}
		if err := createFile(filepath.Join(dir, "go.sum"), ""); err != nil {
			return err
		}
		if err := createFile(filepath.Join(dir, "main_test.go"), "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {\n\tt.Run(\"Test Hello World\", func(t *testing.T) {\n\t\twant := \"Hello, World!\"\n\t\tif got := \"Hello, World!\"; got != want {\n\t\t\tt.Errorf(\"got %s, want %s\", got, want)\n\t\t}\n\t})\n}"); err != nil {
			return err
		}
	case "nextjs-app":
		if err := createFile(filepath.Join(dir, "Dockerfile"), nextjsTemplate); err != nil {
			return err
		}
		// Additional Next.js-specific setup (e.g., package.json, tsconfig.json)
	case "wordpress-site":
		if err := createFile(filepath.Join(dir, "Dockerfile"), wordpressTemplate); err != nil {
			return err
		}
		// Additional WordPress-specific setup
	}

	// Add CI configuration
	if err := util.SetupCI(dir, repo.Name); err != nil {
		return err
	}

	// Commit and push changes
	if err := GitCommitAndPush(dir, username, repo.Name); err != nil {
		return err
	}

	fmt.Println("Successfully initialized repository:", repo.Name)
	return nil
}
