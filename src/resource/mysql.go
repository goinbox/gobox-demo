package resource

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"gdemo/conf"
)

const (
	DBKey = "main"
)

func InitMySQL(config *conf.MySQLConf) error {
	err := mysql.RegisterDB(DBKey,
		mysql.NewDefaultConfig(config.User, config.Pass, config.Host, config.Name, config.Port))
	if err != nil {
		return fmt.Errorf("mysql.RegisterDB error: %w", err)
	}

	return nil
}

func MySQLClient(logger golog.Logger) *mysql.Client {
	client, _ := mysql.NewClientFromPool(DBKey, logger)

	return client
}
