package conf

type ConfigResource[T any] interface {
	Load() error
	Watch()
	Get() T
}

type ConfigOption func() string
