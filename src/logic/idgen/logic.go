package idgen

import (
	"fmt"

	"github.com/goinbox/gomisc"

	"gdemo/model/idgen"
	"gdemo/pcontext"
	"gdemo/resource"
)

type Logic interface {
	GenerateID(ctx *pcontext.Context) int64
}

func NewLogic() Logic {
	return new(logic)
}

var snowflake *gomisc.Snowflake

func Init(idgenName string) error {
	dao := idgen.NewDao(resource.MySQLClient(resource.AccessLogger))
	machineID, err := dao.GenID(idgenName)
	if err != nil {
		return fmt.Errorf("dao.GenID error: %w", err)
	}

	snowflake = gomisc.NewSnowflake(machineID % 1024) // machineIDMax = 1023

	return nil
}

type logic struct {
}

func (l *logic) GenerateID(ctx *pcontext.Context) int64 {
	return snowflake.GenerateID()
}
