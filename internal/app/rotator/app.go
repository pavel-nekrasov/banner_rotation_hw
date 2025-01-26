package rotatorapp

import (
	"context"
	"errors"
	"sync"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/common"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

type (
	bannersCache cache.Cache[domain.BannerID, struct{}]

	slotRotationsCache cache.Cache[domain.SlotID, statData]
	rotationsCache     cache.Cache[domain.GroupID, slotRotationsCache]
)

type App struct {
	logger           common.Logger
	storage          Storage
	config           config.CacheConf
	mutSlotBanners   sync.RWMutex
	slotBannersCache cache.Cache[domain.SlotID, bannersCache]
	mut              sync.Mutex
	rotationsCache   rotationsCache
}

type Storage interface {
	CreateBanner(ctx context.Context, entity domain.Banner) error
	GetBanner(ctx context.Context, entityID string) (domain.Banner, error)
	ListBanners(ctx context.Context) ([]domain.Banner, error)
	UpdateBanner(ctx context.Context, entity domain.Banner) error
	DeleteBanner(ctx context.Context, entityID string) error

	CreateGroup(ctx context.Context, entity domain.Group) error
	GetGroup(ctx context.Context, entityID string) (domain.Group, error)
	ListGroups(ctx context.Context) ([]domain.Group, error)
	UpdateGroup(ctx context.Context, entity domain.Group) error
	DeleteGroup(ctx context.Context, entityID string) error

	CreateSlot(ctx context.Context, entity domain.Slot) error
	GetSlot(ctx context.Context, entityID string) (domain.Slot, error)
	ListSlots(ctx context.Context) ([]domain.Slot, error)
	UpdateSlot(ctx context.Context, entity domain.Slot) error
	DeleteSlot(ctx context.Context, entityID string) error

	AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error
	RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error
	GetSlotBanner(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) (domain.Banner, error)
	ListSlotBanners(ctx context.Context, slotID domain.SlotID) ([]domain.Banner, error)

	ListBannerStats(ctx context.Context, groupID domain.GroupID, slotID domain.SlotID) ([]domain.BannerStat, error)
	IncrementRotationClick(
		ctx context.Context,
		groupID domain.GroupID,
		slotID domain.SlotID,
		bannerID domain.BannerID,
	) error
	IncrementRotationShow(
		ctx context.Context,
		groupID domain.GroupID,
		slotID domain.SlotID,
		bannerID domain.BannerID,
	) error
	GetBannerStat(
		ctx context.Context,
		groupID domain.GroupID,
		slotID domain.SlotID,
		bannerID domain.BannerID,
	) (domain.BannerStat, error)
}

func New(logger common.Logger, storage Storage, config config.CacheConf) *App {
	return &App{
		logger:           logger,
		storage:          storage,
		config:           config,
		slotBannersCache: cache.NewLRUCache[domain.SlotID, bannersCache](config.L1Capacity),
		rotationsCache:   cache.NewLRUCache[domain.GroupID, slotRotationsCache](config.L1Capacity),
	}
}

var (
	errCannotBeEmpty            = errors.New("cannot be empty")
	errNoBannersAssignedForSlot = errors.New("no banners assigned for slot")
	errBannerNoAllowedForSlot   = errors.New("banner not allowed for slot")
)
