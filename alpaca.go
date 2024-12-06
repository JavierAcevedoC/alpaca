package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	Search string `json:"search"`
}

type IconSource struct {
	// Define la estructura de IconSource según tu implementación en Rust
	Name string `json:"Name"`
}

type WindowGeneration struct {
	Generation int `json:"generation"` // Asumiendo que Generation es un entero
	Indice     int `json:"indice"`     // Asumiendo que Indice es un entero
}

type SearchResult struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Keywords    []string          `json:"keywords"`
	Icon        *IconSource       `json:"icon"`
	Exec        *string           `json:"exec"`
	Window      *WindowGeneration `json:"window"`
}

type Response struct {
	Append SearchResult `json:"Append"`
	Fill   string       `json:"Fill,omitempty"`
}

type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func askLLM(search string) string {
	local_url := "http://localhost:11434/api/generate"
	// Crear el cuerpo de la solicitud
	requestBody := ChatRequest{
		Model:  "phi3:mini",
		Prompt: search,
	}

	// Serializar a JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error serializando la solicitud:", err)
		return "Error serialization"
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", local_url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return "Error creating request"
	}
	req.Header.Set("Content-Type", "application/json")

	// Hacer la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error realizando la solicitud:", err)
		return "Error on request"
	}
	defer resp.Body.Close()

	var final_response string
	// Procesar la respuesta línea por línea
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println("Línea recibida:", line)

		var response ChatResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			// fmt.Println("Error parseando la línea:", err)
			continue
		}
		final_response += response.Response
		// fmt.Printf("Respuesta parseada: %+v\n", response)
	}
	// fmt.Println("Respuesta del servidor Ollama:", final_response)
	return final_response

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req Request
		err := json.Unmarshal(scanner.Bytes(), &req)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input: %v\n", err)
			continue
		}

		out := askLLM(req.Search)

		resp := Response{
			Append: SearchResult{Name: out, Icon: &IconSource{Name: "text-x-generic"}},
		}

		respBytes, _ := json.Marshal(resp)
		fmt.Println(string(respBytes))
		fmt.Println(string("\"Finished\""))
	}
}
