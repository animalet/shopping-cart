package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"shopping-cart/client"
	"shopping-cart/server"
)

func main() {
	flag.NewFlagSet("Help", flag.ExitOnError)
	isHelp := flag.Bool("help", true, "help")

	flag.NewFlagSet("Client mode", flag.ExitOnError)
	isClient := flag.Bool("client", false, "client mode")
	endpoint := flag.String("endpoint", "http://localhost:8000", "backend server address")

	flag.NewFlagSet("Server mode", flag.ExitOnError)
	isServer := flag.Bool("server", false, "server mode")
	flag.Parse()

	switch {
	case *isClient:
		cmdExecutor := client.NewCommandExecutor(*endpoint)
		codesString, reader := cmdExecutor.StringCodesForConsole(), bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("\nInput product codes to add to the cart: %s\n\"r\" to remove the cart, \"q\" to quit the client.\n", *codesString)
			if text, err := reader.ReadString('\n'); err != nil {
				fmt.Printf("An error occured reading your input: %s", err)
			} else {
				processInput(cmdExecutor, text)
			}
		}
	case *isServer:
		server.StartServer()
	case *isHelp:
		flag.Usage()
	}
}

func processInput(commandExecutor client.CommandExecutor, text string) {
	defer recoverFromPanic()
	commandExecutor.Execute(text)
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Printf("Error performing operation: %s\n", r)
	}
}
