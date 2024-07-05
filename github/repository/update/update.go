package update

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UpdateRequest struct {
	Message  string `json:"message"`
	Content  string `json:"content"`
	Branch   string `json:"branch"`
	Path     string
	Owner    string
	RepoName string
}

func (t *UpdateRequest) Update() {

	// Codifica o conteúdo do arquivo em base64
	encodedContent := base64.StdEncoding.EncodeToString([]byte(t.Content))

	// Estrutura dos dados da solicitação
	data := map[string]string{
		"message": t.Message,
		"content": encodedContent,
		"branch":  "main", // ou outra branch, se necessário
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao converter os dados para JSON:", err)
		return
	}

	// Build URL do endpoint da API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", t.Owner, t.RepoName, t.Path)

	// Cria uma solicitação HTTP PUT
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao criar a solicitação HTTP:", err)
		return
	}

	// Define os cabeçalhos da solicitação
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Arquivo criado com sucesso.")
	} else if resp.StatusCode == http.StatusOK {
		fmt.Println("Arquivo atualizado com sucesso.")
	} else {
		fmt.Printf("Erro: %d\n", resp.StatusCode)
		fmt.Println(string(body))
	}
}
