package idgen

type Entity struct {
	ID    *int64 `column:"id"`
	Name  string
	MaxID int64 `column:"max_id"`
}
