package entity

type Banner struct {
	ID        uint
	Title     string
	Text      string
	Url       string
	FeatureId uint
	IsActive  bool
}

type BannerWithTag struct {
	Banner
	Tag uint
}

type BannerWithTags struct {
	Banner
	Tags []uint
}

type BannerId struct {
	ID uint
}
