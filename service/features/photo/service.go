package photo

import "github.com/simonesestito/wasaphoto/service/storage"

type Service interface {
}

type ServiceImpl struct {
	Db      Dao
	Storage storage.Storage
}
