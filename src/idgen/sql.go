package idgen

import "github.com/goinbox/mysql"

const GenIdSql = "UPDATE id_gen SET max_id = last_insert_id(max_id + 1) WHERE name = ?"

type SqlIdGenter struct {
	client *mysql.Client
}

func NewSqlIdGenter(client *mysql.Client) *SqlIdGenter {
	return &SqlIdGenter{
		client: client,
	}
}

func (s *SqlIdGenter) SetClient(client *mysql.Client) *SqlIdGenter {
	s.client = client

	return s
}

func (s *SqlIdGenter) GenId(name string) (int64, error) {
	r, err := s.client.Exec(GenIdSql, name)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
