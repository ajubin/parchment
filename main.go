package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"go.bug.st/serial"
)

// Structure du payload attendu
type PrintRequest struct {
	Content string `json:"content"`
}

// Interface qui abstrait l'impression
type Printer interface {
	Print(buffer bytes.Buffer) error
	TestPrint() error
}

type MockPrinter struct{}

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

	fmt.Printf("Envoyé %v octets à l'imprimante\n", n)

	// Pause pour éviter d'envoyer trop rapidement
	time.Sleep(2 * time.Second)
	return nil
}

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

func sanitizeMarkdown(input string) string {
	p := bluemonday.StrictPolicy() // 🚀 **Only allows safe text, no HTML**
	return p.Sanitize(input)
}

func handlePrint(p Printer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		// 🚀 **Limit the request size to 1MB (or adjust as needed)**
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1) // 1MB limit

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Impossible de lire la requête: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req PrintRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Format JSON invalide", http.StatusBadRequest)
			return
		}

		// ✅ **Sanitize Markdown before parsing**
		cleanContent := sanitizeMarkdown(req.Content)

		// Parser le Markdown et générer le buffer
		buffer := parseMarkup(cleanContent)

		// Envoyer au printer
		err = p.Print(buffer)
		if err != nil {
			http.Error(w, "Erreur lors de l'impression", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	}
}

func handleTestPrint(p Printer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		// 🚀 **Limit the request size to 1MB (or adjust as needed)**
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

		p.TestPrint()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	}
}

// 🖨️ **Parse le markup et construit le buffer à imprimer**
func parseMarkup(content string) bytes.Buffer {
	var buffer bytes.Buffer
	// buffer.Write(resetPrinter) // Réinitialisation

	lines := strings.Split(content, "\n")

	// TODO: add markup parsing (maybe with markdown)
	for _, line := range lines {
		buffer.WriteString(line + "\n") // Texte normal

	}
	// Ajouter un saut de ligne final
	return buffer
}

func secureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func main() {
	const PORT = 8080

	var printer Printer

	if strings.ToLower(strings.TrimSpace(getEnv("MODE", "dev"))) == "prod" {

		// En mode production, on utilise l'imprimante série
		fmt.Println("Utilisation de l'imprimante en serie")
		printer = &SerialPrinter{PortName: "/dev/ttyS0", BaudRate: 9600}
	} else {
		// En mode dev, on utilise le mock
		fmt.Println("Utilisation de l'imprimante en mock")
		printer = &MockPrinter{}
	}

	// Créer un multiplexer personnalisé
	mux := http.NewServeMux()
	mux.HandleFunc("/print", handlePrint(printer))
	mux.HandleFunc("/test-print", handleTestPrint(printer))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", PORT),
		ReadTimeout:  10 * time.Second, // ⏳ Prevents slow request attacks
		WriteTimeout: 10 * time.Second, // ⏳ Prevents slow response attacks
		IdleTimeout:  120 * time.Second,
		Handler:      secureMiddleware(mux), // 🔒 Security Middleware
	}

	fmt.Printf("Serveur démarré sur http://localhost:%d", PORT)
	log.Fatal(server.ListenAndServe())

}

func getEnv(key, defaultValue string) string {
	if value, exists := syscall.Getenv(key); exists {
		return value
	}
	return defaultValue
}

// rsync -avz -e ssh ./ pi@192.168.1.97:~/printer-serial
