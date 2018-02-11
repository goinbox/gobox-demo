package idgen

import "github.com/goinbox/mysql"

type SqlIdGenter struct {
	client *mysql.Client

	updateSql string
	selectSql string
}

func NewSqlIdGenter(client *mysql.Client) *SqlIdGenter {
	return &SqlIdGenter{
		client: client,

		updateSql: "UPDATE id_gen SET max_id = last_insert_id(max_id + 1) WHERE name = ?",
		selectSql: "SELECT last_insert_id()",
	}
}

func (s *SqlIdGenter) GenId(name string) (int64, error) {
	_, err := s.client.Exec(s.updateSql, name)
	if err != nil {
		return 0, err
	}

	var id int64
	err = s.client.QueryRow(s.selectSql).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
