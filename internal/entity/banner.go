package entity

type Banner struct {
	ID        uint
	Title     string
	Text      string
	Url       string
	FeatureId uint
}

type BannerWithTag struct {
	Banner
	Tag int
}
