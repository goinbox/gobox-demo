package resource

import (
	"fmt"

	"github.com/goinbox/mysql"

	"gdemo/conf"
)

const (
	dbMySQLKey = "main"
)

func InitMySQL(config *conf.MySQLConf) error {
	mconfig := mysql.NewDefaultConfig(config.User, config.Pass, config.Host, config.Name, config.Port)
	mconfig.ReadTimeout = config.RWTimeout
	mconfig.WriteTimeout = config.RWTimeout

	err := mysql.RegisterDB(dbMySQLKey, mconfig)
	if err != nil {
		return fmt.Errorf("mysql.RegisterDB error: %w", err)
	}

	return nil
}

func MySQLClient() *mysql.Client {
	client, _ := mysql.NewClientFromPool(dbMySQLKey)

	return client
}
