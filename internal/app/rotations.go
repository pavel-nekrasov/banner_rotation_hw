package rotatorapp

import (
	"context"

	appdomain "github.com/pavel-nekrasov/banner_rotation_hw/internal/app/domain"
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

	key := appdomain.GroupSlotKey{GroupID: groupID, SlotID: slotID}
	a.rotationsCache.Lock(key)
	defer a.rotationsCache.Unlock(key)

	statData, err := a.getStats(ctx, key)
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

	key := appdomain.GroupSlotKey{GroupID: groupID, SlotID: slotID}
	a.rotationsCache.Lock(key)
	defer a.rotationsCache.Unlock(key)

	statData, err := a.getStats(ctx, key)
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
	key appdomain.GroupSlotKey,
) (*appdomain.StatData, error) {
	var statData *appdomain.StatData

	statData, ok := a.rotationsCache.Get(key)

	if ok {
		return statData, nil
	}

	slotBanners, err := a.listSlotBanners(ctx, key.SlotID)
	if err != nil {
		return nil, err
	}

	bannerStats, err := a.storage.ListBannerStats(ctx, key.GroupID, key.SlotID)
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
	a.rotationsCache.Set(key, statData)
	return statData, nil
}
