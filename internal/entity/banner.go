package entity

type Banner struct {
	ID        uint
	Title     string
	Text      string
	Url       string
	FeatureId uint
	IsActive  bool
}

func (r *Banner) ConvertToProductionBanner() ProductionBanner {
	return ProductionBanner{
		Title: r.Title,
		Text:  r.Text,
		Url:   r.Url,
	}
}

type ProductionBanner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
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
