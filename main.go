package main
import(
	"bufio"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
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
	// Implementation of cleanInput function
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

func commandMap(cfg *config) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
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

func commandMapBack(cfg *config) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?prev"
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

func main(){

	cliCommands = map[string]cliCommand{
	"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
	"help": {name: "help", description: "Displays a help message", callback: commandHelp},
	"map": {name: "map", description: "Displays location areas", callback: commandMap},
	"mapb": {name: "mapb", description: "Displays location areas", callback: commandMapBack},
	}

	cfg := &config{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan(){
			break
		}
		input := strings.TrimSpace(scanner.Text())
		words := cleanInput(input)
		if len(words) > 0 {
			if cmd, exists := cliCommands[words[0]]; exists {
				cmd.callback(cfg)
			} else {
				fmt.Printf("Unknown command: %s\n", words[0])
			}
		} 
	}
}