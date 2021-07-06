package storage

// Storage service storage interface
type Storage interface {
	WriteLineRate(k float64, line string) error
	GetLineRate(line string) (float64, error)
	CheckConnection() error
}