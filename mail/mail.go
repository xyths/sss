package mail

import (
	"github.com/xyths/sss/cmd/utils"
	"log"
)

func Mail(config *utils.MailConfig, date string) error {
	log.Printf("Mail report for date %s ...", date)
	return nil
}


