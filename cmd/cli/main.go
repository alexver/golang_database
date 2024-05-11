package main

import (
	"bufio"
	"fmt"
	"os"

	database "github.com/alexver/golang_database/internal"
)

func main() {

	db, err := database.CreateCLIDatabase()
	if err != nil {
		panic(err)
	}

	displayHelpScreen(db)

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 5; i++ {
		fmt.Print("database cli > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		result, resultErr := db.ProcessQuery(command)
		if resultErr != nil {
			fmt.Printf("Error: %s\n", resultErr.Error())
			displayHelpScreen(db)

			continue
		}

		fmt.Println(result)
	}
}

func displayHelpScreen(db *database.Database) {
	fmt.Print("\n\nTest Database CLI tool\n\nUsage:\n\t<command> [argument]\n\nThe commands are:\n")
	for _, analyzer := range db.GetAnalyzers() {
		fmt.Printf("\t%-15s\t%s\n", analyzer.Name(), analyzer.Description())
	}
	fmt.Print("\n\n")
}
