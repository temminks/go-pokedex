package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/temminks/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
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
		"explore": {
			name:        "explore",
			description: "Display the names of possible Pokemon encounters",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
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

func commandCatch(args string) error {
	if len(strings.Split(args, " ")) > 1 {
		return errors.New(fmt.Sprintf("only one Pokemon can be caught at a time: `%s` cannot be caught. Use a dash for Pokemon with two-word names, e.g. mime-jr.", args))
	}
	if strings.Trim(args, " ") == "" {
		return errors.New("`catch` requires the Pokemon to catch as an argument.")
	}

	pokemon, err := pokeapi.GetPokemonDetails(args)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	return nil
}

func commandExplore(args string) error {
	if len(strings.Split(args, " ")) > 1 {
		return errors.New(fmt.Sprintf("only one location can be explored at a time: `%s` is an invalid location.", args))
	}
	if strings.Trim(args, " ") == "" {
		return errors.New("`explore` takes a location as argument.")
	}

	location, err := pokeapi.GetLocation(args)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandMap(args string) error {
	if strings.Trim(args, " ") != "" {
		return errors.New(fmt.Sprintf("`map` does not take args `%s`.", args))
	}

	err := getLocations()
	if err != nil {
		return err
	}

	locationOffset += 20
	return nil
}

func commandMapb(args string) error {
	if strings.Trim(args, " ") != "" {
		return errors.New(fmt.Sprintf("`mapb` does not take args `%s`.", args))
	}

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

func commandHelp(args string) error {
	if strings.Trim(args, " ") != "" {
		return errors.New(fmt.Sprintf("`help` does not take args `%s`.", args))
	}

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	fmt.Println()
	return nil
}

func commandExit(args string) error {
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

		commandPart, argsPart, _ := strings.Cut(text, " ")

		command, ok := getCommands()[commandPart]
		if ok {
			err := command.callback(argsPart)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
