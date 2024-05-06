package postgresql

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

type Postgresql struct {
	DB *gorm.DB
}

func NewConnection(config *config.Config) (Postgresql, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Address, config.Port, config.User, config.Password, config.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return Postgresql{}, err
	}
	return Postgresql{DB: db}, nil
}

func (r *Postgresql) Migrate() error {
	v := reflect.ValueOf(entity.Entities{})

	for i := 0; i < v.NumField(); i++ {
		err := r.DB.AutoMigrate(v.Field(i).Interface())
		if err != nil {
			return err
		}
	}

	return nil
}
