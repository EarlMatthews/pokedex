package main
import(
	"bufio"
	"os"
	"fmt"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cliCommands map[string]cliCommand

func cleanInput(text string) []string {
	// Implementation of cleanInput function
	words := strings.Fields(text)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range cliCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}



func main(){

	cliCommands = map[string]cliCommand{
	"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
	"help": {name: "help", description: "Displays a help message", callback: commandHelp},
	}

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
				cmd.callback()
			} else {
				fmt.Printf("Unknown command: %s\n", words[0])
			}
		} 
	}
}