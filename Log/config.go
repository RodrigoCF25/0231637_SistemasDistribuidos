package log

type Segment struct {
	MaxStoreBytes uint64
	MaxIndexBytes uint64
	InitialOffset uint64
}

type Config struct {
	Segment
}

func NewConfig(maxStoreBytes, maxIndexBytes, initialOffset uint64) *Config {
	return &Config{
		Segment: Segment{
			MaxStoreBytes: maxStoreBytes,
			MaxIndexBytes: maxIndexBytes,
			InitialOffset: initialOffset,
		},
	}
}
