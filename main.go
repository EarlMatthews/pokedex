package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EarlMatthews/pokedexcli/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
}

var cliCommands map[string]cliCommand

type apiResponse struct {
	Results  []locationArea `json:"results"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
}

type locationArea struct {
	Name string `json:"name"`
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range cliCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, cache *pokecache.Cache) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	if data, found := cache.Get(url); found {
        var apiResp apiResponse
        if err := json.Unmarshal(data, &apiResp); err != nil {
            return fmt.Errorf("failed to unmarshal cached data: %v", err)
        }
        return displayLocationAreas(apiResp, cfg)
    }

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Marshal the response before adding to cache
    jsonData, err := json.Marshal(apiResp)
    if err != nil {
        return fmt.Errorf("failed to marshal data for cache: %v", err)
    }
    cache.Add(url, jsonData)

	return displayLocationAreas(apiResp, cfg)
}

func commandMapBack(cfg *config, cache *pokecache.Cache) error {
	url := cfg.Previous
	if url == "" {
		fmt.Println("No previous areas available.")
		return nil
	}

	if data, found := cache.Get(url); found {
        var apiResp apiResponse
        if err := json.Unmarshal(data, &apiResp); err != nil {
            return fmt.Errorf("failed to unmarshal cached data: %v", err)
        }
        return displayLocationAreas(apiResp, cfg)
    }

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Marshal the response before adding to cache
    jsonData, err := json.Marshal(apiResp)
    if err != nil {
        return fmt.Errorf("failed to marshal data for cache: %v", err)
    }
    cache.Add(url, jsonData)

	return displayLocationAreas(apiResp, cfg)
}

func displayLocationAreas(apiResp apiResponse, cfg *config) error {
	if len(apiResp.Results) == 0 {
		fmt.Println("No more locations to display.")
		return nil
	}

	fmt.Println("Location Areas:")
	for _, area := range apiResp.Results {
		fmt.Println("-", area.Name)
	}

	cfg.Next = apiResp.Next
	cfg.Previous = apiResp.Previous
	return nil
}

func main() {
	cliCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next location areas",
			callback: func(cfg *config) error {
				cache := pokecache.NewCache(5 * time.Second)
				return commandMap(cfg, cache)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous location areas",
			callback: func(cfg *config) error {
				cache := pokecache.NewCache(5 * time.Second)
				return commandMapBack(cfg, cache)
			},
		},
	}

	cfg := &config{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		cmd, exists := cliCommands[cmdName]
		if !exists {
			fmt.Printf("Unknown command: %s\n", cmdName)
			continue
		}

		if err := cmd.callback(cfg); err != nil {
			fmt.Printf("Error executing command %s: %v\n", cmdName, err)
		}
	}
}