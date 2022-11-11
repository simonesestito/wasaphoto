package database

// AppDatabase is an abstraction over platform *sql.DB
type AppDatabase interface {
	Ping() error
	Version() (int, error)
}
