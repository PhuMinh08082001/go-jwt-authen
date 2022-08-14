package dal

import (
	"fmt"
	"github.com/PhuMinh08082001/go-jwt-authen/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Module = fx.Provide(NewDB)

func NewDB(config config.Configuration) (db *gorm.DB) {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Database.Username, config.Database.Password, config.Database.Name, config.Database.Port, config.Database.Sslmode)
	fmt.Println(dsn)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}
