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
	AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error
	RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error

	CreateBanner(ctx context.Context, entity domain.Banner) (domain.Banner, error)
	UpdateBanner(ctx context.Context, entity domain.Banner) (domain.Banner, error)
	GetBanner(ctx context.Context, id domain.BannerID) (domain.Banner, error)
	DeleteBanner(ctx context.Context, id domain.BannerID) error

	CreateGroup(ctx context.Context, entity domain.Group) (domain.Group, error)
	UpdateGroup(ctx context.Context, entity domain.Group) (domain.Group, error)
	GetGroup(ctx context.Context, id domain.GroupID) (domain.Group, error)
	DeleteGroup(ctx context.Context, id domain.GroupID) error

	CreateSlot(ctx context.Context, entity domain.Slot) (domain.Slot, error)
	UpdateSlot(ctx context.Context, entity domain.Slot) (domain.Slot, error)
	GetSlot(ctx context.Context, id domain.SlotID) (domain.Slot, error)
	DeleteSlot(ctx context.Context, id domain.SlotID) error
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

	bannerID, err := s.app.SelectBanner(ctx, domain.GroupID(groupID), domain.SlotID(slotID))
	if err != nil {
		return nil, err
	}

	return &pb.SelectBannerResponse{BannerID: string(bannerID)}, nil
}

func (s *BannerService) AddBannerToSlot(
	ctx context.Context,
	req *pb.AddBannerToSlotRequest,
) (*empty.Empty, error) {
	slotID := req.GetSlotID()
	if slotID == "" {
		return nil, status.Error(codes.InvalidArgument, "slotID is not specified")
	}

	bannerID := req.GetBannerID()
	if bannerID == "" {
		return nil, status.Error(codes.InvalidArgument, "bannerID is not specified")
	}

	err := s.app.AddBannerToSlot(ctx, domain.SlotID(slotID), domain.BannerID(bannerID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *BannerService) RemoveBannerFromSlot(
	ctx context.Context,
	req *pb.RemoveBannerFromSlotRequest,
) (*empty.Empty, error) {
	slotID := req.GetSlotID()
	if slotID == "" {
		return nil, status.Error(codes.InvalidArgument, "slotID is not specified")
	}

	bannerID := req.GetBannerID()
	if bannerID == "" {
		return nil, status.Error(codes.InvalidArgument, "bannerID is not specified")
	}

	err := s.app.RemoveBannerFromSlot(ctx, domain.SlotID(slotID), domain.BannerID(bannerID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *BannerService) CreateBanner(
	ctx context.Context,
	req *pb.BannerRequest,
) (*pb.SingleBannerResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.CreateBanner(ctx, domain.Banner{
		ID:          domain.BannerID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleBannerResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) UpdateBanner(
	ctx context.Context,
	req *pb.BannerRequest,
) (*pb.SingleBannerResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.UpdateBanner(ctx, domain.Banner{
		ID:          domain.BannerID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleBannerResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) GetBanner(
	ctx context.Context,
	req *pb.BannerIdRequest,
) (*pb.SingleBannerResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	entity, err := s.app.GetBanner(ctx, domain.BannerID(ID))
	if err != nil {
		return nil, err
	}

	return &pb.SingleBannerResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) DeleteBanner(
	ctx context.Context,
	req *pb.BannerIdRequest,
) (*empty.Empty, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	err := s.app.DeleteBanner(ctx, domain.BannerID(ID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *BannerService) CreateGroup(
	ctx context.Context,
	req *pb.GroupRequest,
) (*pb.SingleGroupResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.CreateGroup(ctx, domain.Group{
		ID:          domain.GroupID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleGroupResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) UpdateGroup(
	ctx context.Context,
	req *pb.GroupRequest,
) (*pb.SingleGroupResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.UpdateGroup(ctx, domain.Group{
		ID:          domain.GroupID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleGroupResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) GetGroup(
	ctx context.Context,
	req *pb.GroupIdRequest,
) (*pb.SingleGroupResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	entity, err := s.app.GetGroup(ctx, domain.GroupID(ID))
	if err != nil {
		return nil, err
	}

	return &pb.SingleGroupResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) DeleteGroup(
	ctx context.Context,
	req *pb.GroupIdRequest,
) (*empty.Empty, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	err := s.app.DeleteGroup(ctx, domain.GroupID(ID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *BannerService) CreateSlot(
	ctx context.Context,
	req *pb.SlotRequest,
) (*pb.SingleSlotResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.CreateSlot(ctx, domain.Slot{
		ID:          domain.SlotID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleSlotResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) UpdateSlot(
	ctx context.Context,
	req *pb.SlotRequest,
) (*pb.SingleSlotResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is not specified")
	}

	entity, err := s.app.UpdateSlot(ctx, domain.Slot{
		ID:          domain.SlotID(ID),
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SingleSlotResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) GetSlot(
	ctx context.Context,
	req *pb.SlotIdRequest,
) (*pb.SingleSlotResponse, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	entity, err := s.app.GetSlot(ctx, domain.SlotID(ID))
	if err != nil {
		return nil, err
	}

	return &pb.SingleSlotResponse{
		ID:          string(entity.ID),
		Description: entity.Description,
	}, nil
}

func (s *BannerService) DeleteSlot(
	ctx context.Context,
	req *pb.SlotIdRequest,
) (*empty.Empty, error) {
	ID := req.GetID()
	if ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}

	err := s.app.DeleteSlot(ctx, domain.SlotID(ID))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
