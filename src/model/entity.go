package model

type BaseEntity struct {
	ID       *int64 `column:"id"`
	AddTime  *string
	EditTime *string
}

func (e *BaseEntity) Nil() bool {
	if *e.ID > 0 {
		return false
	}

	return true
}
