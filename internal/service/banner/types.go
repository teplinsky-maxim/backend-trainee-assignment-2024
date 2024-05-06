package banner

type GetUserBannerInput struct {
	TagId            int  `query:"tag_id"`
	FeatureId        int  `query:"feature_id"`
	UseLatestVersion bool `query:"use_latest_version"`
}
