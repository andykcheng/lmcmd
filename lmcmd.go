package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
)

const keyFileName = ".lmcmd.config"

func getKeyFromUser() string {
	fmt.Print("Enter your OpenAI API key: ")
	reader := bufio.NewReader(os.Stdin)
	key, _ := reader.ReadString('\n')
	return strings.TrimSpace(key)
}

// saveKey saves the key to a file in the user's home directory.
func saveKey(key string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	keyFilePath := filepath.Join(homeDir, keyFileName)
	return os.WriteFile(keyFilePath, []byte(key), 0600)
}

func getKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	keyFilePath := filepath.Join(homeDir, keyFileName)

	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		key := getKeyFromUser()
		if err := saveKey(key); err != nil {
			return "", err
		}
		return key, nil
	}

	keyBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(keyBytes)), nil
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIBody struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float32         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	TopP        float32         `json:"top_p"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	key, err := getKey()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	// fmt.Println("Your key is:", key)

	// Get the user input from args
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide a command to generate")
		return
	}
	command := strings.Join(args, " ")

	osRunning := "mac"
	if runtime.GOOS == "windows" {
		fmt.Println("| Windows is not supported yet")
		return
	} else if runtime.GOOS == "linux" {
		osRunning = "linux"
	}
	fmt.Printf("Generating command for %s\n", osRunning)
	url := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	payloadBody := OpenAIBody{
		Model: "gpt-3.5-turbo",
		Messages: []OpenAIMessage{
			{Role: "system", Content: fmt.Sprintf(`You are an AI model that generates shell commands for %s operating system.
			You will only output a JSON format with the command and explanation. Example:
			{
				"command": "ls -l",
				"explanation": "List files with ls command with -l flag for long output format"
			}
			`, osRunning)},
			{Role: "user", Content: command},
		},
		Temperature: 0.1,
		MaxTokens:   500,
		TopP:        1.0,
	}

	payloadJSON, err := json.Marshal(payloadBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := strings.NewReader(string(payloadJSON))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the raw response body for troubleshooting
	// fmt.Println("Raw response body:", string(body))

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return
	}

	if len(response.Choices) > 0 {
		//parse the response.Choices[0].Message.Content and put command into a variable and put explanation into another variable
		// Parse the JSON content
		type Content struct {
			Command     string `json:"command"`
			Explanation string `json:"explanation"`
		}
		var content Content
		if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &content); err != nil {
			fmt.Println(err)
			return
		}
		if content.Command == "" || content.Explanation == "" {
			fmt.Println("Error: JSON content does not have required attributes 'command' and 'explanation'")
			return
		}

		// Assign to variables
		command := content.Command
		explanation := content.Explanation

		// Step 3: Assign to variables
		fmt.Println("Command:", command)
		fmt.Println("Explanation:", explanation)

		if err := clipboard.WriteAll(command); err != nil {
			fmt.Println("Failed to copy command to clipboard:", err)
			return
		}
		fmt.Println("Command copied to clipboard.")
	} else {
		fmt.Println("No command generated.")
	}

}
