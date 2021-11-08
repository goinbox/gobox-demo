package resource

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"gdemo/conf"
)

const (
	dbKey = "main"
)

func InitMySQL(config *conf.MySQLConf) error {
	mconfig := mysql.NewDefaultConfig(config.User, config.Pass, config.Host, config.Name, config.Port)
	mconfig.ReadTimeout = config.RWTimeout
	mconfig.WriteTimeout = config.RWTimeout

	err := mysql.RegisterDB(dbKey, mconfig)
	if err != nil {
		return fmt.Errorf("mysql.RegisterDB error: %w", err)
	}

	return nil
}

func MySQLClient(logger golog.Logger) *mysql.Client {
	client, _ := mysql.NewClientFromPool(dbKey, logger)

	return client
}
