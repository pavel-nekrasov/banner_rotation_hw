package storageintegration

import (
	"context"
	"log"
	"testing"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
	"github.com/stretchr/testify/suite"
)

type StorageIntegrationSuite struct {
	suite.Suite
	dbConn  *storage.Connection
	storage *storage.Storage
}

func TestStorageIntegrationSuite(t *testing.T) {
	suite.Run(t, new(StorageIntegrationSuite))
}

func (s *StorageIntegrationSuite) SetupSuite() {
	config := config.NewRotatorConfig("/app/config/server_config.toml")
	s.dbConn = storage.NewConnection(config.Storage)
	s.storage = storage.NewStorage(s.dbConn)
}

func (s *StorageIntegrationSuite) SetupTest() {
	if err := s.dbConn.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (s *StorageIntegrationSuite) TearDownTest() {
	defer s.dbConn.Close()
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE banner_stats CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE slot_banners CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE banners CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE groups CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE slots CASCADE")
}

func (s *StorageIntegrationSuite) TestBannerCrud() { //nolint:dupl
	itemID := domain.BannerID("banner1")
	item := domain.Banner{
		ID:          itemID,
		Description: "banner description",
	}
	err := s.storage.CreateBanner(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetBanner(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "banner description upd"
	err = s.storage.UpdateBanner(context.Background(), domain.Banner{ID: itemID, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetBanner(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(itemID, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteBanner(context.Background(), itemID)
	s.Require().NoError(err)

	_, err = s.storage.GetBanner(context.Background(), itemID)
	s.Require().Error(err)
}

func (s *StorageIntegrationSuite) TestGroupCrud() { //nolint:dupl
	itemID := domain.GroupID("group1")
	item := domain.Group{
		ID:          itemID,
		Description: "group description",
	}
	err := s.storage.CreateGroup(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetGroup(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "group description upd"
	err = s.storage.UpdateGroup(context.Background(), domain.Group{ID: itemID, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetGroup(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(itemID, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteGroup(context.Background(), itemID)
	s.Require().NoError(err)

	_, err = s.storage.GetGroup(context.Background(), itemID)
	s.Require().Error(err)
}

func (s *StorageIntegrationSuite) TestSlotCrud() { //nolint:dupl
	itemID := domain.SlotID("slot1")
	item := domain.Slot{
		ID:          itemID,
		Description: "slot description",
	}
	err := s.storage.CreateSlot(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetSlot(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "slot description upd"
	err = s.storage.UpdateSlot(context.Background(), domain.Slot{ID: itemID, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetSlot(context.Background(), itemID)
	s.Require().NoError(err)
	s.Require().Equal(itemID, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteSlot(context.Background(), itemID)
	s.Require().NoError(err)

	_, err = s.storage.GetSlot(context.Background(), itemID)
	s.Require().Error(err)
}

func (s *StorageIntegrationSuite) TestSlotBanners() {
	slotID := domain.SlotID("slot1")
	slot := domain.Slot{
		ID:          slotID,
		Description: "description",
	}
	err := s.storage.CreateSlot(context.Background(), slot)
	s.Require().NoError(err)

	bannerID := domain.BannerID("banner1")
	banner := domain.Banner{
		ID:          bannerID,
		Description: "description",
	}
	err = s.storage.CreateBanner(context.Background(), banner)
	s.Require().NoError(err)

	err = s.storage.AddBannerToSlot(context.Background(), slotID, bannerID)
	s.Require().NoError(err)

	banners, err := s.storage.ListSlotBanners(context.Background(), slotID)
	s.Require().NoError(err)
	s.Require().Equal(1, len(banners))
	s.Require().Equal(bannerID, banners[0].ID)

	err = s.storage.RemoveBannerFromSlot(context.Background(), slotID, bannerID)
	s.Require().NoError(err)

	banners, err = s.storage.ListSlotBanners(context.Background(), slotID)
	s.Require().NoError(err)
	s.Require().Equal(0, len(banners))
}

func (s *StorageIntegrationSuite) TestRotations() {
	slotID := domain.SlotID("slot1")
	slot := domain.Slot{
		ID:          slotID,
		Description: "description",
	}
	err := s.storage.CreateSlot(context.Background(), slot)
	s.Require().NoError(err)

	bannerID := domain.BannerID("banner1")
	banner := domain.Banner{
		ID:          bannerID,
		Description: "description",
	}
	err = s.storage.CreateBanner(context.Background(), banner)
	s.Require().NoError(err)

	banner2ID := domain.BannerID("banner2")
	banner2 := domain.Banner{
		ID:          banner2ID,
		Description: "description",
	}
	err = s.storage.CreateBanner(context.Background(), banner2)
	s.Require().NoError(err)

	groupID := domain.GroupID("group1")
	group := domain.Group{
		ID:          groupID,
		Description: "description",
	}

	err = s.storage.CreateGroup(context.Background(), group)
	s.Require().NoError(err)

	err = s.storage.AddBannerToSlot(context.Background(), slotID, bannerID)
	s.Require().NoError(err)
	err = s.storage.AddBannerToSlot(context.Background(), slotID, banner2ID)
	s.Require().NoError(err)

	err = s.storage.IncrementBannerStatClick(context.Background(), groupID, slotID, bannerID)
	s.Require().NoError(err)

	err = s.storage.IncrementBannerStatShow(context.Background(), groupID, slotID, banner2ID)
	s.Require().NoError(err)

	rotations, err := s.storage.ListBannerStats(context.Background(), groupID, slotID)
	s.Require().NoError(err)
	s.Require().Equal(2, len(rotations))

	rotation, err := s.storage.GetBannerStat(context.Background(), groupID, slotID, bannerID)
	s.Require().NoError(err)
	s.Require().Equal(bannerID, rotation.BannerID)
	s.Require().Equal(int64(2), rotation.ClickCount)
	s.Require().Equal(int64(1), rotation.ShowCount)

	rotation, err = s.storage.GetBannerStat(context.Background(), groupID, slotID, banner2ID)
	s.Require().NoError(err)
	s.Require().Equal(banner2ID, rotation.BannerID)
	s.Require().Equal(int64(1), rotation.ClickCount)
	s.Require().Equal(int64(2), rotation.ShowCount)

	err = s.storage.IncrementBannerStatClick(context.Background(), groupID, slotID, bannerID)
	s.Require().NoError(err)
	err = s.storage.IncrementBannerStatShow(context.Background(), groupID, slotID, bannerID)
	s.Require().NoError(err)

	rotation, err = s.storage.GetBannerStat(context.Background(), groupID, slotID, bannerID)
	s.Require().NoError(err)
	s.Require().Equal(bannerID, rotation.BannerID)
	s.Require().Equal(int64(3), rotation.ClickCount)
	s.Require().Equal(int64(2), rotation.ShowCount)
}
