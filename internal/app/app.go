package rotatorapp

import (
	"context"
	"errors"
	"fmt"

	appdomain "github.com/pavel-nekrasov/banner_rotation_hw/internal/app/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/common"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

type (
	bannersCache cache.Cache[domain.BannerID, struct{}]
)

type App struct {
	config              config.CacheConf
	logger              common.Logger
	storage             Storage
	publisher           Publisher
	groupSlotKeys       *cache.ObjectMutex[appdomain.GroupSlotKey]
	slotKeys            *cache.ObjectMutex[domain.SlotID]
	bannersCache        cache.Cache[domain.BannerID, domain.Banner]
	groupsCache         cache.Cache[domain.GroupID, domain.Group]
	slotsCache          cache.Cache[domain.SlotID, domain.Slot]
	allowedBannersCache cache.Cache[domain.SlotID, bannersCache]
	rotationsCache      cache.Cache[appdomain.GroupSlotKey, *appdomain.StatData]
}

type Storage interface {
	CreateBanner(ctx context.Context, entity domain.Banner) error
	GetBanner(ctx context.Context, entityID domain.BannerID) (domain.Banner, error)
	ListBanners(ctx context.Context) ([]domain.Banner, error)
	UpdateBanner(ctx context.Context, entity domain.Banner) error
	DeleteBanner(ctx context.Context, entityID domain.BannerID) error

	CreateGroup(ctx context.Context, entity domain.Group) error
	GetGroup(ctx context.Context, entityID domain.GroupID) (domain.Group, error)
	ListGroups(ctx context.Context) ([]domain.Group, error)
	UpdateGroup(ctx context.Context, entity domain.Group) error
	DeleteGroup(ctx context.Context, entityID domain.GroupID) error

	CreateSlot(ctx context.Context, entity domain.Slot) error
	GetSlot(ctx context.Context, entityID domain.SlotID) (domain.Slot, error)
	ListSlots(ctx context.Context) ([]domain.Slot, error)
	UpdateSlot(ctx context.Context, entity domain.Slot) error
	DeleteSlot(ctx context.Context, entityID domain.SlotID) error

	AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error
	RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error
	GetSlotBanner(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) (domain.Banner, error)
	ListSlotBanners(ctx context.Context, slotID domain.SlotID) ([]domain.Banner, error)

	ListBannerStats(
		ctx context.Context,
		groupID domain.GroupID,
		slotID domain.SlotID,
	) ([]domain.BannerStat, error)
	IncrementBannerStatClick(
		ctx context.Context,
		groupID domain.GroupID,
		slotID domain.SlotID,
		bannerID domain.BannerID,
	) error
	IncrementBannerStatShow(
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

type Publisher interface {
	Publish(data []byte) error
}

func New(
	ctx context.Context,
	logger common.Logger,
	storage Storage,
	publisher Publisher,
	config config.CacheConf,
) *App {
	return &App{
		logger:              logger,
		storage:             storage,
		publisher:           publisher,
		config:              config,
		groupSlotKeys:       cache.NewObjectMutex[appdomain.GroupSlotKey](ctx, logger),
		slotKeys:            cache.NewObjectMutex[domain.SlotID](ctx, logger),
		groupsCache:         cache.NewLRUCache[domain.GroupID, domain.Group](config.L1Capacity),
		slotsCache:          cache.NewLRUCache[domain.SlotID, domain.Slot](config.L1Capacity),
		bannersCache:        cache.NewLRUCache[domain.BannerID, domain.Banner](config.L1Capacity),
		allowedBannersCache: cache.NewLRUCache[domain.SlotID, bannersCache](config.L1Capacity),
		rotationsCache:      cache.NewLRUCache[appdomain.GroupSlotKey, *appdomain.StatData](config.L1Capacity),
	}
}

var (
	ErrCannotBeEmpty            = errors.New("cannot be empty")
	ErrNoBannersAssignedForSlot = errors.New("no banners assigned for slot")
	ErrBannerNoAllowedForSlot   = errors.New("banner not allowed for slot")
)

func (a *App) CacheDebugInfo() {
	a.rotationsCache.Range(func(key appdomain.GroupSlotKey, data *appdomain.StatData) {
		fmt.Printf("\ngroup: %v\n", key.GroupID)
		fmt.Printf("\tslot: %v\n", key.SlotID)
		fmt.Printf("\t\tLast best banner: %v\n", data.BestBanner())
		data.Range(func(key domain.BannerID, data *domain.BannerStat) {
			fmt.Printf("\t\t\tbanner: %v  shows: %v  clicks: %v\n", key, data.ShowCount, data.ClickCount)
		})
	})
}
