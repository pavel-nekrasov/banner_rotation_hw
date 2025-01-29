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

	err = a.storage.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	// очищаем кеш разрешенных баннеров для данного слота
	a.allowedBannersCache.Remove(slotID)
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

	err = a.storage.RemoveBannerFromSlot(ctx, slotID, bannerID)
	if err != nil {
		return err
	}
	// очищаем кеш разрешенных баннеров для данного слота
	a.allowedBannersCache.Remove(slotID)
	// очищаем кэш ротаций чтобы перегрузилась статистика ротаций
	a.rotationsCache.Clear()

	return nil
}

func (a *App) getAllowedBanners(ctx context.Context, slotID domain.SlotID) (bannersCache, error) {
	// лочимся на конкретном слоте чтобы избежать паралельной загрузки баннеров для одного и того же слота
	a.slotKeys.Lock(slotID)
	defer a.slotKeys.Unlock(slotID)

	result, ok := a.allowedBannersCache.Get(slotID)
	if !ok {
		banners, err := a.storage.ListSlotBanners(ctx, slotID)
		if err != nil {
			return nil, err
		}
		result = cache.NewMapCache[domain.BannerID, struct{}]()
		for _, b := range banners {
			result.Set(b.ID, struct{}{})
		}
		a.allowedBannersCache.Set(slotID, result)
	}
	return result, nil
}

func (a *App) checkAllowedSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (bool, error) {
	bannersCache, err := a.getAllowedBanners(ctx, slotID)
	if err != nil {
		return false, err
	}

	_, ok := bannersCache.Get(bannerID)
	return ok, nil
}
