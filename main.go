package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/ajubin/parchment/api"
	p "github.com/ajubin/parchment/printer"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "the server address, default ':8080'")

	flag.Parse()

	// TODO: improve this design, pass port name as flag maybe
	var printer p.Printer
	if strings.ToLower(strings.TrimSpace(getEnv("MODE", "dev"))) == "prod" {

		// En mode production, on utilise l'imprimante série
		fmt.Println("Utilisation de l'imprimante en serie")
		printer = &p.SerialPrinter{PortName: "/dev/ttyS0", BaudRate: 9600}
	} else {
		// En mode dev, on utilise le mock
		fmt.Println("Utilisation de l'imprimante en mock")
		printer = &p.MockPrinter{}

	}

	server := api.NewServer(*listenAddr, printer)

	fmt.Printf("Serveur démarré sur %s", *listenAddr)
	log.Fatal(server.Start())

}

func getEnv(key, defaultValue string) string {
	if value, exists := syscall.Getenv(key); exists {
		return value
	}
	return defaultValue
}
