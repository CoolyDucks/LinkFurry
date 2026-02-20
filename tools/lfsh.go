package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	red    = "\033[31m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	lime   = "\033[92m"
	white  = "\033[37m"
	reset  = "\033[0m"

	configFileName = ".lfsh_config.json"
)

type Config struct {
	Name string `json:"name"`
}

func main() {
	cfg := loadOrCreateConfig()
	reader := bufio.NewReader(os.Stdin)

	for {
		printPrompt(cfg.Name)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := parseInput(input)
		cmd := args[0]

		switch cmd {
		case "exit":
			return
		case "cd":
			changeDir(args)
		case "clear":
			fmt.Print("\033[H\033[2J")
		case "ls":
			listDir(args)
		default:
			runExternal(args)
		}
	}
}

func loadOrCreateConfig() Config {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, configFileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("What is your name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		cfg := Config{Name: name}
		saveConfig(path, cfg)
		return cfg
	}

	file, _ := os.ReadFile(path)
	var cfg Config
	json.Unmarshal(file, &cfg)
	return cfg
}

func saveConfig(path string, cfg Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(path, data, 0644)
}

func printPrompt(name string) {
	wd, err := os.Getwd()
	if err != nil {
		wd = "unknown"
	}

	fmt.Printf("%s%s@linkfurrylinux:%s%s - %s#%s ",
		red, name, wd, reset,
		blue, reset)
}

func parseInput(input string) []string {
	return strings.Fields(input)
}

func changeDir(args []string) {
	if len(args) < 2 {
		return
	}
	err := os.Chdir(args[1])
	if err != nil {
		fmt.Println("cd error:", err)
	}
}

func listDir(args []string) {
	path := "."
	if len(args) > 1 {
		path = args[1]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("ls error:", err)
		return
	}

	for _, entry := range entries {
		info, _ := entry.Info()
		mode := info.Mode()
		name := entry.Name()

		if mode&os.ModeSymlink != 0 {
			fmt.Println(cyan + name + reset)
		} else if mode.IsDir() {
			fmt.Println(lime + name + reset)
		} else {
			fmt.Println(white + name + reset)
		}
	}
}

func runExternal(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}
