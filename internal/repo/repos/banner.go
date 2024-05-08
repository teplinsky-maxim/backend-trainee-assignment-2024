package repos

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/utils/database"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
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

	tagMap, err := fetchBannersIds(bannerIDs, b.postgres.DB)
	if err != nil {
		return []entity.BannerWithTags{}, err
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

func (b *BannerRepo) CreateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool) (entity.BannerId, error) {
	tx := b.postgres.DB.Begin()

	banner := entity.Banner{
		Title:     title,
		Text:      text,
		Url:       url,
		FeatureId: featureId,
		IsActive:  isActive,
	}
	if err := tx.Create(&banner).Error; err != nil {
		tx.Rollback()
		return entity.BannerId{}, err
	}

	bannerTags := make([]entity.BannerTag, len(tagIds))
	for idx, tagId := range tagIds {
		bt := entity.BannerTag{
			BannerId: banner.ID,
			TagId:    tagId,
		}
		bannerTags[idx] = bt
	}
	if err := tx.Create(&bannerTags).Error; err != nil {
		tx.Rollback()
		return entity.BannerId{}, err
	}
	if err := tx.Commit().Error; err != nil {
		return entity.BannerId{}, err
	}
	return entity.BannerId{ID: banner.ID}, nil
}

func (b *BannerRepo) UpdateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool, bannerId uint) error {
	tx := b.postgres.DB.Begin()
	var banner entity.Banner
	err := tx.Model(&banner).Clauses(clause.Returning{}).Where("id = ?", bannerId).Updates(entity.Banner{
		Title:     title,
		Text:      text,
		Url:       url,
		FeatureId: featureId,
		IsActive:  isActive,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if banner.ID != bannerId {
		return ErrBannerNotFound
	}

	tagMap, err := fetchBannersIds([]uint{bannerId}, b.postgres.DB)
	if err != nil {
		return err
	}

	tagIdsMap := lo.SliceToMap(tagIds, func(item uint) (uint, bool) {
		return item, true
	})

	currentBannerTags := tagMap[bannerId]
	currentBannerTagsMap := lo.SliceToMap(currentBannerTags, func(item uint) (uint, bool) {
		return item, true
	})

	tagsToDelete := lo.Filter(currentBannerTags, func(item uint, _ int) bool {
		_, exists := tagIdsMap[item]
		return !exists
	})

	tagsToCreate := lo.Filter(tagIds, func(item uint, _ int) bool {
		_, exists := currentBannerTagsMap[item]
		return !exists
	})

	minLen := min(len(tagsToCreate), len(tagsToDelete))
	maxLen := max(len(tagsToCreate), len(tagsToDelete))
	oldToNew := make(map[uint]uint, minLen)
	for i := 0; i < minLen; i++ {
		oldToNew[tagsToDelete[i]] = tagsToCreate[i]
	}

	updateStmt := generateUpdateCaseStatementsWithQuery(oldToNew, "banner_tags", "tag_id")
	err = tx.Exec(updateStmt).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if maxLen-minLen == 0 {
		return nil
	}

	if len(tagsToDelete) > len(tagsToCreate) {
		tagsLeft := tagsToDelete[maxLen-minLen:]
		err = tx.Model(entity.BannerTag{}).Where("tag_id IN (?)", tagsLeft).Delete(&entity.BannerTag{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tagsLeft := tagsToCreate[maxLen-minLen:]
		bannerTags := make([]entity.BannerTag, maxLen-minLen)
		for idx, tagId := range tagsLeft {
			bannerTags[idx] = entity.BannerTag{
				BannerId: bannerId,
				TagId:    tagId,
			}
		}
		// Здесь будет создана транзакция еще одна (на batch insert), но она тоже должна откатиться
		// в случае отката первоначальной транзакции
		err = tx.Create(bannerTags).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func generateUpdateCaseStatementsWithQuery(oldToNew map[uint]uint, bannerTableName, bannerIdColName string) string {
	q := fmt.Sprintf(`
UPDATE %s SET %s = CASE
`, bannerTableName, bannerIdColName)
	cases := make([]string, len(oldToNew))
	i := 0
	for oldId, newId := range oldToNew {
		cases[i] = fmt.Sprintf("WHEN %s = %d THEN %d", bannerIdColName, oldId, newId)
		i++
	}
	q = q + strings.Join(cases, " ")
	keys := convertSliceToSQLSyntax(lo.Keys(oldToNew))
	q = q + fmt.Sprintf(" END WHERE %s.%s IN %v", bannerTableName, bannerIdColName, keys)
	return q
}

func convertSliceToSQLSyntax(slice []uint) string {
	var values []string
	for _, val := range slice {
		values = append(values, fmt.Sprintf("%d", val))
	}
	return fmt.Sprintf("(%s)", strings.Join(values, ", "))
}

func fetchBannersIds(bannersIds []uint, db *gorm.DB) (map[uint][]uint, error) {
	var bannerTags []entity.BannerTag
	q := db.Table("banner_tags").Where("banner_id IN (?)", bannersIds)
	err := q.Scan(&bannerTags).Error
	if err != nil {
		return map[uint][]uint{}, err
	}

	tagMap := make(map[uint][]uint, len(bannersIds))
	for _, bannerId := range bannersIds {
		tagMap[bannerId] = []uint{}
	}
	for _, bt := range bannerTags {
		tagMap[bt.BannerId] = append(tagMap[bt.BannerId], bt.TagId)
	}
	return tagMap, nil
}

func NewBannerRepo(postgres postgresql.Postgresql) *BannerRepo {
	return &BannerRepo{
		postgres: postgres,
	}
}
