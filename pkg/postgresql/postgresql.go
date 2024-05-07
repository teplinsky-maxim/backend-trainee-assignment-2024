package postgresql

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

const TestsSchemaName = "tests"

type Postgresql struct {
	DB *gorm.DB
}

type PostgresqlTest Postgresql

func newConnection(config *config.Config, testConnection bool) (Postgresql, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Address, config.Port, config.User, config.Password, config.Database,
	)

	var db *gorm.DB
	var err error
	if !testConnection {
		db, err = gorm.Open(postgres.Open(dsn))
	} else {
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   TestsSchemaName + ".",
				SingularTable: false,
			},
		})
	}

	if err != nil {
		return Postgresql{}, err
	}
	return Postgresql{DB: db}, nil
}

func NewConnection(config *config.Config) (Postgresql, error) {
	return newConnection(config, false)
}

func NewTestConnection(config *config.Config) (PostgresqlTest, error) {
	conn, err := newConnection(config, true)
	conn.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + TestsSchemaName)
	conn.DB.Exec("SET search_path to " + TestsSchemaName)
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}
	return PostgresqlTest{
		DB: conn.DB,
	}, err
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

func (r *PostgresqlTest) SetUp(tableName string) {
	r.DB.Exec("TRUNCATE " + tableName)
}

func (r *PostgresqlTest) TearDown(tableName string) {
	r.DB.Exec("TRUNCATE " + tableName)
}
