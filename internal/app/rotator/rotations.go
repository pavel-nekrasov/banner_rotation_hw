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
	ok, err := a.CheckSlotBanner(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}
	a.mutRotations.Lock()
	defer a.mutRotations.Unlock()

	statData, err := a.loadStatData(ctx, groupID, slotID)
	if err != nil {
		return err
	}
	err = a.storage.IncrementRotationClick(ctx, groupID, slotID, bannerID)
	if err != nil {
		return err
	}
	rotation, err := a.storage.GetRotation(ctx, groupID, slotID, bannerID)
	if err != nil {
		return err
	}
	statData.rotations[rotation.BannerID] = rotation

	// TODO: add queue publish

	return nil
}

func (a *App) SelectBanner(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (domain.Banner, bool, error) {
	a.mutRotations.Lock()
	defer a.mutRotations.Unlock()

	statData, err := a.loadStatData(ctx, groupID, slotID)
	if err != nil {
		return domain.Banner{}, false, err
	}

	// implement calc logic

	// TODO: add queue publish
	return domain.Banner{}, false, nil
}

func (a *App) loadStatData(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) (StatData, error) {
	l2Cache, ok := a.rotationsCache.Get(groupID)
	if !ok {
		l2Cache = cache.NewCache[domain.SlotID, StatData](a.config.L2Capacity)
		a.rotationsCache.Set(groupID, l2Cache)
	}

	statData, ok := l2Cache.Get(slotID)

	if !ok {
		rotations, err := a.storage.ListRotations(ctx, groupID, slotID)
		if err != nil {
			return statData, err
		}
		statData.rotations = map[domain.BannerID]domain.Rotation{}
		for _, r := range rotations {
			statData.rotations[r.BannerID] = r
		}
		l2Cache.Set(slotID, statData)
	}
	return statData, nil
}
