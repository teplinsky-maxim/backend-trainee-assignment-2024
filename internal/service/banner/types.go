package banner

type GetUserBannerInput struct {
	TagId            uint `query:"tag_id"`
	FeatureId        uint `query:"feature_id"`
	UseLatestVersion bool `query:"use_latest_version"`
}
