package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

func (s *Storage) ListRotations(
	ctx context.Context,
	groupID models.GroupID,
	slotID models.SlotID,
) ([]models.Rotation, error) {
	result := make([]models.Rotation, 0)
	rows, err := s.DB.QueryContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		WHERE r.group_id = $1 AND r.slot_id = $2`,
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
		var entity models.Rotation
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
	groupID models.GroupID,
	slotID models.SlotID,
	bannerID models.BannerID,
) (models.Rotation, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		WHERE r.banner_id = $1 AND r.group_id = $2 AND r.slot_id = $3`,
		bannerID,
		groupID,
		slotID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return models.Rotation{},
			customerrors.NotFound{
				Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
			}
	}
	if row.Err() != nil {
		return models.Rotation{}, row.Err()
	}

	var entity models.Rotation
	err := row.Scan(&entity.BannerID,
		&entity.ShowCount,
		&entity.ClickCount,
	)
	if err != nil {
		return models.Rotation{}, err
	}

	return entity, nil
}

func (s *Storage) IncrementRotationClick(
	ctx context.Context,
	groupID models.GroupID,
	slotID models.SlotID,
	bannerID models.BannerID,
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
	groupID models.GroupID,
	slotID models.SlotID,
	bannerID models.BannerID,
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
