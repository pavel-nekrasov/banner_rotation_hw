package rotatorapp

import (
	"context"

	appdomain "github.com/pavel-nekrasov/banner_rotation_hw/internal/app/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) ClickBanner(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	var err error

	err = a.checkGroup(ctx, groupID)
	if err != nil {
		return err
	}

	err = a.checkSlot(ctx, slotID)
	if err != nil {
		return err
	}

	err = a.checkBanner(ctx, bannerID)
	if err != nil {
		return err
	}

	ok, err := a.checkSlotBanner(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	if !ok {
		return ErrBannerNoAllowedForSlot
	}

	// a.mut.Lock()
	// defer a.mut.Unlock()

	statData, err := a.getStats(ctx, groupID, slotID)
	if err != nil {
		return err
	}

	err = a.storage.IncrementBannerStatClick(ctx, groupID, slotID, bannerID)
	if err != nil {
		return err
	}

	statData.IncrementClick(bannerID)
	// TODO: add queue publish

	return nil
}

func (a *App) SelectBanner(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (domain.BannerID, error) {
	var emptyBannerID domain.BannerID
	var err error

	err = a.checkGroup(ctx, groupID)
	if err != nil {
		return emptyBannerID, err
	}

	err = a.checkSlot(ctx, slotID)
	if err != nil {
		return emptyBannerID, err
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	statData, err := a.getStats(ctx, groupID, slotID)
	if err != nil {
		return emptyBannerID, err
	}
	bestBannerID := statData.BestBanner()

	err = a.storage.IncrementBannerStatShow(ctx, groupID, slotID, bestBannerID)
	if err != nil {
		return emptyBannerID, err
	}

	statData.IncrementShow(bestBannerID)
	// TODO: add queue publish

	return bestBannerID, nil
}

func (a *App) getStats(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (*appdomain.StatData, error) {
	var statData *appdomain.StatData

	l2Cache, ok := a.rotationsCache.Get(groupID)
	if !ok {
		l2Cache = cache.NewLRUCache[domain.SlotID, *appdomain.StatData](a.config.L2Capacity)
		a.rotationsCache.Set(groupID, l2Cache)
	}

	statData, ok = l2Cache.Get(slotID)
	if ok {
		return statData, nil
	}

	slotBanners, err := a.listSlotBanners(ctx, slotID)
	if err != nil {
		return nil, err
	}

	bannerStats, err := a.storage.ListBannerStats(ctx, groupID, slotID)
	if err != nil {
		return nil, err
	}

	statData = appdomain.NewStatData()

	slotBanners.Range(func(key domain.BannerID, _ struct{}) {
		statData.Set(key, &domain.BannerStat{
			BannerID:   key,
			ShowCount:  1,
			ClickCount: 1,
		})
	})

	for _, row := range bannerStats {
		statData.Set(row.BannerID, &row)
	}
	if statData.Empty() {
		return nil, ErrNoBannersAssignedForSlot
	}
	statData.Recalculate()
	l2Cache.Set(slotID, statData)
	return statData, nil
}
