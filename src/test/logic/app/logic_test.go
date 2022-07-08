package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gdemo/logic/app"
	"gdemo/test"
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)
}

func TestInitAppCacheData(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Log(app.LatestCacheTime())
		time.Sleep(time.Second * 4)
	}
}

func TestListAllApps(t *testing.T) {
	l := app.NewLogic()

	for i, item := range l.ListAllApps(test.Context()) {
		t.Log(i, item)
	}
}

func TestAppByName(t *testing.T) {
	l := app.NewLogic()
	ctx := test.Context()

	for i, item := range l.ListAllApps(test.Context()) {
		item = l.AppByName(ctx, item.Name)
		t.Log(i, item)
	}
}

func TestGenerateToken(t *testing.T) {
	l := app.NewLogic()
	ctx := test.Context()

	t.Log(l.GenerateToken(ctx, "demo"))
}
