package printer

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type SerialPrinter struct {
	PortName string
	BaudRate int
}

// Implémentation pour envoyer sur l'imprimante série
func (p *SerialPrinter) Print(buffer bytes.Buffer) error {
	mode := &serial.Mode{BaudRate: p.BaudRate}
	port, err := serial.Open(p.PortName, mode)
	if err != nil {
		log.Println("Erreur ouverture port série:", err)
		return err
	}
	defer port.Close()

	_, err = port.Write(buffer.Bytes())
	if err != nil {
		log.Println("Erreur envoi imprimante:", err)
		return err
	}

	log.Println("Texte envoyé à l'imprimante")
	return nil
}

func (p *SerialPrinter) TestPrint() error {
	const serialPort = "/dev/ttyS0" // Change si nécessaire
	// Commandes ESC/POS pour la taille du texte
	var (
		resetPrinter  = []byte{0x1B, 0x40}       // Réinitialiser imprimante
		setSmallText  = []byte{0x1B, 0x21, 0x00} // Texte Petit (1x)
		setMediumText = []byte{0x1B, 0x21, 0x10} // Texte Moyen (2x hauteur)
		setLargeText  = []byte{0x1B, 0x21, 0x20} // Texte Large (2x largeur)
		setHugeText   = []byte{0x1B, 0x21, 0x30} // Texte Très Grand (2x hauteur et largeur)
		setBoldOn     = []byte{0x1B, 0x45, 0x01} // Activer le texte en gras
		setBoldOff    = []byte{0x1B, 0x45, 0x00} // Désactiver le texte en gras
		newLine       = []byte{0x0A}             // Saut de ligne
	)

	// Jeux de caractères ESC/POS
	var charsets = [][]byte{
		{0x1B, 0x74, 0x00}, // Charset 0 - Standard
		{0x1B, 0x74, 0x01}, // Charset 1 - Alternative
		{0x1B, 0x74, 0x02}, // Charset 2 - Spécial
		{0x1B, 0x74, 0x03}, // Charset 3 - Autre
		{0x1B, 0x74, 0x04}, // Charset 4 - Japonais (Exemple)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open(serialPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	var buffer bytes.Buffer

	// Réinitialisation
	buffer.Write(resetPrinter)
	buffer.Write(newLine)

	// 1️⃣ Test des tailles de texte
	buffer.Write(setSmallText)
	buffer.WriteString("Texte Petit\n")
	buffer.Write(setMediumText)
	buffer.WriteString("Texte Moyen\n")
	buffer.Write(setLargeText)
	buffer.WriteString("Texte Large\n")
	buffer.Write(setHugeText)
	buffer.WriteString("Texte Très Grand\n")

	// 2️⃣ Test du texte en gras
	buffer.Write(setBoldOn)
	buffer.WriteString("Texte en Gras\n")
	buffer.Write(setBoldOff)
	buffer.WriteString("Texte Normal\n")

	// 3️⃣ Test des jeux de caractères
	for i, charset := range charsets {
		buffer.Write(charset)
		buffer.WriteString(fmt.Sprintf("Charset %d : ABC123\n", i))
	}

	// 4️⃣ Ajout de nouvelles lignes
	buffer.Write(newLine)
	buffer.Write(newLine)

	// // **Boucle pour tester tous les charsets de 0x00 à 0xFF**
	// buffer.Write(resetPrinter)
	// for i := 0x00; i <= 0xFF; i++ {
	// 	// Sélection du charset ESC t n
	// 	charsetCmd := []byte{0x1B, 0x74, byte(i)}
	// 	buffer.Write(charsetCmd)

	// 	// Impression du numéro du charset et d'un texte test
	// 	buffer.WriteString(fmt.Sprintf("Charset %02X : ABC123 éèêçö ñüä\n", i))
	// 	buffer.Write(newLine)

	// 	// Pause pour éviter d'envoyer trop rapidement
	// 	time.Sleep(100 * time.Millisecond)
	// }

	// Envoi au port série
	n, err := port.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Envoyé %v octets à l'imprimante\n", n)

	// Pause pour éviter d'envoyer trop rapidement
	time.Sleep(1 * time.Second)
	return nil
}
