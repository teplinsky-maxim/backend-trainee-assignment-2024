package cache

import "avito-backend-2024-trainee/internal/entity"

type BannerCache interface {
	Get(tagId, featureId uint) (entity.ProductionBanner, error)
	Set(featureId, tagId uint, banner entity.ProductionBanner)
}
