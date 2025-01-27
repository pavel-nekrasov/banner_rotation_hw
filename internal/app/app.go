package rotatorapp

import (
	"context"
	"errors"
	"fmt"
	"sync"

	appdomain "github.com/pavel-nekrasov/banner_rotation_hw/internal/app/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/common"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

type (
	bannersCache cache.Cache[domain.BannerID, struct{}]

	slotRotationsCache cache.Cache[domain.SlotID, *appdomain.StatData]
	rotationsCache     cache.Cache[domain.GroupID, slotRotationsCache]
)

type App struct {
	logger           common.Logger
	storage          Storage
	config           config.CacheConf
	mutSlotBanners   sync.RWMutex
	groupsCache      cache.Cache[domain.GroupID, domain.Group]
	slotsCache       cache.Cache[domain.SlotID, domain.Slot]
	bannersCache     cache.Cache[domain.BannerID, domain.Banner]
	slotBannersCache cache.Cache[domain.SlotID, bannersCache]
	mut              sync.RWMutex
	rotationsCache   rotationsCache
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

	ListBannerStats(ctx context.Context, groupID domain.GroupID, slotID domain.SlotID) ([]domain.BannerStat, error)
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

func New(logger common.Logger, storage Storage, config config.CacheConf) *App {
	return &App{
		logger:           logger,
		storage:          storage,
		config:           config,
		groupsCache:      cache.NewLRUCache[domain.GroupID, domain.Group](config.L1Capacity),
		slotsCache:       cache.NewLRUCache[domain.SlotID, domain.Slot](config.L1Capacity),
		bannersCache:     cache.NewLRUCache[domain.BannerID, domain.Banner](config.L1Capacity),
		slotBannersCache: cache.NewLRUCache[domain.SlotID, bannersCache](config.L1Capacity),
		rotationsCache:   cache.NewLRUCache[domain.GroupID, slotRotationsCache](config.L1Capacity),
	}
}

var (
	ErrCannotBeEmpty            = errors.New("cannot be empty")
	ErrNoBannersAssignedForSlot = errors.New("no banners assigned for slot")
	ErrBannerNoAllowedForSlot   = errors.New("banner not allowed for slot")
)

func (a *App) Debug() {
	a.rotationsCache.Range(func(key domain.GroupID, data slotRotationsCache) {
		fmt.Printf("\ngroup: %v\n", key)
		data.Range(func(key domain.SlotID, data *appdomain.StatData) {
			fmt.Printf("\tslot: %v\n", key)
			fmt.Printf("\t\tLast best banner: %v\n", data.BestBanner())
			data.Range(func(key domain.BannerID, data *domain.BannerStat) {
				fmt.Printf("\t\t\tbanner: %v  shows: %v  clicks: %v\n", key, data.ShowCount, data.ClickCount)
			})
		})
	})
}
