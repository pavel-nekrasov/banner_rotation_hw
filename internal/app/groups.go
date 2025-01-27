package rotatorapp //nolint:dupl

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) CreateGroup(ctx context.Context, model domain.Group) (domain.Group, error) {
	if model.ID == "" {
		return domain.Group{}, customerrors.ValidationError{Field: "ID", Err: ErrCannotBeEmpty}
	}
	if model.Description == "" {
		return domain.Group{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.CreateGroup(ctx, model)
	if err != nil {
		return domain.Group{}, err
	}
	return model, nil
}

func (a *App) GetGroup(ctx context.Context, entityID domain.GroupID) (domain.Group, error) {
	return a.storage.GetGroup(ctx, entityID)
}

func (a *App) ListGroups(ctx context.Context) ([]domain.Group, error) {
	return a.storage.ListGroups(ctx)
}

func (a *App) UpdateGroup(ctx context.Context, model domain.Group) (domain.Group, error) {
	if model.Description == "" {
		return domain.Group{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.UpdateGroup(ctx, model)
	if err != nil {
		return domain.Group{}, err
	}
	return model, nil
}

func (a *App) DeleteGroup(ctx context.Context, entityID domain.GroupID) error {
	err := a.storage.DeleteGroup(ctx, entityID)
	if err != nil {
		return err
	}
	a.groupsCache.Remove(entityID)
	return nil
}

func (a *App) checkGroup(ctx context.Context, entityID domain.GroupID) error {
	_, ok := a.groupsCache.Get(entityID)

	if ok {
		return nil
	}

	entity, err := a.storage.GetGroup(ctx, entityID)
	if err != nil {
		return err
	}

	a.groupsCache.Set(entityID, entity)
	return nil
}
