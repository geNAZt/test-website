package storage

type Storage interface {
	Store([]byte, string) (bool, error)
	Exists(string) bool
	Delete(string) (bool, error)
	GetUrl(string) (string, error)
}

func GetStorage() Storage {
	return new(fileSystemStorage)
}
