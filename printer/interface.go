package printer

// Interface qui abstrait l'impression
type Printer interface {
	Print(content string) error
	TestPrint() error
}
