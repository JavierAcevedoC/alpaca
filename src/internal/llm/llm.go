package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func AskLLM(search string) string {
	question := strings.Split(search, "ia") //from the pop-launcher instance.
	req, shouldReturn, returnValue := gotKnowledge(question[1])
	if shouldReturn {
		return returnValue
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error on request:", err)
		return err.Error()
	}

	defer resp.Body.Close()

	return cleanAndFormatting(resp.Body)
}

func gotKnowledge(search string) (*http.Request, bool, string) {
	local_url := "http://localhost:11434/api/generate"

	requestBody := ChatRequest{
		Model:  "phi3:mini",
		Prompt: search,
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		fmt.Println("Error serialization:", err)
		return nil, true, err.Error()
	}

	req, err := http.NewRequest("POST", local_url, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, true, err.Error()
	}

	req.Header.Set("Content-Type", "application/json")
	return req, false, ""
}

func cleanAndFormatting(body io.ReadCloser) string {

	const break_new_line_format = 20
	var cleaned_response string
	var count_for_new_line int = 0

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()

		var response ChatResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue
		}

		if count_for_new_line == break_new_line_format {
			cleaned_response += "\n"
			count_for_new_line = 0
		} else {
			count_for_new_line++
		}
		cleaned_response += response.Response
	}
	return cleaned_response
}
