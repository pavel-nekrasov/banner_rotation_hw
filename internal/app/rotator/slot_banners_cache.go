package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) CacheGetSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (domain.Banner, bool) {
	banners, ok := a.slotBannersCache.Get(slotID)
	if ok {
		return banners.Get(bannerID)
	}
	return domain.Banner{}, false
}

func (a *App) CacheAddSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	banner domain.Banner,
) {
	banners, ok := a.slotBannersCache.Get(slotID)
	if !ok {
		banners = cache.NewCache[domain.BannerID, domain.Banner](a.config.L2Capacity)
		a.slotBannersCache.Set(slotID, banners)
	}
	banners.Set(banner.ID, banner)
}

func (a *App) CacheRemoveSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) {
	banners, ok := a.slotBannersCache.Get(slotID)
	if ok {
		banners.Remove(bannerID)
	}
}
