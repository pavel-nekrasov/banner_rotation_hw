package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

func (a *App) CreateBanner(ctx context.Context, model models.Banner) (models.Banner, error) {
	if model.ID == "" {
		return models.Banner{}, customerrors.ValidationError{Field: "ID", Err: errCannotBeEmpty}
	}
	if model.Description == "" {
		return models.Banner{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.CreateBanner(ctx, model)
	if err != nil {
		return models.Banner{}, err
	}
	return model, nil
}

func (a *App) GetBanner(ctx context.Context, eventID string) (models.Banner, error) {
	return a.storage.GetBanner(ctx, eventID)
}

func (a *App) ListBanners(ctx context.Context) ([]models.Banner, error) {
	return a.storage.ListBanners(ctx)
}

func (a *App) UpdateBanner(ctx context.Context, model models.Banner) (models.Banner, error) {
	if model.Description == "" {
		return models.Banner{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.UpdateBanner(ctx, model)
	if err != nil {
		return models.Banner{}, err
	}
	return model, nil
}

func (a *App) DeleteBanner(ctx context.Context, eventID string) error {
	return a.storage.DeleteBanner(ctx, eventID)
}

func (a *App) CreateGroup(ctx context.Context, model models.Group) (models.Group, error) {
	if model.ID == "" {
		return models.Group{}, customerrors.ValidationError{Field: "ID", Err: errCannotBeEmpty}
	}
	if model.Description == "" {
		return models.Group{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.CreateGroup(ctx, model)
	if err != nil {
		return models.Group{}, err
	}
	return model, nil
}

func (a *App) GetGroup(ctx context.Context, eventID string) (models.Group, error) {
	return a.storage.GetGroup(ctx, eventID)
}

func (a *App) ListGroups(ctx context.Context) ([]models.Group, error) {
	return a.storage.ListGroups(ctx)
}

func (a *App) UpdateGroup(ctx context.Context, model models.Group) (models.Group, error) {
	if model.Description == "" {
		return models.Group{}, customerrors.ValidationError{Field: "Description", Err: errCannotBeEmpty}
	}
	err := a.storage.UpdateGroup(ctx, model)
	if err != nil {
		return models.Group{}, err
	}
	return model, nil
}

func (a *App) DeleteGroup(ctx context.Context, eventID string) error {
	return a.storage.DeleteGroup(ctx, eventID)
}

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

func (a *App) GetSlot(ctx context.Context, eventID string) (models.Slot, error) {
	return a.storage.GetSlot(ctx, eventID)
}

func (a *App) ListSlots(ctx context.Context) ([]models.Slot, error) {
	return a.storage.ListSlots(ctx)
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

func (a *App) DeleteSlot(ctx context.Context, eventID string) error {
	return a.storage.DeleteSlot(ctx, eventID)
}
