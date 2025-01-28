package rotatorapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	statData, err := a.getStats(ctx, key)
	if err != nil {
		return err
	}

	// лочимся на конкретном бакете чтобы не мешать обработке событий с другим Group/Slot
	statData.Lock()
	defer statData.Unlock()

	if statData.Empty() {
		return ErrNoBannersAssignedForSlot
	}

	statData.Recalculate()

	err = a.storage.IncrementBannerStatClick(ctx, groupID, slotID, bannerID)
	if err != nil {
		return err
	}

	statData.IncrementClick(bannerID)

	a.notify(appdomain.NotificationEvent{
		EventType: "click",
		GroupID:   key.GroupID,
		SlotID:    key.SlotID,
		BannerID:  bannerID,
		Timestamp: time.Now().UnixMilli(),
	})

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
	statData, err := a.getStats(ctx, key)
	if err != nil {
		return emptyBannerID, err
	}

	// лочимся на конкретном бакете чтобы не мешать обработке событий с другим Group/Slot
	statData.Lock()
	defer statData.Unlock()

	if statData.Empty() {
		return emptyBannerID, ErrNoBannersAssignedForSlot
	}

	statData.Recalculate()
	bestBannerID := statData.BestBanner()

	err = a.storage.IncrementBannerStatShow(ctx, groupID, slotID, bestBannerID)
	if err != nil {
		return emptyBannerID, err
	}

	statData.IncrementShow(bestBannerID)

	a.notify(appdomain.NotificationEvent{
		EventType: "show",
		GroupID:   key.GroupID,
		SlotID:    key.SlotID,
		BannerID:  bestBannerID,
		Timestamp: time.Now().UnixMilli(),
	})

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

	a.mut.Lock()
	defer a.mut.Unlock()

	// double check
	statData, ok = a.rotationsCache.Get(key)
	if ok {
		return statData, nil
	}

	allowedBanners, err := a.listSlotBanners(ctx, key.SlotID)
	if err != nil {
		return nil, err
	}

	bannerStats, err := a.storage.ListBannerStats(ctx, key.GroupID, key.SlotID)
	if err != nil {
		return nil, err
	}

	statData = appdomain.NewStatData()
	// заполняем из списка разрешенных баннеров дефолтными значениями
	allowedBanners.Range(func(key domain.BannerID, _ struct{}) {
		statData.Set(key, &domain.BannerStat{
			BannerID:   key,
			ShowCount:  1,
			ClickCount: 1,
		})
	})
	// заполняем данными из БД
	for _, row := range bannerStats {
		statData.Set(row.BannerID, &row)
	}
	a.rotationsCache.Set(key, statData)
	return statData, nil
}

func (a *App) notify(payload appdomain.NotificationEvent) error {
	if a.publisher == nil {
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}
	err = a.publisher.Publish(data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}
