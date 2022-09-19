package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool
	Loc       string
}
type Server struct {
	Host string
	Port int
}

var Config = &struct {
	Mysql  Mysql  `toml:"mysql"`
	Server Server `toml:"server"`
}{}

func init() {
	_, err := toml.DecodeFile("./conf/conf.toml", &Config)
	if err != nil {
		panic(err)
	}
}

func MysqlConnectString() string {
	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Config.Mysql.Username, Config.Mysql.Password, Config.Mysql.Host, Config.Mysql.Port, Config.Mysql.Database, Config.Mysql.Charset, Config.Mysql.ParseTime, Config.Mysql.Loc)
	return str
}
