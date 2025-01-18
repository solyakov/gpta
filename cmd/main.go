package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

const (
	apiURL        = "https://api.openai.com/v1/chat/completions"
	maxOutputSize = 4096
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type Options struct {
	Key         string `short:"k" long:"key" env:"OPENAI_API_KEY" description:"OpenAI API key" required:"true"`
	Task        string `short:"t" long:"task" env:"GPTA_TASK" description:"Task to perform"`
	Interactive bool   `short:"i" long:"interactive" description:"Interactive mode (ask for confirmation before executing commands)"`
	Shell       string `short:"s" long:"shell" env:"GPTA_SHELL" description:"Shell to use for executing commands" default:"/bin/sh"`
	Model       string `short:"m" long:"model" env:"GPTA_MODEL" description:"GPT model to use" default:"gpt-4o"`
	Verbose     bool   `short:"v" long:"verbose" description:"Verbose output"`
	Config      string `short:"c" long:"config" env:"GPTA_CONFIG" description:"Configuration file" default:"~/gpta.system"`
}

func execute(input, shell string) string {
	cmd := exec.Command(shell, "-c", input)
	res, err := cmd.CombinedOutput()

	if err != nil {
		res = append(res, []byte(fmt.Sprintf("Error: %s\n", err.Error()))...)
	}

	if len(res) == 0 {
		res = []byte("No output\n")
	}

	if len(res) > maxOutputSize {
		t := []byte("Output truncated\n")
		res = append(res[:maxOutputSize-len(t)], t...)
	}

	return string(res)
}

func confirm(command string) bool {
	f, err := os.Open("/dev/tty")
	if err != nil {
		log.Fatalf("Unable to open /dev/tty for interactive confirmation: %s", err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	fmt.Printf("Execute '%s' [Y/n]: ", command)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading user input: %s", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes" || response == ""
}

func request(apiURL, apiKey string, requestBody ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	var chatResponse ChatResponse
	if err = json.Unmarshal(body, &chatResponse); err != nil {
		return nil, fmt.Errorf("error decoding response JSON: %w", err)
	}

	return &chatResponse, nil
}

func main() {
	var opts Options
	args, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		}
		log.Fatalf("Error parsing flags: %s", err)
	}

	task := opts.Task
	files := args

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Error stating stdin: %s", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		files = append(files, "/dev/stdin")
	}

	for _, fileName := range files {
		content, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Error reading file '%s': %s", fileName, err)
		}
		task += fmt.Sprintf("\nThe following input was provided via file '%s':\n%s", fileName, string(content))
	}

	if strings.TrimSpace(task) == "" {
		log.Fatal("Please provide a task for the agent to perform.")
	}

	messages := []Message{
		{Role: "system", Content: `
You are a console application connected to a shell. Your purpose is to accurately and efficiently fulfill the task given to you by the user.

You may only produce two types of messages:

1. shell:<command>
  - Executes "` + opts.Shell + ` -c <command>" and returns its output.
  - The output is shown to both you and the user. Do NOT duplicate this output to the user.
  - The output is limited to ` + strconv.Itoa(maxOutputSize) + ` bytes; if it exceeds this limit, it is truncated.
  - If the output is truncated, empty, or an error occurs, an appropriate message will appear in the output.
  - Example: "shell:ls"

2. exit:<code>
  - Terminates the application with the specified exit code.
  - Use "exit:0" for success or "exit:<non-zero>" for failure.
  - Do NOT use "shell:exit <code>" to terminate.
  - Only terminate when the task is complete.

Additional Rules and Constraints:
  - BEFORE performing the task, read "` + opts.Config + `" for additional instructions or context.
  - Do NOT ask the user for input; no interactive capabilities are available.
  - ALL responses MUST be in one of the two formats above: "shell:<command>" or "exit:<code>".
  - Do NOT produce any other text outside these commands. Do NOT explain your reasoning to the user.
  - If you need to display a message to the user, use a shell command such as "shell:echo <message>".
  - Do NOT combine multiple commands in a single response. Respond with ONE command per message, either:
    - "shell:<command>"
    - "exit:<code>"
  - Strictly follow these formats and constraints.
`},
		{Role: "user", Content: task},
	}

	for {
		requestBody := ChatRequest{
			Model:    opts.Model,
			Messages: messages,
		}

		response, err := request(apiURL, opts.Key, requestBody)
		if err != nil {
			log.Fatalf("Error sending API request: %s", err)
		}

		if len(response.Choices) == 0 {
			log.Fatalf("No choices received in response.")
		}

		content := response.Choices[0].Message.Content

		messages = append(messages, Message{Role: "assistant", Content: content})

		switch {

		case strings.HasPrefix(content, "exit:"):
			code := strings.TrimPrefix(content, "exit:")
			code = strings.TrimSpace(code)
			n, err := strconv.Atoi(code)
			if err != nil {
				log.Fatalf("Invalid exit code: %s", code)
			}
			if opts.Verbose {
				log.Printf("Exiting with code: %d", n)
			}
			os.Exit(n)

		case strings.HasPrefix(content, "shell:"):
			command := strings.TrimPrefix(content, "shell:")
			command = strings.TrimSpace(command)
			if opts.Interactive && !confirm(command) {
				log.Fatalf("Aborted by user.")
			}
			if opts.Verbose {
				log.Printf("Executing: %s", command)
			}
			output := execute(command, opts.Shell)
			fmt.Print(output)
			messages = append(messages, Message{Role: "user", Content: output})
			continue

		default:
			log.Fatalf("Invalid response: %s", content)
		}
	}
}
