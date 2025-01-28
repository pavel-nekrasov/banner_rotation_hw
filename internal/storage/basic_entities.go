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

func (s *Storage) CreateSlot(ctx context.Context, entity domain.Slot) error {
	res, err := s.DB.Exec(ctx, `INSERT INTO slots 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetSlot(ctx context.Context, entityID domain.SlotID) (domain.Slot, error) {
	var entity domain.Slot
	var description sql.NullString

	err := s.DB.QueryRow(ctx,
		`SELECT id, description 
		FROM slots WHERE id = $1`,
		entityID,
	).Scan(&entity.ID,
		&description)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Slot{}, customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entityID)}
	}
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
	rows, err := s.DB.Query(ctx,
		`SELECT id, description 
		FROM slots `,
	)
	if errors.Is(err, pgx.ErrNoRows) {
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
	res, err := s.DB.Exec(ctx, `UPDATE slots 
		SET description = $2 
		WHERE id = $1`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) DeleteSlot(ctx context.Context, entityID domain.SlotID) error {
	res, err := s.DB.Exec(ctx, "delete from slots where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Slot with id = \"%v\" not found", entityID)}
	}

	return nil
}

func (s *Storage) CreateBanner(ctx context.Context, entity domain.Banner) error {
	res, err := s.DB.Exec(ctx, `INSERT INTO banners 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetBanner(ctx context.Context, entityID domain.BannerID) (domain.Banner, error) {
	var entity domain.Banner
	var description sql.NullString

	err := s.DB.QueryRow(ctx,
		`SELECT id, description 
		FROM banners WHERE id = $1`,
		entityID,
	).Scan(
		&entity.ID,
		&description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Banner{}, customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entityID)}
	}
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
	rows, err := s.DB.Query(ctx,
		`SELECT id, description 
		FROM banners `,
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

func (s *Storage) UpdateBanner(ctx context.Context, entity domain.Banner) error {
	res, err := s.DB.Exec(ctx, `UPDATE banners 
		SET description = $2 
		WHERE id = $1`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) DeleteBanner(ctx context.Context, entityID domain.BannerID) error {
	res, err := s.DB.Exec(ctx, "delete from banners where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Banner with id = \"%v\" not found", entityID)}
	}

	return nil
}

func (s *Storage) CreateGroup(ctx context.Context, entity domain.Group) error {
	res, err := s.DB.Exec(ctx, `INSERT INTO groups 
		(id, description) 
		VALUES ($1, $2)`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) GetGroup(ctx context.Context, entityID domain.GroupID) (domain.Group, error) {
	var entity domain.Group
	var description sql.NullString

	err := s.DB.QueryRow(ctx,
		`SELECT id, description 
		FROM groups WHERE id = $1`,
		entityID,
	).Scan(
		&entity.ID,
		&description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Group{}, customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entityID)}
	}
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
	rows, err := s.DB.Query(ctx,
		`SELECT id, description 
		FROM groups `,
	)
	if errors.Is(err, pgx.ErrNoRows) {
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
	res, err := s.DB.Exec(ctx, `UPDATE groups 
		SET description = $2 
		WHERE id = $1`,
		entity.ID,
		entity.Description,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entity.ID)}
	}

	return nil
}

func (s *Storage) DeleteGroup(ctx context.Context, entityID domain.GroupID) error {
	res, err := s.DB.Exec(ctx, "delete from groups where id = $1", entityID)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Group with id = \"%v\" not found", entityID)}
	}

	return nil
}
