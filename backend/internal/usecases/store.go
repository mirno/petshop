package usecases

type Store[T any] interface {
	Save(key string, value T) error
	Keys() ([]string, error)
	Get(key string) (T, error)
	Delete(key string) error
}
