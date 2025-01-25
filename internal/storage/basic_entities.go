package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (s *Storage) CreateSlot(ctx context.Context, entity domain.Slot) error {
	res, err := s.DB.ExecContext(ctx, `INSERT INTO slots 
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

func (s *Storage) GetSlot(ctx context.Context, entityID domain.SlotID) (domain.Slot, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT id, description 
		FROM slots WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return domain.Slot{}, customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return domain.Slot{}, row.Err()
	}
	var entity domain.Slot
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return domain.Slot{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListSlots(ctx context.Context) ([]domain.Slot, error) {
	result := make([]domain.Slot, 0)
	rows, err := s.DB.QueryContext(ctx,
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
		var entity domain.Slot
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

func (s *Storage) UpdateSlot(ctx context.Context, entity domain.Slot) error {
	res, err := s.DB.ExecContext(ctx, `UPDATE slots 
		SET description = $2 
		WHERE id = $1`,
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

func (s *Storage) DeleteSlot(ctx context.Context, entityID domain.SlotID) error {
	res, err := s.DB.ExecContext(ctx, "delete from slots where id = $1", entityID)
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

func (s *Storage) CreateBanner(ctx context.Context, entity domain.Banner) error {
	res, err := s.DB.ExecContext(ctx, `INSERT INTO banners 
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

func (s *Storage) GetBanner(ctx context.Context, entityID domain.BannerID) (domain.Banner, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT id, description 
		FROM banners WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return domain.Banner{}, customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return domain.Banner{}, row.Err()
	}
	var entity domain.Banner
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return domain.Banner{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListBanners(ctx context.Context) ([]domain.Banner, error) {
	result := make([]domain.Banner, 0)
	rows, err := s.DB.QueryContext(ctx,
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

func (s *Storage) UpdateBanner(ctx context.Context, entity domain.Banner) error {
	res, err := s.DB.ExecContext(ctx, `UPDATE banners 
		SET description = $2 
		WHERE id = $1`,
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

func (s *Storage) DeleteBanner(ctx context.Context, entityID domain.BannerID) error {
	res, err := s.DB.ExecContext(ctx, "delete from banners where id = $1", entityID)
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

func (s *Storage) CreateGroup(ctx context.Context, entity domain.Group) error {
	res, err := s.DB.ExecContext(ctx, `INSERT INTO groups 
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

func (s *Storage) GetGroup(ctx context.Context, entityID domain.GroupID) (domain.Group, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT id, description 
		FROM groups WHERE id = $1`,
		entityID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return domain.Group{}, customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entityID)}
	}
	if row.Err() != nil {
		return domain.Group{}, row.Err()
	}
	var entity domain.Group
	var description sql.NullString
	err := row.Scan(&entity.ID,
		&description,
	)
	if err != nil {
		return domain.Group{}, err
	}

	if description.Valid {
		entity.Description = description.String
	}

	return entity, nil
}

func (s *Storage) ListGroups(ctx context.Context) ([]domain.Group, error) {
	result := make([]domain.Group, 0)
	rows, err := s.DB.QueryContext(ctx,
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
		var entity domain.Group
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

func (s *Storage) UpdateGroup(ctx context.Context, entity domain.Group) error {
	res, err := s.DB.ExecContext(ctx, `UPDATE groups 
		SET description = $2 
		WHERE id = $1`,
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

func (s *Storage) DeleteGroup(ctx context.Context, entityID domain.GroupID) error {
	res, err := s.DB.ExecContext(ctx, "delete from groups where id = $1", entityID)
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
