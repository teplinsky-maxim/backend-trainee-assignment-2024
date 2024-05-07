package repos

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"errors"
	"strconv"
)

var (
	ErrBannerNotFound      = errors.New("banner not found")
	BannerScanError        = errors.New("scanning to banner went wrong")
	BannerIsNotActiveError = errors.New("banner is not active")
)

type BannerRepo struct {
	postgres postgresql.Postgresql
}

func (b *BannerRepo) GetUserBanner(ctx context.Context, tagId uint, featureId uint, useLatestVersion bool) (entity.BannerWithTag, error) {
	query := `
SELECT b.id, b.title, b.text, b.url, b.feature_id, bt.tag_id as tag, b.is_active
FROM banners b
         JOIN banner_tags bt ON b.id = bt.banner_id
WHERE bt.tag_id = $1 AND b.feature_id = $2;
`
	var result entity.BannerWithTag
	row := b.postgres.DB.Raw(query, tagId, featureId)
	row.Scan(&result)
	// TODO: переделать на нормальную проверку понимания того, нашли мы результат или нет
	if result.ID == 0 {
		return entity.BannerWithTag{}, ErrBannerNotFound
	}
	if result.IsActive == false {
		userRole := ctx.Value(auth.RoleCtxField).(auth.Role)
		if userRole == auth.ADMIN {
			return result, nil
		} else if userRole == auth.USER {
			return entity.BannerWithTag{}, BannerIsNotActiveError
		} else {
			panic("Unhandled user role " + strconv.Itoa(int(userRole)))
		}
	}
	return result, nil
}

func NewBannerRepo(postgres postgresql.Postgresql) *BannerRepo {
	return &BannerRepo{
		postgres: postgres,
	}
}
