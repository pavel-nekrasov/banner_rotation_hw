package rotatorapp

import (
	"context"
	"errors"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	if slotID == "" {
		return customerrors.ParamError{Param: "slotID", Err: errCannotBeEmpty}
	}
	if bannerID == "" {
		return customerrors.ParamError{Param: "bannerID", Err: errCannotBeEmpty}
	}
	err := a.storage.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	if slotID == "" {
		return customerrors.ParamError{Param: "slotID", Err: errCannotBeEmpty}
	}
	if bannerID == "" {
		return customerrors.ParamError{Param: "bannerID", Err: errCannotBeEmpty}
	}
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	err := a.storage.RemoveBannerFromSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	a.CacheRemoveSlotBanner(ctx, slotID, bannerID)
	a.rotationsCache.Clear()
	return nil
}

func (a *App) CheckSlotBanner(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) (bool, error) {
	if slotID == "" {
		return false, customerrors.ParamError{Param: "slotID", Err: errCannotBeEmpty}
	}
	if bannerID == "" {
		return false, customerrors.ParamError{Param: "bannerID", Err: errCannotBeEmpty}
	}
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	_, ok := a.CacheGetSlotBanner(ctx, slotID, bannerID)
	if ok {
		return true, nil
	}
	banner, err := a.storage.GetSlotBanner(ctx, slotID, bannerID)
	var notFoundErr customerrors.NotFound
	if errors.As(err, &notFoundErr) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	a.CacheAddSlotBanner(ctx, slotID, banner)
	return true, nil
}
