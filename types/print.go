package types

type PrintRequest struct {
	Content string `json:"content"`
}

func ValidatePrintRequet(p *PrintRequest) bool { return true }
