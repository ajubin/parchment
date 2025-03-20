package printer

import "bytes"

// Interface qui abstrait l'impression
type Printer interface {
	Print(buffer bytes.Buffer) error
	TestPrint() error
}
