package repos

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/utils/database"
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

func (b *BannerRepo) GetBanner(ctx context.Context, tagId, featureId, limit, offset *uint) ([]entity.BannerWithTags, error) {
	tableName, _ := database.GetTableName(b.postgres, entity.Banner{})
	q := b.postgres.DB.Table(tableName).
		Select("banners.id, banners.title, banners.text, banners.url, banners.feature_id, banners.is_active").
		Joins("LEFT JOIN banner_tags ON banners.id = banner_tags.banner_id")
	if tagId != nil && featureId != nil {
		q = q.Where("banner_tags.tag_id = ? AND banners.feature_id = ?", *tagId, *featureId)
	} else if tagId != nil {
		q = q.Where("banner_tags.tag_id = ?", *tagId)
	} else if featureId != nil {
		q = q.Where("banners.feature_id = ?", *featureId)
	}

	q = q.Group("banners.id")

	if limit != nil {
		q = q.Limit(int(*limit))
	}
	if offset != nil {
		q = q.Offset(int(*offset))
	}

	// Сначала вытаскиваем все баннеры, потому получаем id найденных баннеров
	var banners []entity.Banner
	if err := q.Find(&banners).Error; err != nil {
		return nil, err
	}

	var bannerIDs []uint
	for _, banner := range banners {
		bannerIDs = append(bannerIDs, banner.ID)
	}

	// Создать tmp-структуру чтобы зафетчить только tagId?
	var bannerTags []entity.BannerTag
	q = b.postgres.DB.Table("banner_tags").Where("banner_id IN (?)", bannerIDs)
	err := q.Scan(&bannerTags).Error
	if err != nil {
		return []entity.BannerWithTags{}, err
	}

	tagMap := make(map[uint][]uint)
	for _, bt := range bannerTags {
		tagMap[bt.BannerId] = append(tagMap[bt.BannerId], bt.TagId)
	}

	var result []entity.BannerWithTags
	for _, banner := range banners {
		result = append(result, entity.BannerWithTags{
			Banner: banner,
			Tags:   tagMap[banner.ID],
		})
	}

	return result, nil
}

func NewBannerRepo(postgres postgresql.Postgresql) *BannerRepo {
	return &BannerRepo{
		postgres: postgres,
	}
}
