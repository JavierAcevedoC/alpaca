package main

import (
	"alpaca/src/internal/llm"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	Search string `json:"search"`
}

type IconSource struct {
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

const Clear = "\"Clear\""
const Finished = "\"Finished\""

func main() {
	createMock()
	generateResponse()
}

func generateResponse() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req Request
		err := json.Unmarshal(scanner.Bytes(), &req)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input: %v\n", err)
			continue
		}
		// lets go
		out := llm.AskLLM(req.Search)

		fmt.Println(string(Clear)) //clear pop-launcher ui list
		resp := Response{
			Append: SearchResult{Name: out, Icon: &IconSource{Name: "text-x-generic"}},
		}

		respBytes, _ := json.Marshal(resp)
		fmt.Println(string(respBytes))
		fmt.Println(string(Finished))
	}
}

func createMock() {
	// mockup
	resp := Response{
		Append: SearchResult{Name: "Generating...", Icon: &IconSource{Name: "edit-find"}},
	}

	respBytes, _ := json.Marshal(resp)
	fmt.Println(string(respBytes))
	fmt.Println(string(Finished))
}
