package app

import (
	"context"
	"fmt"

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
	IsInList(ctx context.Context, address, listType string) (bool, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) IsInBlacklist(ctx context.Context, ip string) bool {
	isInList, err := a.storage.IsInList(ctx, ip, storage.Blacklist)
	if err != nil {
		a.logger.Error(fmt.Sprintf("storage error: %s", err))
		return false
	}

	return isInList
}

func (a *App) IsInWhitelist(ctx context.Context, ip string) bool {
	isInList, err := a.storage.IsInList(ctx, ip, storage.Whitelist)
	if err != nil {
		a.logger.Error(fmt.Sprintf("storage error: %s", err))
		return false
	}

	return isInList
}

func (a *App) HasLimits(ctx context.Context, login, password, ip string) bool {
	return false
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
