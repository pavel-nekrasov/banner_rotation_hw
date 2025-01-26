package rotatorapp

import (
	"math"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

type statData struct {
	totalShowCount int64
	bestBannerID   domain.BannerID
	bannerStats    cache.Cache[domain.BannerID, domain.BannerStat]
}

func (s *statData) empty() bool {
	return s.bannerStats.Len() == 0
}

func (s *statData) recalculate() {
	s.totalShowCount = 0
	var logTotalShowCount float64

	s.bannerStats.Range(func(_ domain.BannerID, data domain.BannerStat) {
		s.totalShowCount += data.ShowCount
	})
	logTotalShowCount = 2 * math.Log(float64(s.totalShowCount))

	var maxCoef float64
	s.bannerStats.Range(func(_ domain.BannerID, data domain.BannerStat) {
		coef := data.Ratio() + math.Sqrt(logTotalShowCount/float64(data.ShowCount))
		if coef > maxCoef {
			maxCoef = coef
			s.bestBannerID = data.BannerID
		}
	})
}
