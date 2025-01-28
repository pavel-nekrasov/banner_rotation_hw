package appdomain

import "github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"

type NotificationEvent struct {
	EventType string
	GroupID   domain.GroupID
	SlotID    domain.SlotID
	BannerID  domain.BannerID
	Timestamp int64
}
