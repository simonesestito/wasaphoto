package storage

// Storage is an abstraction and can be any storage
// (filesystem, S3 bucket, CDN, in-memory, ...)
type Storage interface {
	SaveFile(name string, data []byte) error
	DeleteFile(path string) error
}
