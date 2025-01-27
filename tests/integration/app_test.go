package storageintegration

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
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
	GroupCount  = 10
	SlotCount   = 20
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
	s.storage.DB.ExecContext(context.Background(), "TRUNCATE banner_stats CASCADE")
	s.storage.DB.ExecContext(context.Background(), "TRUNCATE slot_banners CASCADE")
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

// в тесте кликаем X раз на один и тот же баннер (Z) в одном и том же слоте,
// а потом показываем баннеры в том же слоте Y раз
//  1. Перебор всех: после большого количества показов, каждый баннер должен быть показан хотя один раз (Z)
//  2. Выбор популярных: если на один из баннеров кликают,
//     у него должно быть существенно больше показов чем у остальных
func (s *AppIntegrationSuite) TestSingleLogicCheckCrud() {
	const groupID = domain.GroupID("Group0ID")
	const slotID = domain.SlotID("Slot0ID")
	const bannerID = domain.BannerID("Banner0ID")
	const numberOfClicks = 100
	const numberOfShows = 5000
	const thresholdClickedBanner = 400 // мин кол-во показов по баннеру по которому кликали
	const thresholdNotClicked = 100    // макс кол-во показова по баннерам по которым не кликали
	const thresholdAtLeastOneShow = 1  // баннеры по которым не кликали должны быть показаны как минимум 1 раз
	const batchSize = 50               // операции делаем пачками по 50 горутин

	// добавляем все баннеры в показ для первого слота
	for i := 0; i < BannerCount; i++ {
		err := s.app.AddBannerToSlot(context.Background(), slotID, domain.BannerID(fmt.Sprintf("Banner%vID", i)))
		s.Suite.Require().NoError(err)
	}

	// по  баннеру 0 делаем numberOfClicks кликов
	for i := 0; i < numberOfClicks; i += batchSize {
		var wg sync.WaitGroup

		for j := 0; j < batchSize; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.app.ClickBanner(context.Background(), groupID, slotID, bannerID)
			}()
		}

		wg.Wait()
	}

	// делаем numberOfShows показов
	for i := 0; i < numberOfShows; i += batchSize {
		var wg sync.WaitGroup

		for j := 0; j < batchSize; j++ {
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
	/*
		fmt.Println("database data:")
		for _, bs := range bannerStats {
			fmt.Printf("\tbs: %v\n", bs)
		}
	*/
	for _, bs := range bannerStats {
		if bs.BannerID == bannerID {
			s.Suite.Require().True(bs.ShowCount >= thresholdClickedBanner, bs.BannerID)
			s.Suite.Require().Equal(numberOfClicks+1, int(bs.ClickCount))
		} else {
			s.Suite.Require().True(bs.ShowCount > thresholdAtLeastOneShow, bs.BannerID)
			s.Suite.Require().True(bs.ShowCount <= thresholdNotClicked, bs.BannerID)
			s.Suite.Require().Equal(1, int(bs.ClickCount))
		}
		totalShows += bs.ShowCount
		totalClicks += bs.ClickCount
	}

	// проверка целостности данных:
	// кликали X раз - значит сумма кликов сохраненная в базе для соответствующих записей
	// должна быть X + кол-во баннеров для показа
	// по умолчанию кол-во кликов равно 1 даже если по баннеру не кликали, чтобы избежать деления на 0
	// показывали Y раз - начит сумма показов сохраненная в базе для соответствующих записей
	// должна быть Y + кол-во баннеров для показа
	// по умолчанию кол-во показов равно 1 даже если по баннер не разу не показывали, чтобы избежать деления на 0
	s.Suite.Require().Equal(BannerCount+numberOfClicks, int(totalClicks))
	s.Suite.Require().Equal(BannerCount+numberOfShows, int(totalShows))
}

// в тесте кликаем X раз случайно на разные баннеры в разных слотах
//
//	проверка правильности работы в условиях конкурентности
func (s *AppIntegrationSuite) TestMultiLogicCheckCrud() {
	const numberOfClicks = 1000
	const numberOfShows = 10000
	const batchSize = 50 // операции делаем пачками по 50 горутин

	// создаем рандомайзеры
	randBanner, err := rand.Int(rand.Reader, big.NewInt(BannerCount))
	s.Suite.Require().NoError(err)

	randSlot, err := rand.Int(rand.Reader, big.NewInt(SlotCount))
	s.Suite.Require().NoError(err)

	randGroup, err := rand.Int(rand.Reader, big.NewInt(GroupCount))
	s.Suite.Require().NoError(err)

	// добавляем все баннеры в показ для всех слотов
	for j := 0; j < SlotCount; j++ {
		for i := 0; i < BannerCount; i++ {
			err := s.app.AddBannerToSlot(
				context.Background(),
				domain.SlotID(fmt.Sprintf("Slot%vID", j)),
				domain.BannerID(fmt.Sprintf("Banner%vID", i)))
			s.Suite.Require().NoError(err)
		}
	}

	var wgTasks sync.WaitGroup
	// делаем numberOfClicks кликов
	wgTasks.Add(1)
	go func() {
		for i := 0; i < numberOfClicks; i += batchSize {
			var wgBatch sync.WaitGroup
			for j := 0; j < batchSize; j++ {
				wgBatch.Add(1)
				go func() {
					defer wgBatch.Done()
					s.app.ClickBanner(
						context.Background(),
						domain.GroupID(fmt.Sprintf("Group%vID", randGroup.Int64())),
						domain.SlotID(fmt.Sprintf("Slot%vID", randSlot.Int64())),
						domain.BannerID(fmt.Sprintf("Banner%vID", randBanner.Int64())),
					)
				}()
			}

			wgBatch.Wait()
		}
		wgTasks.Done()
	}()

	// делаем numberOfShows показов
	wgTasks.Add(1)
	go func() {
		for i := 0; i < numberOfShows; i += batchSize {
			var wgBatch sync.WaitGroup

			for j := 0; j < batchSize; j++ {
				wgBatch.Add(1)

				go func() {
					defer wgBatch.Done()
					_, err := s.app.SelectBanner(
						context.Background(),
						domain.GroupID(fmt.Sprintf("Group%vID", randGroup.Int64())),
						domain.SlotID(fmt.Sprintf("Slot%vID", randSlot.Int64())),
					)
					s.Suite.Require().NoError(err)
				}()
			}

			wgBatch.Wait()
		}
		wgTasks.Done()
	}()
	wgTasks.Wait()
	totalShows := int64(0)
	totalClicks := int64(0)

	for i := 0; i < GroupCount; i++ {
		for j := 0; j < SlotCount; j++ {
			bannerStats, err := s.storage.ListBannerStats(
				context.Background(),
				domain.GroupID(fmt.Sprintf("Group%vID", i)),
				domain.SlotID(fmt.Sprintf("Slot%vID", j)),
			)
			s.Suite.Require().NoError(err)
			for _, bs := range bannerStats {
				totalShows += bs.ShowCount
				totalClicks += bs.ClickCount
			}
		}
	}

	// проверка целостности данных:
	// кликали X раз - значит сумма кликов сохраненная в базе для соответствующих записей
	// должна быть X + кол-во баннеров для показа
	// по умолчанию кол-во кликов равно 1 даже если по баннеру не кликали, чтобы избежать деления на 0
	// показывали Y раз - значит сумма показов сохраненная в базе для соответствующих записей
	// должна быть Y + кол-во баннеров для показа
	// по умолчанию кол-во показов равно 1 даже если  баннер ни разу не показывали, чтобы избежать деления на 0
	s.Suite.Require().Equal(BannerCount+numberOfClicks, int(totalClicks))
	s.Suite.Require().Equal(BannerCount+numberOfShows, int(totalShows))
}
