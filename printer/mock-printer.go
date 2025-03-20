package printer

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

type MockPrinter struct{}

// Implémentation pour afficher le texte dans la console (au lieu d'imprimer)
func (p *MockPrinter) Print(buffer bytes.Buffer) error {
	log.Println("🔹 MOCK PRINTER OUTPUT 🔹")
	log.Println(buffer.String())
	log.Println("🔹 FIN 🔹")
	return nil
}

func (p *MockPrinter) TestPrint() error {
	return nil
}
