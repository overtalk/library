package mysql

import (
	"database/sql"
	"time"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

// Config implement MySQL config
type Config struct {
	User    string `json:"mysql_user" env:"MYSQL_USER"`
	Pwd     string `json:"mysql_pwd" env:"MYSQL_PWD"`
	Host    string `json:"mysql_host" env:"MYSQL_HOST" envDefault:"127.0.0.1"`
	DB      string `json:"mysql_db" env:"MYSQL_DB"`
	Port    int    `json:"mysql_port" env:"MYSQL_PORT" envDefault:"3306"`
	MaxConn int    `json:"mysql_max_conn" env:"MYSQL_MAX_CONN" envDefault:"1000"`
}

// Connect is the constructor
func (c Config) Connect() (*sql.DB, error) {
	return manager.New(c.DB, c.User, c.Pwd, c.Host).Set(
		manager.SetCharset("utf8"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second),
	).Port(c.Port).Open(true)
}
