package printer

import (
	log "github.com/sirupsen/logrus"
)

type MockPrinter struct{}

// Implémentation pour afficher le texte dans la console (au lieu d'imprimer)
func (p *MockPrinter) Print(content string) error {
	log.Println("🔹 MOCK PRINTER OUTPUT 🔹")
	log.Println(content)
	log.Println("🔹 FIN 🔹")
	return nil
}

func (p *MockPrinter) TestPrint() error {
	return nil
}
