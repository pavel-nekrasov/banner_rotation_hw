package rotatorapp //nolint:dupl

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) CreateSlot(ctx context.Context, model domain.Slot) (domain.Slot, error) {
	if model.ID == "" {
		return domain.Slot{}, customerrors.ValidationError{Field: "ID", Err: ErrCannotBeEmpty}
	}
	if model.Description == "" {
		return domain.Slot{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.CreateSlot(ctx, model)
	if err != nil {
		return domain.Slot{}, err
	}
	return model, nil
}

func (a *App) GetSlot(ctx context.Context, entityID domain.SlotID) (domain.Slot, error) {
	return a.storage.GetSlot(ctx, entityID)
}

func (a *App) ListSlots(ctx context.Context) ([]domain.Slot, error) {
	return a.storage.ListSlots(ctx)
}

func (a *App) UpdateSlot(ctx context.Context, model domain.Slot) (domain.Slot, error) {
	if model.Description == "" {
		return domain.Slot{}, customerrors.ValidationError{Field: "Description", Err: ErrCannotBeEmpty}
	}
	err := a.storage.UpdateSlot(ctx, model)
	if err != nil {
		return domain.Slot{}, err
	}
	return model, nil
}

func (a *App) DeleteSlot(ctx context.Context, entityID domain.SlotID) error {
	err := a.storage.DeleteSlot(ctx, entityID)
	if err != nil {
		return err
	}
	a.slotsCache.Remove(entityID)
	return nil
}

func (a *App) checkSlot(ctx context.Context, entityID domain.SlotID) error {
	_, ok := a.slotsCache.Get(entityID)

	if ok {
		return nil
	}

	entity, err := a.storage.GetSlot(ctx, entityID)
	if err != nil {
		return err
	}

	a.slotsCache.Set(entityID, entity)
	return nil
}
