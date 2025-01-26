package rotatorapp

import (
	"context"

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
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	err := a.storage.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if ok {
		bannersCache[bannerID] = struct{}{}
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
	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if ok {
		delete(bannersCache, bannerID)
	}

	a.rotationsCache.Range(func(_ domain.GroupID, data slotRotationsCache) {
		stats, ok := data.Get(slotID)
		if ok {
			delete(stats.bannerStats, bannerID)
			stats.recalculate()
		}
	})
	return nil
}

func (a *App) listSlotBanners(ctx context.Context, slotID domain.SlotID) (map[domain.BannerID]struct{}, error) {
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if !ok {
		banners, err := a.storage.ListSlotBanners(ctx, slotID)
		if err != nil {
			return nil, err
		}
		bannersCache = make(map[domain.BannerID]struct{})
		for _, b := range banners {
			bannersCache[b.ID] = struct{}{}
		}
		a.slotBannersCache.Set(slotID, bannersCache)
	}
	return bannersCache, nil
}

func (a *App) checkSlotBanner(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) (bool, error) {
	bannersCache, err := a.listSlotBanners(ctx, slotID)
	if err != nil {
		return false, nil
	}
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	_, ok := bannersCache[bannerID]
	if ok {
		return true, nil
	}
	return true, nil
}
