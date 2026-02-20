package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	cyan  = "\033[36m"
	lime  = "\033[92m"
	white = "\033[37m"
	reset = "\033[0m"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("lflist error: failed to read directory '%s': %v\n", path, err)
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)

		fi, err := os.Lstat(fullPath)
		if err != nil {
			fmt.Printf("lflist warning: cannot access '%s': %v\n", fullPath, err)
			continue
		}

		switch {
		case fi.Mode()&os.ModeSymlink != 0:
			fmt.Println(cyan + name + reset)
		case fi.IsDir():
			fmt.Println(lime + name + reset)
		default:
			fmt.Println(white + name + reset)
		}
	}
}
