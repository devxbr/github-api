package scaffold

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	branch    = "main" // ou outra branch, se necessário
	commitMsg = "Initial commit"
)

type Blob struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Tree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
}

type TemplateRequest struct {
	Files    []string
	RepoName string
}

func (t *TemplateRequest) Scaffold() {
	blobs := make(map[string]Blob)

	for _, file := range t.Files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}

		// Cria um blob para cada arquivo
		blobSHA, err := t.createBlob(content)
		if err != nil {
			fmt.Println("Erro ao criar o blob:", err)
			return
		}

		blobs[file] = blobSHA
	}

	// Cria uma árvore que referencia os blobs
	treeSHA, err := t.createTree(blobs)
	if err != nil {
		fmt.Println("Erro ao criar a árvore:", err)
		return
	}

	// Cria um commit que referencia a árvore
	err = t.createCommit(treeSHA)
	if err != nil {
		fmt.Println("Erro ao criar o commit:", err)
	} else {
		fmt.Println("Commit criado com sucesso.")
	}
}

func (t *TemplateRequest) createBlob(content []byte) (Blob, error) {
	encodedContent := base64.StdEncoding.EncodeToString(content)

	data := map[string]string{
		"content":  encodedContent,
		"encoding": "base64",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return Blob{}, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/blobs", os.Getenv("ORGANIZATION_NAME"), t.RepoName)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Blob{}, err
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Blob{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Blob{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return Blob{}, fmt.Errorf("erro ao criar o blob: %s", string(body))
	}

	var blob Blob
	err = json.Unmarshal(body, &blob)
	if err != nil {
		return Blob{}, err
	}

	return blob, nil
}

func (t *TemplateRequest) createTree(blobs map[string]Blob) (string, error) {
	var treeEntries []Tree

	for file, blob := range blobs {
		treeEntries = append(treeEntries, Tree{
			Path: file,
			Mode: "100644",
			Type: "blob",
			SHA:  blob.SHA,
		})
	}

	data := map[string]interface{}{
		"tree": treeEntries,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees", os.Getenv("ORGANIZATION_NAME"), t.RepoName)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("erro ao criar a árvore: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result["sha"].(string), nil
}

func (t *TemplateRequest) createCommit(treeSHA string) error {
	// Obtém o SHA do último commit na branch especificada
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", os.Getenv("ORGANIZATION_NAME"), t.RepoName, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao obter o SHA do último commit: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	lastCommitSHA := result["object"].(map[string]interface{})["sha"].(string)

	// Cria o novo commit
	data := map[string]interface{}{
		"message": commitMsg,
		"tree":    treeSHA,
		"parents": []string{lastCommitSHA},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits", os.Getenv("ORGANIZATION_NAME"), t.RepoName)

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro ao criar o commit: %s", string(body))
	}

	var commitResult map[string]interface{}
	err = json.Unmarshal(body, &commitResult)
	if err != nil {
		return err
	}

	newCommitSHA := commitResult["sha"].(string)

	// Atualiza o ref da branch para apontar para o novo commit
	data = map[string]interface{}{
		"sha":   newCommitSHA,
		"force": false,
	}

	jsonData, err = json.Marshal(data)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", os.Getenv("ORGANIZATION_NAME"), t.RepoName, branch)

	req, err = http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+os.Getenv("TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao atualizar o ref da branch: %s", string(body))
	}

	return nil
}
