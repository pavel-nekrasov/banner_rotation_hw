package integration

import (
	"context"
	"log"
	"testing"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
	"github.com/stretchr/testify/suite"
)

type StorageIntegrationSuite struct {
	suite.Suite
	storage *storage.Storage
}

func TestStorageIntegrationSuite(t *testing.T) {
	suite.Run(t, new(StorageIntegrationSuite))
}

func (s *StorageIntegrationSuite) SetupSuite() {
	config := config.NewServerConfig("/app/config/server_config.toml")
	s.storage = storage.New(config.Storage.Host, config.Storage.Port, config.Storage.DBName, config.Storage.User, config.Storage.Password)
}

func (s *StorageIntegrationSuite) SetupTest() {
	if err := s.storage.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (s *StorageIntegrationSuite) TearDownTest() {
	defer s.storage.Close(context.Background())
	s.storage.Db.ExecContext(context.Background(), "TRUNCATE banners CASCADE")
	s.storage.Db.ExecContext(context.Background(), "TRUNCATE groups CASCADE")
	s.storage.Db.ExecContext(context.Background(), "TRUNCATE slots CASCADE")
}

func (s *StorageIntegrationSuite) TestBannerCrud() {
	itemId := "banner1"
	item := models.Banner{
		ID:          itemId,
		Description: "description",
	}
	err := s.storage.CreateBanner(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetBanner(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "description updated"
	err = s.storage.UpdateBanner(context.Background(), models.Banner{ID: itemId, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetBanner(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(itemId, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteBanner(context.Background(), itemId)
	s.Require().NoError(err)

	_, err = s.storage.GetBanner(context.Background(), itemId)
	s.Require().Error(err)
}

func (s *StorageIntegrationSuite) TestGroupCrud() {
	itemId := "group1"
	item := models.Group{
		ID:          itemId,
		Description: "description",
	}
	err := s.storage.CreateGroup(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetGroup(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "description updated"
	err = s.storage.UpdateGroup(context.Background(), models.Group{ID: itemId, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetGroup(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(itemId, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteGroup(context.Background(), itemId)
	s.Require().NoError(err)

	_, err = s.storage.GetGroup(context.Background(), itemId)
	s.Require().Error(err)
}

func (s *StorageIntegrationSuite) TestSlotCrud() {
	itemId := "slot1"
	item := models.Slot{
		ID:          itemId,
		Description: "description",
	}
	err := s.storage.CreateSlot(context.Background(), item)
	s.Require().NoError(err)

	dbItem, err := s.storage.GetSlot(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(item.ID, dbItem.ID)
	s.Require().Equal(item.Description, dbItem.Description)

	updatedDescription := "description updated"
	err = s.storage.UpdateSlot(context.Background(), models.Slot{ID: itemId, Description: updatedDescription})
	s.Require().NoError(err)

	dbItem, err = s.storage.GetSlot(context.Background(), itemId)
	s.Require().NoError(err)
	s.Require().Equal(itemId, dbItem.ID)
	s.Require().Equal(updatedDescription, dbItem.Description)

	err = s.storage.DeleteSlot(context.Background(), itemId)
	s.Require().NoError(err)

	_, err = s.storage.GetSlot(context.Background(), itemId)
	s.Require().Error(err)
}
