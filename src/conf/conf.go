package conf

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/BurntSushi/toml"
)

func Init(confDir string) error {
	err := parseServerConf(confDir)
	if err != nil {
		return err
	}

	ServerConf.Hostname, _ = os.Hostname()
	curUser, _ := user.Current()
	ServerConf.Username = curUser.Username

	ServerConf.Redis.ConnectTimeout = time.Duration(ServerConf.Redis.ConnectTimeoutSeconds) * time.Second
	ServerConf.Redis.RWTimeout = time.Duration(ServerConf.Redis.RWTimeoutSeconds) * time.Second
	ServerConf.Redis.PoolKeepAliveInterval = time.Duration(ServerConf.Redis.PoolKeepAliveIntervalSeconds) * time.Second
	ServerConf.Redis.PoolClientMaxIdleTime = time.Duration(ServerConf.Redis.PoolClientMaxIdleSeconds) * time.Second

	ServerConf.MySQL.RWTimeout = time.Duration(ServerConf.MySQL.RWTimeoutSeconds) * time.Second

	return nil
}

func parseServerConf(confDir string) error {
	err := parseConfFile(confDir + "/common/server_conf.toml")
	if err != nil {
		return fmt.Errorf("parse server_conf.toml error: %w", err)
	}

	err = parseConfFile(confDir + "/rewrite/server_conf_rewrite.toml")
	if err != nil {
		return fmt.Errorf("parse server_conf_rewrite.toml error: %w", err)
	}

	return nil
}

func parseConfFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(content, &ServerConf)
}
