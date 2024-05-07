package banner

type GetUserBannerInput struct {
	TagId            uint `query:"tag_id"`
	FeatureId        uint `query:"feature_id"`
	UseLatestVersion bool `query:"use_latest_version"`
}

type GetBannerInput struct {
	TagId     *uint `query:"tag_id"`
	FeatureId *uint `query:"feature_id"`
	Limit     *uint `query:"limit"`
	Offset    *uint `query:"offset"`
}
