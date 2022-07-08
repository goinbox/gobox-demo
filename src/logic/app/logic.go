package app

import (
	"fmt"
	"time"

	"gdemo/misc"
	"gdemo/pcontext"
)

type Logic interface {
	ListAllApps(ctx *pcontext.Context) []*App
	AppByName(ctx *pcontext.Context, appName string) *App
	GenerateToken(ctx *pcontext.Context, appName string) string
}

type logic struct {
}

func NewLogic() *logic {
	return &logic{}
}

func (l *logic) ListAllApps(ctx *pcontext.Context) []*App {
	var result []*App
	for _, app := range appsCacheData.appMap {
		result = append(result, app)
	}

	return result
}

func (l *logic) AppByName(ctx *pcontext.Context, appName string) *App {
	return appsCacheData.appMap[appName]
}

func (l *logic) GenerateToken(ctx *pcontext.Context, appName string) string {
	token := fmt.Sprintf("gdemo-app-%s-t-%d", appName, time.Now().Second())

	return misc.Sha256(token)
}
