package banner

type GetUserBannerInput struct {
	TagId           uint `query:"tag_id"`
	FeatureId       uint `query:"feature_id"`
	UseLastRevision bool `query:"use_last_revision"`
}

type GetBannerInput struct {
	TagId     *uint `query:"tag_id"`
	FeatureId *uint `query:"feature_id"`
	Limit     *uint `query:"limit"`
	Offset    *uint `query:"offset"`
}

type CreateBannerInput struct {
	TagIds    []uint `json:"tag_ids"`
	FeatureId uint   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		Url   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_active"`
}

type UpdateBannerInput struct {
	CreateBannerInput
}

type DeleteBannerInput struct {
}
