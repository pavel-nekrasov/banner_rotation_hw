package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/customerrors"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/domain"
)

func (s *Storage) ListBannerStats(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
) ([]domain.BannerStat, error) {
	result := make([]domain.BannerStat, 0)
	rows, err := s.DB.Query(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM banner_stats r 
		INNER JOIN slot_banners sb ON sb.banner_id = r.banner_id AND sb.slot_id = r.slot_id 
		WHERE r.group_id = $1 AND r.slot_id = $2 `,
		groupID,
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
		var entity domain.BannerStat
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

func (s *Storage) GetBannerStat(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) (domain.BannerStat, error) {
	var entity domain.BannerStat
	err := s.DB.QueryRow(ctx,
		`SELECT r.banner_id, r.show_count, r.click_count 
		FROM banner_stats r 
		INNER JOIN slot_banners sb ON sb.banner_id = r.banner_id AND sb.slot_id = r.slot_id 
		WHERE r.banner_id = $1 AND r.group_id = $2 AND r.slot_id = $3`,
		bannerID,
		groupID,
		slotID,
	).Scan(&entity.BannerID,
		&entity.ShowCount,
		&entity.ClickCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.BannerStat{},
			customerrors.NotFound{
				Message: fmt.Sprintf("Banner stat entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
			}
	}
	if err != nil {
		return domain.BannerStat{}, err
	}

	return entity, nil
}

func (s *Storage) BannerStatTotals(
	ctx context.Context,
) (int64, int64, error) {
	var showSum, clickSum, recCnt int64
	err := s.DB.QueryRow(ctx,
		`SELECT sum(r.show_count) AS show_sum, sum(r.click_count) AS click_sum, count(r.id) AS rec_count 
		FROM banner_stats r 
		INNER JOIN slot_banners sb ON sb.banner_id = r.banner_id AND sb.slot_id = r.slot_id`,
	).Scan(&showSum,
		&clickSum,
		&recCnt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, customerrors.NotFound{}
	}
	if err != nil {
		return 0, 0, err
	}

	return showSum - recCnt, clickSum - recCnt, nil
}

func (s *Storage) IncrementBannerStatClick(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	res, err := s.DB.Exec(ctx, `INSERT INTO banner_stats (banner_id, group_id, slot_id, click_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET click_count = banner_stats.click_count + 1;`,
		bannerID,
		groupID,
		slotID,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Banner stat entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
		}
	}

	return nil
}

func (s *Storage) IncrementBannerStatShow(
	ctx context.Context,
	groupID domain.GroupID,
	slotID domain.SlotID,
	bannerID domain.BannerID,
) error {
	res, err := s.DB.Exec(ctx, `INSERT INTO banner_stats (banner_id, group_id, slot_id, show_count)
    VALUES ($1, $2, $3, 2)
    ON CONFLICT (banner_id, group_id, slot_id) DO UPDATE SET show_count = banner_stats.show_count + 1;`,
		bannerID,
		groupID,
		slotID,
	)
	if err != nil {
		return err
	}

	cnt := res.RowsAffected()
	if cnt == 0 {
		return customerrors.NotFound{
			Message: fmt.Sprintf("Banner stat entry for banner/group/slot %v/%v/%v not found", bannerID, groupID, slotID),
		}
	}

	return nil
}
