package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/temminks/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var locationOffset int

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 map locations",
			callback:    commandMapb,
		},
	}
}

func getLocations() error {
	locations, err := pokeapi.GetLocations(locationOffset)
	if err != nil {
		return err
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMap() error {
	err := getLocations()
	if err != nil {
		return err
	}

	locationOffset += 20
	return nil
}

func commandMapb() error {
	locationOffset -= 20
	if locationOffset < 0 {
		locationOffset = 0
		return errors.New("Cannot go further back...")
	}
	err := getLocations()
	if err != nil {
		return err
	}
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		text := reader.Text()
		if len(text) == 0 {
			continue
		}
		command, ok := getCommands()[text]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
