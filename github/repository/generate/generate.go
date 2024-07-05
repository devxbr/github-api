package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RepoTemplateRequest struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

func (r *RepoTemplateRequest) Generate() {

	templateOwner := "devxbr"
	templateRepo := "template"

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/generate", templateOwner, templateRepo)

	jsonData, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Erro ao codificar os dados do JSON:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.baptiste-preview+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("Repositório criado com sucesso!")
	} else {
		fmt.Printf("Falha ao criar repositório: %s\n", resp.Status)
	}
}
