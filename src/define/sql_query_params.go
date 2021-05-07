package define

import "github.com/goinbox/mysql"

type SqlQueryParams struct {
	CondItems []*mysql.SqlColQueryItem

	OrderBy string
	Offset  int64
	Cnt     int64
}
