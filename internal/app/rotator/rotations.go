package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) RegisterClick(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	ok, err := a.checkSlotBanner(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	if !ok {
		return errBannerNoAllowedForSlot
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	data, err := a.loadStats(ctx, groupID, slotID)
	if err != nil {
		return err
	}

	err = a.storage.IncrementRotationClick(ctx, groupID, slotID, bannerID)
	if err != nil {
		return err
	}

	bannerStat, ok := data.bannerStats.Get(bannerID)

	if ok {
		bannerStat.ClickCount++
		data.recalculate()
		// TODO: add queue publish
	}

	return nil
}

func (a *App) ShowBanner(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (domain.BannerID, bool, error) {
	a.mut.Lock()
	defer a.mut.Unlock()

	var bannerID domain.BannerID

	data, err := a.loadStats(ctx, groupID, slotID)
	if err != nil {
		return bannerID, false, err
	}

	if data.empty() {
		return bannerID, false, errNoBannersAssignedForSlot
	}

	err = a.storage.IncrementRotationShow(ctx, groupID, slotID, data.bestBannerID)
	if err != nil {
		return bannerID, false, err
	}
	bannerID = data.bestBannerID
	bannerStat, ok := data.bannerStats.Get(bannerID)

	if ok {
		bannerStat.ShowCount++
		data.recalculate()
		// TODO: add queue publish
	}

	return bannerID, false, nil
}

func (a *App) loadStats(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (statData, error) {
	l2Cache, ok := a.rotationsCache.Get(groupID)
	if !ok {
		l2Cache = cache.NewLRUCache[domain.SlotID, statData](a.config.L2Capacity)
		a.rotationsCache.Set(groupID, l2Cache)
	}

	data, ok := l2Cache.Get(slotID)

	if !ok {
		bannerStats, err := a.storage.ListBannerStats(ctx, groupID, slotID)
		if err != nil {
			return statData{}, err
		}

		data.bannerStats = cache.NewMapCache[domain.BannerID, domain.BannerStat](a.config.L2Capacity)

		slotBanners, err := a.listSlotBanners(ctx, slotID)
		if err != nil {
			return statData{}, err
		}

		slotBanners.Range(func(key domain.BannerID, _ struct{}) {
			data.bannerStats.Set(key, domain.BannerStat{
				BannerID:   key,
				ShowCount:  1,
				ClickCount: 1,
			})
		})

		for _, row := range bannerStats {
			data.bannerStats.Set(row.BannerID, row)
		}
		if !data.empty() {
			data.recalculate()
		}

		l2Cache.Set(slotID, data)
	}
	return data, nil
}
