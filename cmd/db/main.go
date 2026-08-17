package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/HeyReyHR/keyvaluedb/internal"
	"github.com/HeyReyHR/keyvaluedb/internal/compute/parser"
	"github.com/HeyReyHR/keyvaluedb/internal/storage"
	"github.com/HeyReyHR/keyvaluedb/internal/storage/engine"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	p := parser.NewParser(l)
	e := engine.NewEngine()
	s := storage.NewStorage(e)

	db, err := internal.NewDatabase(p, s, l)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("[db] ")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Print("[db] ")
			continue
		}
		if line == "exit" {
			break
		}
		result := db.HandleQuery(line)
		fmt.Println("[db]", result)
		fmt.Print("[db] ")
	}
}
