package database

// AppDatabase is an abstraction over platform *sql.DB
type AppDatabase interface {
	Ping() error
	Version() (int, error)
	QueryStructRow(destPointer any, query string, args ...any) error
}
