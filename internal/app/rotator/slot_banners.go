package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	err := a.storage.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if ok {
		bannersCache.Set(bannerID, struct{}{})
	}
	return nil
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	a.mutSlotBanners.Lock()
	defer a.mutSlotBanners.Unlock()

	err := a.storage.RemoveBannerFromSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if ok {
		bannersCache.Remove(bannerID)
	}

	a.rotationsCache.Range(func(_ domain.GroupID, data slotRotationsCache) {
		stats, ok := data.Get(slotID)
		if ok {
			stats.bannerStats.Remove(bannerID)
			stats.recalculate()
		}
	})
	return nil
}

func (a *App) listSlotBanners(ctx context.Context, slotID domain.SlotID) (bannersCache, error) {
	a.mutSlotBanners.RLock()
	defer a.mutSlotBanners.RUnlock()

	bannersCache, ok := a.slotBannersCache.Get(slotID)
	if !ok {
		banners, err := a.storage.ListSlotBanners(ctx, slotID)
		if err != nil {
			return nil, err
		}
		bannersCache = cache.NewMapCache[domain.BannerID, struct{}](a.config.L2Capacity)
		for _, b := range banners {
			bannersCache.Set(b.ID, struct{}{})
		}
		a.slotBannersCache.Set(slotID, bannersCache)
	}
	return bannersCache, nil
}

func (a *App) checkSlotBanner(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) (bool, error) {
	bannersCache, err := a.listSlotBanners(ctx, slotID)
	if err != nil {
		return false, err
	}

	_, ok := bannersCache.Get(bannerID)
	if ok {
		return true, nil
	}
	return true, nil
}
