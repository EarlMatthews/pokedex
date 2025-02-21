package main
import(
	"bufio"
	"os"
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	// Implementation of cleanInput function
	words := strings.Fields(text)
	return words
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan(){
			break
		}
		input := strings.TrimSpace(scanner.Text())
		words := strings.Fields(strings.ToLower(input))
		if len(words) > 0 {
			fmt.Printf("Your command was: %s\n",words[0])
		}
	}
}