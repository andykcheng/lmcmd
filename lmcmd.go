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
	fmt.Print("Enter your key: ")
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
		// File does not exist, ask user for key
		key := getKeyFromUser()
		if err := saveKey(key); err != nil {
			return "", err
		}
		return key, nil
	}

	// File exists, read key
	keyBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(keyBytes)), nil
}

type TogetherBodyMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type TogetherBody struct {
	Model       string                `json:"model"`
	Messages    []TogetherBodyMessage `json:"messages"`
	Temperature float32               `json:"temperature"`
	MaxTokens   int                   `json:"max_tokens"`
}
type TogetherOutput struct {
	Explanation string `json:"explanation"`
	Command     string `json:"command"`
}

func main() {
	key, err := getKey()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	fmt.Println("Your key is:", key)

	// get the user input from args
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
	url := "https://api.together.xyz/v1/chat/completions"
	method := "POST"

	payloadBody := TogetherBody{
		Model: "meta-llama/Llama-3-70b-chat-hf",
		Messages: []TogetherBodyMessage{
			{Role: "system", Content: fmt.Sprintf(`You are a command generator which will generate %s commands according to user requirements. 
			
			Give the command and explanation of the command especially the parameters used.

			Output in a json format like the following:
				{
					"command": "ls -l",
					"explanation": "List files in long format"
				}
				Do not output anything else except the json structure.
			`, osRunning)},
			{Role: "user", Content: command},
		},
		Temperature: 0.8,
		MaxTokens:   500,
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

	//parse only the first choice message
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	outputFromLllm := response.Choices[0].Message.Content
	// parse outputFromLlm to get the json format into a TogetherOutput struct
	var togetherOutput TogetherOutput
	err = json.Unmarshal([]byte(outputFromLllm), &togetherOutput)
	if err != nil {
		fmt.Println(err)
		return
	}
	const (
		ColorReset  = "\033[0m"
		ColorRed    = "\033[31m"
		ColorGreen  = "\033[32m"
		ColorYellow = "\033[33m"
		ColorBlue   = "\033[34m"
		ColorPurple = "\033[35m"
		ColorCyan   = "\033[36m"
		ColorWhite  = "\033[37m"
	)

	fmt.Println(ColorGreen, togetherOutput.Explanation, ColorReset, "\n")
	fmt.Println(ColorBlue, togetherOutput.Command, ColorReset, "\n")
	fmt.Println(ColorGreen, "Command copied to clipboard", ColorReset, "\n")

	err = clipboard.WriteAll(togetherOutput.Command)
	if err != nil {
		fmt.Println("Failed to copy to clipboard:", err)
		return
	}
	fmt.Println("Text copied to clipboard successfully")

}
