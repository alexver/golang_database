package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexver/golang_database/internal/compute"
	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/network"
)

func main() {

	var net = flag.String("network", "tcp", "Network type")
	var host = flag.String("host", "127.0.0.1", "Server IP address")
	var port = flag.Int("port", 8080, "Server port number")

	flag.Parse()

	client := network.CreateClient(*net, *host, *port)

	analyzers := []compute.AnalyzerInterface{
		analyzer.NewGet(),
		analyzer.NewSet(),
		analyzer.NewDel(),
		analyzer.NewExit(),
		analyzer.NewPing(),
	}

	displayHelpScreen(analyzers)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("database > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Read user input error: %s", err)
		}

		command = strings.Trim(command, " \t\n")
		if command == analyzer.COMMAND_EXIT_1 || command == analyzer.COMMAND_EXIT_2 {
			break
		}

		result := client.CallClient(command)

		fmt.Println(result)
	}
}

func displayHelpScreen(analyzers []compute.AnalyzerInterface) {
	fmt.Print("\n\nTest Database client\n\nUsage:\n\t<command> [argument]\n\nThe commands are:\n")
	for _, analyzer := range analyzers {
		fmt.Printf("\t%-15s\t%s\n", analyzer.Name(), analyzer.Description())
	}
	fmt.Print("\n\n")
}
