package app

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"

	"gdemo/resource"
)

const (
	updateAppsCacheInterval = time.Second * 60
)

type App struct {
	Name  string
	Token string
}

type appsCache struct {
	cacheTime time.Time
	appMap    map[string]*App
}

var appsCacheData *appsCache

func InitAppsCacheData(dir string) error {
	cache, err := loadAppsCache(dir)
	if err != nil {
		return fmt.Errorf("loadAppsCache error: %w", err)
	}
	appsCacheData = cache

	go func() {
		for {
			time.Sleep(updateAppsCacheInterval)

			cache, err = loadAppsCache(dir)
			if err != nil {
				resource.AccessLogger.Error("loadAppsCache error", golog.ErrorField(err))
			} else {
				resource.AccessLogger.Info("update appsCacheData")
				appsCacheData = cache
			}
		}
	}()

	return nil
}

func loadAppsCache(dir string) (*appsCache, error) {
	files, err := gomisc.ListFilesInDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ListFilesInDir error: %w", err)
	}

	cache := &appsCache{
		cacheTime: time.Now(),
		appMap:    make(map[string]*App),
	}
	for _, path := range files {
		matched, _ := filepath.Match("*.toml", filepath.Base(path))
		if !matched {
			continue
		}

		app, err := parseAppFile(path)
		if err != nil {
			return nil, fmt.Errorf("parseAppFile error: %w", err)
		}

		cache.appMap[app.Name] = app
	}

	return cache, nil
}

func parseAppFile(path string) (*App, error) {
	app := new(App)
	_, err := toml.DecodeFile(path, app)
	if err != nil {
		return nil, fmt.Errorf("toml.DecodeFile error: %w", err)
	}

	return app, nil
}

func LatestCacheTime() time.Time {
	return appsCacheData.cacheTime
}
