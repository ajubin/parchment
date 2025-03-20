package main

import (
	"flag"

	"github.com/ajubin/parchment/api"
	p "github.com/ajubin/parchment/printer"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Note, config can be improved with this pattern https://www.kaznacheev.me/posts/en/clean-way-pass-configs-go-application/
	listenAddr := flag.String("listenaddr", ":8080", "the server address, default ':8080'")
	serialPort := flag.String("serialPort", "", "path to the serial port of the printer, eg: /dev/ttyS0. Will use mock printer if nothing provided")
	apiUser := flag.String("apiUser", "admin", "the user of basic auth to access protected routes, defaults: admin")
	apiToken := flag.String("apiToken", "admin", "the password of basic auth to access protected routes, defaults: admin")
	flag.Parse()

	if *apiUser == "admin" && *apiToken == "admin" {
		log.Warn("No basic-auth info provided. In production, provide a --apiUser and --apiToken to protect routes")
	}

	// TODO: improve this design, pass port name as flag maybe
	var printer p.Printer
	if *serialPort == "" {
		log.Println("Utilisation de l'imprimante en mock")
		printer = &p.MockPrinter{}
	} else {
		log.Println("Utilisation de l'imprimante en serie sur le port", *serialPort)
		printer = &p.SerialPrinter{PortName: *serialPort, BaudRate: 9600}

	}

	server := api.NewServer(*listenAddr, printer, *api.NewBasicAuth(*apiUser, *apiToken))

	log.Printf("Serveur démarré sur %s", *listenAddr)
	log.Fatal(server.Start())

}
