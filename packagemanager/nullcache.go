package packagemanager

type (
	NullCache struct {
	}
)

func (*NullCache) Put(key string, value []byte) error {
	return nil
}

func (*NullCache) Get(key string) ([]byte, error) {
	return nil, nil
}
