package app

type App struct {
	logger  Logger
	storage Storage
}

//go:generate mockery --name=Logger
type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

//go:generate mockery --name=Storage
type Storage interface{}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}
