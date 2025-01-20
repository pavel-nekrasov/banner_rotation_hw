package rotatorapp

import (
	"context"
	"errors"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/common"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

type App struct {
	logger  common.Logger
	storage Storage
}

type Storage interface {
	CreateBanner(ctx context.Context, entity models.Banner) error
	GetBanner(ctx context.Context, entityID string) (models.Banner, error)
	ListBanners(ctx context.Context) ([]models.Banner, error)
	UpdateBanner(ctx context.Context, entity models.Banner) error
	DeleteBanner(ctx context.Context, entityID string) error

	CreateGroup(ctx context.Context, entity models.Group) error
	GetGroup(ctx context.Context, entityID string) (models.Group, error)
	ListGroups(ctx context.Context) ([]models.Group, error)
	UpdateGroup(ctx context.Context, entity models.Group) error
	DeleteGroup(ctx context.Context, entityID string) error

	CreateSlot(ctx context.Context, entity models.Slot) error
	GetSlot(ctx context.Context, entityID string) (models.Slot, error)
	ListSlots(ctx context.Context) ([]models.Slot, error)
	UpdateSlot(ctx context.Context, entity models.Slot) error
	DeleteSlot(ctx context.Context, entityID string) error
}

func New(logger common.Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

var (
	errCannotBeEmpty = errors.New("cannot be empty")
)

/*
func (a *App) CreateSlot(ctx context.Context, model models.Slot) (models.Slot, error) {
	if model.ID == "" {
		return models.Slot{}, customerrors.ValidationError{Field: "ID", Err: errCannotBeEmpty}
	}
	if model.Description == "" {
		return models.Slot{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.CreateSlot(ctx, model)
	if err != nil {
		return models.Slot{}, err
	}
	return model, nil
}

func (a *App) UpdateSlot(ctx context.Context, model models.Slot) (models.Slot, error) {
	if model.Description == "" {
		return models.Slot{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.UpdateSlot(ctx, model)
	if err != nil {
		return models.Slot{}, err
	}
	return model, nil
}

func (a *App) GetEvent(ctx context.Context, eventID string) (models.Slot, error) {
	return a.storage.GetSlot(ctx, eventID)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.storage.DeleteSlot(ctx, eventID)
}
*/
/*
func (a *App) ListEventsForDate(ctx context.Context, ownerEmail string, date int64) ([]model.Event, error) {
	dt := time.Unix(date, 0)
	return a.storage.ListOwnerEventsForPeriod(ctx, ownerEmail, dt, dt.AddDate(0, 0, 1))
}

func (a *App) ListEventsForWeek(ctx context.Context, ownerEmail string, date int64) ([]model.Event, error) {
	dt := time.Unix(date, 0)
	return a.storage.ListOwnerEventsForPeriod(ctx, ownerEmail, dt, dt.AddDate(0, 0, 7))
}

func (a *App) validateAttributes(ctx context.Context, dto contracts.Event) (model.Event, error) {
	if dto.Title == "" {
		return model.Event{}, customerrors.ValidationError{Field: "Title", Err: errCannotBeEmpty}
	}

	startTime := time.Unix(dto.StartTime, 0)
	endTime := time.Unix(dto.EndTime, 0)

	if startTime.After(endTime) {
		return model.Event{}, customerrors.ValidationError{Field: "StartTime", Err: errWrongPeriod}
	}

	if startTime.Day() != endTime.Day() || startTime.Month() != endTime.Month() || startTime.Year() != endTime.Year() {
		return model.Event{}, customerrors.ValidationError{Field: "StartTime", Err: errNotSameDay}
	}

	if dto.OwnerEmail == "" {
		return model.Event{}, customerrors.ValidationError{Field: "OwnerEmail", Err: errCannotBeEmpty}
	}

	var notifyTime time.Time
	if dto.NotifyBefore != "" {
		duration, err := time.ParseDuration(dto.NotifyBefore)
		if err != nil {
			return model.Event{}, customerrors.ValidationError{Field: "Notify", Err: errWrongDateFormat}
		}

		notifyTime = startTime.Add(-duration)
	}

	existingEvents, err := a.storage.ListOwnerEventsForPeriod(ctx, dto.OwnerEmail, startTime, endTime)
	if err != nil {
		return model.Event{}, err
	}
	for _, ev := range existingEvents {
		if ev.ID != dto.ID && (startTime.Equal(ev.StartTime) || startTime.After(ev.StartTime)) &&
			(endTime.Equal(ev.EndTime) || endTime.Before(ev.EndTime)) {
			return model.Event{}, customerrors.ValidationError{Field: "StartTime|EndTime", Err: errPeriodIsBusy}
		}
	}

	return model.Event{
		ID:           dto.ID,
		Title:        dto.Title,
		StartTime:    startTime,
		EndTime:      endTime,
		Description:  dto.Description,
		OwnerEmail:   dto.OwnerEmail,
		NotifyBefore: dto.NotifyBefore,
		NotifyTime:   notifyTime,
	}, nil
}
*/
