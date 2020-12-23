package define

type SqlQueryParams struct {
	ParamsStructPtr interface{}
	Exists          map[string]bool
	Conditions      map[string]string

	OrderBy string
	Offset  int64
	Cnt     int64
}
