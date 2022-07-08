package misc

import (
	"testing"
)

func init() {
	MysqlUpdateColumnValueConvertFuncMap["convert_is_student"] = convertIsStudent
}

type modifyParams struct {
	ID        int `column:"id"`
	Name      *string
	Age       *int
	IsStudent *bool  `mysql_update_column:"convert_is_student"`
	Version   string `mysql_update_column:"omit"`
}

func convertIsStudent(v interface{}) interface{} {
	if v.(bool) == true {
		return 1
	}
	return 0
}

func TestMakeMysqlUpdateColumns(t *testing.T) {
	id := 1
	name := "tom"
	age := 10
	isStudent := true

	params := &modifyParams{
		ID:        id,
		Name:      &name,
		Age:       &age,
		IsStudent: &isStudent,
	}

	for i, item := range MakeMysqlUpdateColumns(params) {
		t.Log(i, item)
	}
}
