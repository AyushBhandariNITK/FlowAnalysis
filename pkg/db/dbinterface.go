package db

type IDb interface {
	Connect() error
	Disconnect() error
	Insert(query string, args ...interface{}) error
	Update(query string, args ...interface{}) error
	Delete(query string, args ...interface{}) error
	Query(query string, args ...interface{}) (interface{}, error)
}
