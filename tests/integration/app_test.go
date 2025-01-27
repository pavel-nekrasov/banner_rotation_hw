package storageintegration

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	rotatorapp "github.com/pavel-nekrasov/banner_rotation_hw/internal/app"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/logger"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
	"github.com/stretchr/testify/suite"
)

type AppIntegrationSuite struct {
	suite.Suite
	storage *storage.Storage
	app     *rotatorapp.App
}

const (
	SlotCount   = 20
	GroupCount  = 10
	BannerCount = 50
)

func TestSAppIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AppIntegrationSuite))
}

func (s *AppIntegrationSuite) SetupSuite() {
	config := config.NewRotatorConfig("/app/config/server_config.toml")
	logg := logger.New(config.Logger.Level, config.Logger.Output)
	s.storage = storage.New(
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.DBName,
		config.Storage.User,
		config.Storage.Password,
	)
	s.app = rotatorapp.New(logg, s.storage, config.Cache)
}

func (s *AppIntegrationSuite) SetupTest() {
	if err := s.storage.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < GroupCount; i++ {
		s.storage.CreateGroup(
			context.Background(),
			domain.Group{
				ID:          domain.GroupID(fmt.Sprintf("Group%vID", i)),
				Description: fmt.Sprintf("Group description %v", i),
			})
	}

	for i := 0; i < SlotCount; i++ {
		s.storage.CreateSlot(
			context.Background(),
			domain.Slot{
				ID:          domain.SlotID(fmt.Sprintf("Slot%vID", i)),
				Description: fmt.Sprintf("Slot description %v", i),
			})
	}

	for i := 0; i < BannerCount; i++ {
		s.storage.CreateBanner(
			context.Background(),
			domain.Banner{
				ID:          domain.BannerID(fmt.Sprintf("Banner%vID", i)),
				Description: fmt.Sprintf("Banner description %v", i),
			})
	}
}

func (s *AppIntegrationSuite) TearDownTest() {
	defer s.storage.Close(context.Background())
	s.app.Debug()
	s.storage.DB.ExecContext(context.Background(), "TRUNCATE banners CASCADE")
	s.storage.DB.ExecContext(context.Background(), "TRUNCATE groups CASCADE")
	s.storage.DB.ExecContext(context.Background(), "TRUNCATE slots CASCADE")
}

func (s *AppIntegrationSuite) TestSlotBannersNegativeCrud() {
	var notFoundErr customerrors.NotFound

	err := s.app.AddBannerToSlot(context.Background(), "wrong slot", "Banner0ID")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorAs(err, &notFoundErr)

	err = s.app.AddBannerToSlot(context.Background(), "Slot1ID", "wrong banner")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorAs(err, &notFoundErr)

	err = s.app.ClickBanner(context.Background(), "GroupID", "wrong slot", "wrong banner")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorAs(err, &notFoundErr)

	_, err = s.app.SelectBanner(context.Background(), "GroupID", "wrong slot")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorAs(err, &notFoundErr)

	_, err = s.app.SelectBanner(context.Background(), "Group1ID", "Slot1ID")
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrNoBannersAssignedForSlot)
	s.Suite.Require().Error(err)

	err = s.app.ClickBanner(context.Background(), "Group1ID", "Slot1ID", "Banner1ID")
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrBannerNoAllowedForSlot)
	s.Suite.Require().Error(err)

	err = s.app.AddBannerToSlot(context.Background(), "Slot1ID", "Banner1ID")
	s.Suite.Require().NoError(err)

	err = s.app.ClickBanner(context.Background(), "Group1ID", "Slot1ID", "Banner2ID")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrBannerNoAllowedForSlot)
}

func (s *AppIntegrationSuite) TestStatisticCheckCrud() {
	groupID := domain.GroupID("Group0ID")
	slotID := domain.SlotID("Slot0ID")
	bannerID := domain.BannerID("Banner0ID")
	numberOfClicks := 100
	numberOfShows := 1000
	shiftSize := 50
	// добавляем баннеры в показ для первого слота
	for i := 0; i < BannerCount; i++ {
		err := s.app.AddBannerToSlot(context.Background(), slotID, domain.BannerID(fmt.Sprintf("Banner%vID", i)))
		s.Suite.Require().NoError(err)
	}

	// по  баннеру 0 делаем numberOfClicks кликов
	for i := 0; i < numberOfClicks; i += shiftSize {
		var wg sync.WaitGroup

		for j := 0; j < shiftSize; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.app.ClickBanner(context.Background(), groupID, slotID, bannerID)
			}()
		}

		wg.Wait()
	}

	// делаем numberOfShows показов
	for i := 0; i < numberOfShows; i += shiftSize {
		var wg sync.WaitGroup

		for j := 0; j < shiftSize; j++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				_, err := s.app.SelectBanner(context.Background(), groupID, slotID)
				s.Suite.Require().NoError(err)
			}()
		}

		wg.Wait()
	}

	bannerStats, err := s.storage.ListBannerStats(context.Background(), groupID, slotID)
	s.Suite.Require().NoError(err)

	totalShows := int64(0)
	totalClicks := int64(0)
	fmt.Println("database data:")
	for _, bs := range bannerStats {
		fmt.Printf("\tbs: %v\n", bs)
	}
	for _, bs := range bannerStats {
		//
		if bs.BannerID == bannerID {
			s.Suite.Require().True(bs.ShowCount >= 100, bs.BannerID)
			s.Suite.Require().Equal(numberOfClicks+1, int(bs.ClickCount))
		} else {
			s.Suite.Require().True(bs.ShowCount >= 2, bs.BannerID)
			s.Suite.Require().True(bs.ShowCount <= 50, bs.BannerID)
			s.Suite.Require().Equal(1, int(bs.ClickCount))
		}
		totalShows += bs.ShowCount
		totalClicks += bs.ClickCount
	}
	s.Suite.Require().Equal(BannerCount+numberOfClicks, int(totalClicks))
	s.Suite.Require().Equal(BannerCount+numberOfShows, int(totalShows))
}
