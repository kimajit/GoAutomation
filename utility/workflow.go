package utility

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func SetupCI(dir, repoName string) error {
	ciContent := ""
	switch repoName {
	case "go-app":
		ciContent = `name: Go CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test ./...

    - name: Lint
      run: go run github.com/golangci/golangci-lint/cmd/golangci-lint run

    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    - name: Build Docker image
      run: docker build -t myusername/go-app:${{ github.sha }} .

    - name: Push Docker image
      run: docker push myusername/go-app:${{ github.sha }}
`
	case "nextjs-app":
		ciContent = `name: Next.js CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Node.js
      uses: actions/setup-node@v2
      with:
        node-version: 18

    - name: Install dependencies
      run: yarn install

    - name: Lint
      run: yarn lint

    - name: Run tests
      run: yarn test

    - name: Build
      run: yarn build

    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    - name: Build Docker image
      run: docker build -t myusername/nextjs-app:${{ github.sha }} .

    - name: Push Docker image
      run: docker push myusername/nextjs-app:${{ github.sha }}
`
	case "wordpress-site":
		ciContent = `name: WordPress CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up PHP
      uses: shivammathur/setup-php@v2
      with:
        php-version: 8.0

    - name: Install dependencies
      run: composer install

    - name: Lint
      run: vendor/bin/phpcs --standard=WordPress wp-content/

    - name: Run tests
      run: vendor/bin/phpunit

    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    - name: Build Docker image
      run: docker build -t myusername/wordpress-site:${{ github.sha }} .

    - name: Push Docker image
      run: docker push myusername/wordpress-site:${{ github.sha }}
`
	}

	ciDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(ciDir, 0755); err != nil {
		return err
	}

	ciFile := filepath.Join(ciDir, "ci.yml")
	return createFile(ciFile, ciContent)
}
func createFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}
