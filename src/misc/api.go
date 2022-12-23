package misc

import (
	"fmt"
)

func ApiSign(token string, ts int64, payload []byte) string {
	return Sha256(fmt.Sprintf("%s%d%s", token, ts, payload))
}

type CommonListParams struct {
	Offset int64 `validate:"omitempty,min=0"`
	Limit  int64 `validate:"omitempty,min=0"`

	OrderBy string
	Order   string
}

func NewDefaultCommonListParams() *CommonListParams {
	return &CommonListParams{
		Offset:  0,
		Limit:   20,
		OrderBy: "id",
		Order:   "desc",
	}
}

func (p *CommonListParams) MysqlOrderBy() string {
	return fmt.Sprintf("%s %s", p.OrderBy, p.Order)
}
