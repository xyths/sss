package message

import "github.com/xyths/sss/stake"

type Message struct {
	Date     string       `json:"date"`
	To       string       `json:"to"`
	Report   stake.Report `json:"report"`
	Attempts int          `json:"attempts"`
}
