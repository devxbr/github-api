package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.github.com"
)

type RepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
	HasIssues   bool   `json:"has_issues"`
	HasWiki     bool   `json:"has_wiki"`
	AutoInit    bool   `json:"auto_init"`
}

func (r *RepoRequest) Create() {
	jsonData, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Erro ao converter os dados para JSON:", err)
		return
	}

	var url string
	if os.Getenv("ORGANIZATION_NAME") == "" {
		url = fmt.Sprintf("%s/user/repos", apiURL)
	} else {
		url = fmt.Sprintf("%s/orgs/%s/repos", apiURL, os.Getenv("ORGANIZATION_NAME"))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao criar a solicitação HTTP:", err)
		return
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a solicitação:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Erro: %d\n", resp.StatusCode)
		return
	}
	fmt.Println("Repositório criado com sucesso.")
}
