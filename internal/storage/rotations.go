package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (s *Storage) ListRotations(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) ([]domain.Rotation, error) {
	result := make([]domain.Rotation, 0)
	rows, err := s.DB.QueryContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		INNER JOIN slot_banners sb ON sb.banner_id = r.banner_id AND sb.slot_id = r.slot_id 
		WHERE r.group_id = $1 AND r.slot_id = $2 `,
		groupID,
		slotID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity domain.Rotation
		err := rows.Scan(&entity.BannerID,
			&entity.ShowCount,
			&entity.ClickCount,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}

	return result, rows.Err()
}

func (s *Storage) GetRotation(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (domain.Rotation, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		INNER JOIN slot_banners sb ON sb.banner_id = r.banner_id AND sb.slot_id = r.slot_id 
		WHERE r.banner_id = $1 AND r.group_id = $2 AND r.slot_id = $3`,
		bannerID,
		groupID,
		slotID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return domain.Rotation{},
			customerrors.NotFound{
				Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
			}
	}
	if row.Err() != nil {
		return domain.Rotation{}, row.Err()
	}

	var entity domain.Rotation
	err := row.Scan(&entity.BannerID,
		&entity.ShowCount,
		&entity.ClickCount,
	)
	if err != nil {
		return domain.Rotation{}, err
	}

	return entity, nil
}

func (s *Storage) IncrementRotationClick(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	res, err := s.DB.ExecContext(ctx, `INSERT INTO rotations (banner_id, group_id, slot_id, click_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET click_count = rotations.click_count + 1;`,
		bannerID,
		groupID,
		slotID,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
		}
	}

	return nil
}

func (s *Storage) IncrementRotationShow(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	res, err := s.DB.ExecContext(ctx, `INSERT INTO rotations (banner_id, group_id, slot_id, show_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET show_count = rotations.show_count + 1;`,
		bannerID,
		groupID,
		slotID,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
		}
	}

	return nil
}
