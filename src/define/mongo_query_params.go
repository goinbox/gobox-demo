package define

type MongoQueryParams struct {
	ParamsStructPtr interface{}
	Exists          map[string]bool
	Conditions      map[string]string

	OrderBy []string
	Offset  int
	Cnt     int
}
