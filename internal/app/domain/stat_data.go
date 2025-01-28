package appdomain

import (
	"math"
	"sync"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

type StatData struct {
	mut            sync.RWMutex
	totalShowCount int64
	bestBannerID   domain.BannerID
	bannerStats    cache.Cache[domain.BannerID, *domain.BannerStat]
}

func NewStatData() *StatData {
	return &StatData{
		bannerStats: cache.NewMapCache[domain.BannerID, *domain.BannerStat](),
	}
}

func (s *StatData) Lock() {
	s.mut.Lock()
}

func (s *StatData) Unlock() {
	s.mut.Unlock()
}

func (s *StatData) Empty() bool {
	return s.bannerStats.Len() == 0
}

func (s *StatData) BestBanner() domain.BannerID {
	return s.bestBannerID
}

func (s *StatData) Recalculate() {
	s.bestBannerID = ""
	s.totalShowCount = 0
	var logTotalShowCount float64

	s.bannerStats.Range(func(_ domain.BannerID, data *domain.BannerStat) {
		s.totalShowCount += data.ShowCount
	})
	logTotalShowCount = 2 * math.Log(float64(s.totalShowCount))

	var maxCoef float64
	s.bannerStats.Range(func(_ domain.BannerID, data *domain.BannerStat) {
		coef := data.Ratio() + math.Sqrt(logTotalShowCount/float64(data.ShowCount))
		if coef > maxCoef {
			maxCoef = coef
			s.bestBannerID = data.BannerID
		}
	})
}

func (s *StatData) Set(bannerID domain.BannerID, bannerStat *domain.BannerStat) {
	s.bannerStats.Set(bannerID, bannerStat)
}

func (s *StatData) Remove(bannerID domain.BannerID) {
	s.bannerStats.Remove(bannerID)
}

func (s *StatData) IncrementClick(bannerID domain.BannerID) {
	data, ok := s.bannerStats.Get(bannerID)
	if ok {
		data.ClickCount++
	}
}

func (s *StatData) IncrementShow(bannerID domain.BannerID) {
	data, ok := s.bannerStats.Get(bannerID)
	if ok {
		data.ShowCount++
	}
}

func (s *StatData) Range(handler func(key domain.BannerID, data *domain.BannerStat)) {
	s.bannerStats.Range(handler)
}
