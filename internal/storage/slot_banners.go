package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

func (s *Storage) AddBannerToSlot(ctx context.Context, slotId models.SlotId, bannerId models.BannerId) error {
	res, err := s.Db.ExecContext(ctx, `INSERT INTO slot_banners 
		(slot_id, banner_id) 
		VALUES ($1, $2)`,
		slotId,
		bannerId,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Failed to add banner Id = \"%v\" to slot id = \"%v\"", bannerId, slotId)}
	}

	return nil
}

func (s *Storage) RemoveBannerFromSlot(ctx context.Context, slotId models.SlotId, bannerId models.BannerId) error {
	res, err := s.Db.ExecContext(ctx, "DELETE FROM slot_banners WHERE slot_id = $1 AND banner_id = $2", slotId, bannerId)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner id = \"%v\" not found in slot id = \"%v\"", bannerId, slotId)}
	}

	return nil
}

func (s *Storage) ListSlotBanners(ctx context.Context, slotId models.SlotId) ([]models.Banner, error) {
	result := make([]models.Banner, 0)
	rows, err := s.Db.QueryContext(ctx,
		`SELECT b.id, b.description 
		FROM banners b 
		INNER JOIN slot_banners sb ON sb.banner_id = b.id 
		WHERE sb.slot_id = $1`,
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
		var entity models.Banner
		var description sql.NullString
		err := rows.Scan(&entity.ID,
			&description,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			entity.Description = description.String
		}

		result = append(result, entity)
	}

	return result, rows.Err()
}
