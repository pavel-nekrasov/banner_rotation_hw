package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

func (s *Storage) ListRotations(ctx context.Context, groupId models.GroupId, slotId models.SlotId) ([]models.Rotation, error) {
	result := make([]models.Rotation, 0)
	rows, err := s.Db.QueryContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		WHERE r.group_id = $1 AND r.slot_id = $2`,
		groupId,
		slotId,
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

func (s *Storage) GetRotation(ctx context.Context, bannerID models.BannerId, groupID models.GroupId, slotID models.SlotId) (models.Rotation, error) {
	row, err := s.Db.QueryContext(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM rotations r 
		WHERE r.banner_id = $1 AND r.group_id = $2 AND r.slot_id = $3`,
		bannerID,
		groupID,
		slotID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Rotation{}, customerrors.NotFound{Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID)}
	}
	if err != nil {
		return models.Rotation{}, err
	}

	defer row.Close()

	var entity models.Rotation
	err = row.Scan(&entity.BannerID,
		&entity.ShowCount,
		&entity.ClickCount,
	)
	if err != nil {
		return models.Rotation{}, err
	}

	return entity, nil
}

func (s *Storage) IncrementRotationClick(ctx context.Context, bannerID models.BannerId, groupID models.GroupId, slotID models.SlotId) error {
	res, err := s.Db.ExecContext(ctx, `INSERT INTO rotations (banner_id, group_id, slot_id, click_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET click_count = click_count + 1;`,
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
		return customerrors.NotFound{Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID)}
	}

	return nil
}

func (s *Storage) IncrementRotationShow(ctx context.Context, bannerID models.BannerId, groupID models.GroupId, slotID models.SlotId) error {
	res, err := s.Db.ExecContext(ctx, `INSERT INTO rotations (banner_id, group_id, slot_id, show_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET show_count = show_count + 1;`,
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
		return customerrors.NotFound{Message: fmt.Sprintf("Rotation entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID)}
	}

	return nil
}
