package rotatorapp

import (
	"context"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (a *App) CacheGetRotation(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) {

}
