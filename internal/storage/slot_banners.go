package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (s *Storage) AddBannerToSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	res, err := s.conn.DB.Exec(ctx, `INSERT INTO slot_banners 
		(slot_id, banner_id) 
		VALUES ($1, $2)`,
		slotID,
		bannerID,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Failed to add banner Id = \"%v\" to slot id = \"%v\"", bannerID, slotID),
		}
	}

	return nil
}

func (s *Storage) RemoveBannerFromSlot(ctx context.Context, slotID domain.SlotID, bannerID domain.BannerID) error {
	res, err := s.conn.DB.Exec(ctx, "DELETE FROM slot_banners WHERE slot_id = $1 AND banner_id = $2", slotID, bannerID)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Banner id = \"%v\" not found in slot id = \"%v\"", bannerID, slotID),
		}
	}

	return nil
}

func (s *Storage) GetSlotBanner(
	ctx context.Context,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (domain.Banner, error) {
	var entity domain.Banner
	var description sql.NullString

	err := s.conn.DB.QueryRow(ctx,
		`SELECT b.id, b.description 
		FROM banners b 
		INNER JOIN slot_banners sb ON sb.banner_id = b.id 
		WHERE sb.slot_id = $1 AND sb.banner_id = $2`,
		slotID,
		bannerID,
	).Scan(
		&entity.ID,
		&description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Banner{}, customerrors.NotFound{
			Message: fmt.Sprintf("Slot id = \"%v\" does not containt banner id = %v", slotID, bannerID),
		}
	}
	if err != nil {
		return domain.Banner{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListSlotBanners(ctx context.Context, slotID domain.SlotID) ([]domain.Banner, error) {
	result := make([]domain.Banner, 0)
	rows, err := s.conn.DB.Query(ctx,
		`SELECT b.id, b.description 
		FROM banners b 
		INNER JOIN slot_banners sb ON sb.banner_id = b.id 
		WHERE sb.slot_id = $1`,
		slotID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity domain.Banner
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
