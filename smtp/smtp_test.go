package smtp_test

import (
	"testing"
	"time"

	"github.com/caarlos0/env"

	. "web-layout/utils/smtp"
)

func TestSMTP(t *testing.T) {
	cfg := SMTPCfg{}
	if err := env.Parse(&cfg); err != nil {
		t.Error(err)
	}

	cfg.Init()

	if err := cfg.SendMail(Mail{
		Subject: " Sausage Log",
		Detail:  "sausage" + time.Now().String(),
	}); err != nil {
		t.Error(err)
		return
	}
}
