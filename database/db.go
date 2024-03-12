package database

import (
	"fmt"

	zconstant "lazy-auth/app/constant"
	"lazy-auth/app/repository"
	"lazy-auth/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitDatabase(configEnv config.ConfigEnv) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v",
		configEnv.DataBaseHost,
		configEnv.DataBaseUser,
		configEnv.DataBasePassword,
		configEnv.DataBaseName,
		configEnv.DataBasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	if configEnv.DataBaseAutoMigrate {
		fmt.Println("[GORM] [WARNING] Automatically migrate your schema, this is NOT safe.")
		db.AutoMigrate(
			&repository.Role{},
			&repository.User{},
			&repository.Session{},
		)
	}

	// Initial role
	roles := zconstant.GetDefaultRoles()
	for _, role := range roles {
		prepareCreateRole := repository.Role{Name: role}
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&prepareCreateRole)
	}

	return db
}
