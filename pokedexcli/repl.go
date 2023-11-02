package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/commands"
	"pokedexcli/internal/config"
	"strings"
)

func tokenizeInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}

func repl(config *config.ApplicationData) {
	scanner := bufio.NewScanner(os.Stdin)
	menu := commands.NewMenu(config)

	for {
		fmt.Printf("prompt> ")
		scanner.Scan()
		text := scanner.Text()

		fields := tokenizeInput(text)
		if len(fields) == 0 {
			continue
		}

		if command, ok := menu.MenuCommands[fields[0]]; ok {
			if err := command.CommandCallback(fields); err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println("Unrecognised command")
			continue
		}
	}
}
