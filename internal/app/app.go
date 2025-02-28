package app

import (
	"context"

	"github.com/evg555/antibrutforce/internal/storage"
)

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
type Storage interface {
	Save(ctx context.Context, subnet storage.Subnet, listType string) error
	Find(ctx context.Context, address, listType string) (*storage.Subnet, error)
	Delete(ctx context.Context, address, listType string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) AddIpWhitelist(ctx context.Context, subnet string) error {
	return a.storage.Save(ctx, storage.Subnet{Address: subnet}, storage.Whitelist)
}

func (a *App) DeleteIpWhitelist(ctx context.Context, subnet string) error {
	return a.storage.Delete(ctx, subnet, storage.Whitelist)
}

func (a *App) AddIpBlacklist(ctx context.Context, subnet string) error {
	return a.storage.Save(ctx, storage.Subnet{Address: subnet}, storage.Blacklist)
}

func (a *App) DeleteIpBlacklist(ctx context.Context, subnet string) error {
	return a.storage.Delete(ctx, subnet, storage.Blacklist)
}
