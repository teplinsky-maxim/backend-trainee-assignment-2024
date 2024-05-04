package service

type Services struct {
	Banner Banner
}
type ServiceDependencies struct {
}

type Banner interface {
	GetUserBanner()
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Banner: NewBannerService(),
	}
}
