package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/server/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application interface {
	ClickBanner(ctx context.Context, groupID domain.GroupID, slotID domain.SlotID, bannerID domain.BannerID) error
	SelectBanner(ctx context.Context, groupID domain.GroupID, slotID domain.SlotID) (domain.BannerID, error)
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type BannerService struct {
	logger Logger
	app    Application
	pb.UnimplementedBannersServer
}

func NewBannerService(logger Logger, app Application) *BannerService {
	return &BannerService{logger: logger, app: app}
}

func (s *BannerService) ClickBanner(
	ctx context.Context,
	req *pb.ClickBannerRequest,
) (*empty.Empty, error) {
	groupID := req.GetGroupID()
	if groupID == "" {
		return nil, status.Error(codes.InvalidArgument, "groupID is not specified")
	}

	slotID := req.GetSlotID()
	if slotID == "" {
		return nil, status.Error(codes.InvalidArgument, "slotID is not specified")
	}

	bannerID := req.GetBannerID()
	if bannerID == "" {
		return nil, status.Error(codes.InvalidArgument, "bannerID is not specified")
	}

	err := s.app.ClickBanner(ctx, domain.GroupID(groupID), domain.SlotID(slotID), domain.BannerID(bannerID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *BannerService) SelectBanner(
	ctx context.Context,
	req *pb.SelectBannerRequest,
) (*pb.SelectBannerResponse, error) {
	groupID := req.GetGroupID()
	if groupID == "" {
		return nil, status.Error(codes.InvalidArgument, "groupID is not specified")
	}

	slotID := req.GetSlotID()
	if slotID == "" {
		return nil, status.Error(codes.InvalidArgument, "slotID is not specified")
	}

	bannerID, err := s.app.SelectBanner(ctx, domain.GroupID(req.GroupID), domain.SlotID(req.SlotID))
	if err != nil {
		return nil, err
	}

	return &pb.SelectBannerResponse{BannerID: string(bannerID)}, nil
}
