package storageintegration

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

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
	logger    *logger.Logger
	config    config.ServerConfig
	dbConn    *storage.Connection
	storage   *storage.Storage
	app       *rotatorapp.App
	random    *rand.Rand
	ctx       context.Context
	ctxCancel context.CancelFunc
}

const (
	GroupCount  = 5
	SlotCount   = 10
	BannerCount = 20
)

func TestAppIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AppIntegrationSuite))
}

func (s *AppIntegrationSuite) SetupSuite() {
	s.random = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	s.config = config.NewRotatorConfig("/app/config/server_config.toml")
	s.logger = logger.New(s.config.Logger.Level, s.config.Logger.Output)
	s.dbConn = storage.NewConnection(
		s.config.Storage.Host,
		s.config.Storage.Port,
		s.config.Storage.DBName,
		s.config.Storage.User,
		s.config.Storage.Password,
	)
	s.storage = storage.NewStorage(s.dbConn)
}

func (s *AppIntegrationSuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.app = rotatorapp.New(s.ctx, s.logger, s.storage, nil, s.config.Cache)
	if err := s.dbConn.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	s.cleanupDB()
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
	defer s.ctxCancel()
	defer s.dbConn.Close()

	s.cleanupDB()
}

func (s *AppIntegrationSuite) cleanupDB() {
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE banner_stats CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE slot_banners CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE banners CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE groups CASCADE")
	s.dbConn.DB.Exec(context.Background(), "TRUNCATE slots CASCADE")
}

func (s *AppIntegrationSuite) TestSlotBannersNegative() {
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
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrNoBannersAssignedForSlot)

	err = s.app.ClickBanner(context.Background(), "Group1ID", "Slot1ID", "Banner1ID")
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrBannerNoAllowedForSlot)
	s.Suite.Require().Error(err)

	err = s.app.AddBannerToSlot(context.Background(), "Slot1ID", "Banner1ID")
	s.Suite.Require().NoError(err)

	err = s.app.ClickBanner(context.Background(), "Group1ID", "Slot1ID", "Banner2ID")
	s.Suite.Require().Error(err)
	s.Suite.Require().ErrorIs(err, rotatorapp.ErrBannerNoAllowedForSlot)
}

// в тесте кликаем X раз на один и тот же баннер  в одном и том же слоте,
// а потом показываем баннеры в том же слоте Y раз
//  1. Перебор всех: после большого количества показов, каждый баннер должен быть показан хотя один раз
//  2. Выбор популярных: если на один из баннеров кликают,
//     у него должно быть существенно больше показов чем у остальных
func (s *AppIntegrationSuite) TestSingleLogicCheck() {
	const groupID = domain.GroupID("Group0ID")
	const slotID = domain.SlotID("Slot0ID")
	const bannerID = domain.BannerID("Banner0ID")
	const numberOfClicks = 100
	const numberOfShows = 5000
	const thresholdPopularMin = 750   // мин кол-во показов по баннеру по которому кликали
	const thresholdUnpopularMax = 250 // макс кол-во показова по баннерам по которым не кликали
	const thresholdAtLeastOneShow = 1 // баннеры по которым не кликали должны быть показаны как минимум 1 раз
	const batchSize = 50              // операции делаем пачками по 50 горутин

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

	// загружаем данные по группе/слоту по которым кликали/показывали
	bannerStats, err := s.storage.ListBannerStats(context.Background(), groupID, slotID)
	s.Suite.Require().NoError(err)
	/*
		fmt.Println("database data:")
		for _, bs := range bannerStats {
			fmt.Printf("\tbs: %v\n", bs)
		}*/
	for _, bs := range bannerStats {
		if bs.BannerID == bannerID {
			s.Suite.Require().True(
				bs.ShowCount >= thresholdPopularMin,
				"популярный баннер должен показываться часто",
			)
			s.Suite.Require().Equal(
				numberOfClicks+1, int(bs.ClickCount),
				"кол-во кликов у популярного баннера должно быть правильным",
			)
		} else {
			s.Suite.Require().True(
				bs.ShowCount > thresholdAtLeastOneShow,
				"непопулярный баннер должен быть показан хотя бы раз",
			)
			s.Suite.Require().True(
				bs.ShowCount <= thresholdUnpopularMax,
				"непопулярный баннер должен быть иногда показываться",
			)
			s.Suite.Require().Equal(1, int(bs.ClickCount))
		}
	}

	// берем из БД общее кол-во кликов/показов
	actualShows, actualClicks, err := s.storage.BannerStatTotals(context.Background())
	s.Suite.Require().NoError(err)

	s.Suite.Require().Equal(int64(numberOfClicks), actualClicks, "общее кол-во кликов по всем баннерам")
	s.Suite.Require().Equal(int64(numberOfShows), actualShows, "общее кол-во показов по всем баннерам")
}

// в тесте кликаем X раз случайно на разные баннеры в разных слотах, и показываем баннеры Y раз
// причем клики и показы идут одновременно
//
//	проверка правильности работы в условиях конкурентности
func (s *AppIntegrationSuite) TestMultiLogicCheck() {
	const numberOfClicks = 10000
	const numberOfShows = 100000
	const batchSize = 50 // операции делаем пачками по 50 горутин

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
						domain.GroupID(fmt.Sprintf("Group%vID", s.random.Intn(GroupCount))),
						domain.SlotID(fmt.Sprintf("Slot%vID", s.random.Intn(SlotCount))),
						domain.BannerID(fmt.Sprintf("Banner%vID", s.random.Intn(BannerCount))),
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
						domain.GroupID(fmt.Sprintf("Group%vID", s.random.Intn(GroupCount))),
						domain.SlotID(fmt.Sprintf("Slot%vID", s.random.Intn(SlotCount))),
					)
					s.Suite.Require().NoError(err)
				}()
			}

			wgBatch.Wait()
		}
		wgTasks.Done()
	}()
	wgTasks.Wait()

	// берем из БД общее кол-во кликов/показов
	actualShows, actualClicks, err := s.storage.BannerStatTotals(context.Background())
	s.Suite.Require().NoError(err)

	s.Suite.Require().Equal(int64(numberOfClicks), actualClicks, "общее кол-во кликов по всем баннерам")
	s.Suite.Require().Equal(int64(numberOfShows), actualShows, "общее кол-во показов по всем баннерам")
}
