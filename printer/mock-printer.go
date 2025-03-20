package printer

import (
	"bytes"
	"fmt"
)

type MockPrinter struct{}

// Implémentation pour afficher le texte dans la console (au lieu d'imprimer)
func (p *MockPrinter) Print(buffer bytes.Buffer) error {
	fmt.Println("🔹 MOCK PRINTER OUTPUT 🔹")
	fmt.Println(buffer.String())
	fmt.Println("🔹 FIN 🔹")
	return nil
}

func (p *MockPrinter) TestPrint() error {
	return nil
}
