package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"time"

	"github.com/BurntSushi/toml"

	"gdemo/perror"
)

func Init(confDir string) *perror.Error {
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

	ServerConf.Mysql.RWTimeout = time.Duration(ServerConf.Mysql.RWTimeoutSeconds) * time.Second

	return nil
}

func parseServerConf(confDir string) *perror.Error {
	err := parseConfFile(confDir + "/server_conf.toml")
	if err != nil {
		return perror.New(perror.ESysParseServerConfFail,
			fmt.Sprintf("parse server_conf.toml error: %v", err))
	}

	err = parseConfFile(confDir + "/server_conf_rewrite.toml")
	if err != nil {
		return perror.New(perror.ESysParseServerConfFail,
			fmt.Sprintf("parse server_conf_rewrite.toml error: %v", err))
	}

	return nil
}

func parseConfFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(content, &ServerConf)
}
