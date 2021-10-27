package mysql

type Dao interface {
	Insert(entities ...interface{}) error
}
