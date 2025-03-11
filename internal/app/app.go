package app

import (
	"context"
	"fmt"

	"github.com/evg555/antibrutforce/internal/ratelimiter"
	"github.com/evg555/antibrutforce/internal/storage"
)

type App struct {
	logger      Logger
	storage     Storage
	rateLimiter *ratelimiter.AuthRateLimiter
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	Save(ctx context.Context, subnet storage.Subnet, listType string) error
	Find(ctx context.Context, address, listType string) (*storage.Subnet, error)
	Delete(ctx context.Context, address, listType string) error
	IsInList(ctx context.Context, address, listType string) (bool, error)
}

func New(logger Logger, storage Storage, rateLimiter *ratelimiter.AuthRateLimiter) *App {
	return &App{
		logger:      logger,
		storage:     storage,
		rateLimiter: rateLimiter,
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

func (a *App) HasLimits(login, password, ip string) bool {
	return a.rateLimiter.AllowAttempt(login, password, ip)
}

func (a *App) AddIPWhitelist(ctx context.Context, subnet string) error {
	return a.storage.Save(ctx, storage.Subnet{Address: subnet}, storage.Whitelist)
}

func (a *App) DeleteIPWhitelist(ctx context.Context, subnet string) error {
	return a.storage.Delete(ctx, subnet, storage.Whitelist)
}

func (a *App) AddIPBlacklist(ctx context.Context, subnet string) error {
	return a.storage.Save(ctx, storage.Subnet{Address: subnet}, storage.Blacklist)
}

func (a *App) DeleteIPBlacklist(ctx context.Context, subnet string) error {
	return a.storage.Delete(ctx, subnet, storage.Blacklist)
}

func (a *App) ResetBucket(password, ip string) {
	a.rateLimiter.ResetBucket(password, ip)
}
