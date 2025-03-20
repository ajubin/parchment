package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ajubin/parchment/api"
	p "github.com/ajubin/parchment/printer"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "the server address, default ':8080'")
	serialPort := flag.String("serialPort", "", "path to the serial port of the printer, eg: /dev/ttyS0. Will use mock printer if nothing provided")
	flag.Parse()

	// TODO: improve this design, pass port name as flag maybe
	var printer p.Printer
	if *serialPort == "" {
		fmt.Println("Utilisation de l'imprimante en mock")
		printer = &p.MockPrinter{}
	} else {
		fmt.Println("Utilisation de l'imprimante en serie sur le port", *serialPort)
		printer = &p.SerialPrinter{PortName: *serialPort, BaudRate: 9600}

	}

	server := api.NewServer(*listenAddr, printer)

	fmt.Printf("Serveur démarré sur %s", *listenAddr)
	log.Fatal(server.Start())

}
