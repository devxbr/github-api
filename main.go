package main

import (
	"log"

	"github.com/devxbr/github/repository/create"
	"github.com/devxbr/github/repository/scaffold"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error on load .env: %v", err)
	}

	newRepoName := "repositorio-test"
	newRepoDescription := "repository test for github api"

	repo := create.RepoRequest{
		Name:        newRepoName,
		Description: newRepoDescription,
		Private:     false,
		HasIssues:   false,
		HasWiki:     false,
		AutoInit:    true,
	}
	repo.Create()
	print("Repository created successfully!")

	g := scaffold.TemplateRequest{
		Files:    []string{"chart/values.qa.yaml", "chart/values.yaml"},
		RepoName: newRepoName,
	}

	g.Scaffold()
	print("Scaffolded done!")

}
