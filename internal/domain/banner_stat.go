package domain

type BannerStat struct {
	BannerID   BannerID
	ShowCount  int64
	ClickCount int64
}

func (b *BannerStat) Ratio() float64 {
	return float64(b.ClickCount) / float64(b.ShowCount)
}
