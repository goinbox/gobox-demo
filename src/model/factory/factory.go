package factory

type DaoFactory interface {
}

var DefaultDaoFactory DaoFactory = new(daoFactory)

type daoFactory struct {
}
