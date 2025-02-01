package rotatorapp //nolint:dupl

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) CreateBanner(ctx context.Context, model domain.Banner) (domain.Banner, error) {
	if model.ID == "" {
		return domain.Banner{}, customerrors.ValidationError{Field: "ID", Err: ErrCannotBeEmpty}
	}
	if model.Description == "" {
		return domain.Banner{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.CreateBanner(ctx, model)
	if err != nil {
		return domain.Banner{}, err
	}
	return model, nil
}

func (a *App) GetBanner(ctx context.Context, entityID domain.BannerID) (domain.Banner, error) {
	return a.storage.GetBanner(ctx, entityID)
}

func (a *App) ListBanners(ctx context.Context) ([]domain.Banner, error) {
	return a.storage.ListBanners(ctx)
}

func (a *App) UpdateBanner(ctx context.Context, model domain.Banner) (domain.Banner, error) {
	if model.Description == "" {
		return domain.Banner{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.UpdateBanner(ctx, model)
	if err != nil {
		return domain.Banner{}, err
	}
	return model, nil
}

func (a *App) DeleteBanner(ctx context.Context, entityID domain.BannerID) error {
	err := a.storage.DeleteBanner(ctx, entityID)
	if err != nil {
		return err
	}
	a.bannersCache.Remove(entityID)
	return nil
}

func (a *App) checkBanner(ctx context.Context, entityID domain.BannerID) error {
	_, ok := a.bannersCache.Get(entityID)

	if ok {
		return nil
	}

	entity, err := a.storage.GetBanner(ctx, entityID)
	if err != nil {
		return err
	}

	a.bannersCache.Set(entityID, entity)
	return nil
}
