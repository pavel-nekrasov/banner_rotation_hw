package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/models"
)

func (s *Storage) CreateSlot(ctx context.Context, entity models.Slot) error {
	res, err := s.db.ExecContext(ctx, `INSERT INTO slots 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetSlot(ctx context.Context, entityID string) (models.Slot, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, description 
		FROM slots WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return models.Slot{}, customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return models.Slot{}, row.Err()
	}
	var entity models.Slot
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return models.Slot{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListSlots(
	ctx context.Context,
) ([]models.Slot, error) {
	result := make([]models.Slot, 0)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, description 
		FROM slots `,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity models.Slot
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

func (s *Storage) UpdateSlot(ctx context.Context, event models.Slot) error {
	res, err := s.db.ExecContext(ctx, `UPDATE slots 
		SET description = $2 
		WHERE id = $1`,
		event.ID,
		event.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", event.ID)}
	}

	return nil
}

func (s *Storage) DeleteSlot(ctx context.Context, entityID string) error {
	res, err := s.db.ExecContext(ctx, "delete from slots where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entityID)}
	}

	return nil
}

func (s *Storage) CreateBanner(ctx context.Context, entity models.Banner) error {
	res, err := s.db.ExecContext(ctx, `INSERT INTO banners 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetBanner(ctx context.Context, entityID string) (models.Banner, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, description 
		FROM banners WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return models.Banner{}, customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return models.Banner{}, row.Err()
	}
	var entity models.Banner
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return models.Banner{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListBanners(
	ctx context.Context,
) ([]models.Banner, error) {
	result := make([]models.Banner, 0)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, description 
		FROM banners `,
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

func (s *Storage) UpdateBanner(ctx context.Context, event models.Banner) error {
	res, err := s.db.ExecContext(ctx, `UPDATE banners 
		SET description = $2 
		WHERE id = $1`,
		event.ID,
		event.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", event.ID)}
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, entityID string) error {
	res, err := s.db.ExecContext(ctx, "delete from banners where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entityID)}
	}

	return nil
}

func (s *Storage) CreateGroup(ctx context.Context, entity models.Group) error {
	res, err := s.db.ExecContext(ctx, `INSERT INTO groups 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetGroup(ctx context.Context, entityID string) (models.Group, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, description 
		FROM groups WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return models.Group{}, customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return models.Group{}, row.Err()
	}
	var entity models.Group
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return models.Group{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListGroups(
	ctx context.Context,
) ([]models.Group, error) {
	result := make([]models.Group, 0)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, description 
		FROM groups `,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity models.Group
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

func (s *Storage) UpdateGroup(ctx context.Context, event models.Group) error {
	res, err := s.db.ExecContext(ctx, `UPDATE groups 
		SET description = $2 
		WHERE id = $1`,
		event.ID,
		event.Description,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", event.ID)}
	}

	return nil
}

func (s *Storage) DeleteGroup(ctx context.Context, entityID string) error {
	res, err := s.db.ExecContext(ctx, "delete from groups where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entityID)}
	}

	return nil
}
