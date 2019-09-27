package mysql_test

import (
	"testing"

	. "web-layout/utils/mysql"
)

func TestModelBuilder_Build(t *testing.T) {
	m := ModelBuilder{
		User: "root",
		Pwd:  "12345678",
		DB:   "sausage_shooter",
		Tables: []string{
			`player`,
			`nodes`,
			`openid_mapping`,
			`player_global_mail`,
			`player_invitation`,
			`pvp_result`,
			`single_mail`,
		},
	}
	m.Build()
}
