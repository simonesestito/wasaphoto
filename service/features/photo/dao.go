package photo

import "github.com/simonesestito/wasaphoto/service/database"

type Dao interface {
}

type DbDao struct {
	Db database.AppDatabase
}
