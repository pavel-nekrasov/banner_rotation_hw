package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	var err error

	err = a.checkSlot(ctx, slotID)
	if err != nil {
		return err
	}

	err = a.checkBanner(ctx, bannerID)
	if err != nil {
		return err
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	err = a.storage.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	allowedBannersCache, ok := a.allowedSlotBannersCache.Get(slotID)
	if ok {
		allowedBannersCache.Set(bannerID, struct{}{})
	}
	// очищаем кэш ротаций чтобы перегрузилась статистика ротаций
	a.rotationsCache.Clear()

	return nil
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	var err error

	err = a.checkSlot(ctx, slotID)
	if err != nil {
		return err
	}

	err = a.checkBanner(ctx, bannerID)
	if err != nil {
		return err
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	err = a.storage.RemoveBannerFromSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	allowedBannersCache, ok := a.allowedSlotBannersCache.Get(slotID)
	if ok {
		allowedBannersCache.Remove(bannerID)
	}
	// очищаем кэш ротаций чтобы перегрузилась статистика ротаций
	a.rotationsCache.Clear()

	return nil
}

func (a *App) listAllowedSlotBanners(ctx context.Context, slotID domain.SlotID) (bannersCache, error) {
	bannersCache, ok := a.allowedSlotBannersCache.Get(slotID)
	if !ok {
		banners, err := a.storage.ListSlotBanners(ctx, slotID)
		if err != nil {
			return nil, err
		}
		bannersCache = cache.NewMapCache[domain.BannerID, struct{}]()
		for _, b := range banners {
			bannersCache.Set(b.ID, struct{}{})
		}
		a.allowedSlotBannersCache.Set(slotID, bannersCache)
	}
	return bannersCache, nil
}

func (a *App) checkAllowedSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (bool, error) {
	bannersCache, err := a.listAllowedSlotBanners(ctx, slotID)
	if err != nil {
		return false, err
	}

	_, ok := bannersCache.Get(bannerID)
	return ok, nil
}
