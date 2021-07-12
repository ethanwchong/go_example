package repository

import (
	"fmt"
	"log"
	"videoproject/entity"
	"videoproject/util"

	_ "database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	jinzhuGorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"gorm.io/driver/postgres"
	gormGorm "gorm.io/gorm"
)

type VideoRepository interface {
	Save(video entity.Video)
	Update(video entity.Video)
	Delete(video entity.Video)
	FindAll() []entity.Video
	CloseDB()
}

type database struct {
	connection1 *jinzhuGorm.DB
	connection2 *gormGorm.DB
}

func NewVideoRepository() VideoRepository {

	fmt.Printf("Environment variable USEDB: %s \n", GetConfig().UseDB)
	if GetConfig().UseDB == "postgre_cloud_proxy" {
		gormPostgresDbCloudProxyDB := GetGormPostgresDBCloudProxy()
		gormPostgresDbCloudProxyDB.AutoMigrate(&entity.Video{}, &entity.Person{})
		return &database{
			connection1: nil,
			connection2: gormPostgresDbCloudProxyDB,
		}
	} else if GetConfig().UseDB == "postgre_remote" {
		gormPostgresRemoteDB := GetGormPostgresRemoteDB()
		gormPostgresRemoteDB.AutoMigrate(&entity.Video{}, &entity.Person{})
		return &database{
			connection1: nil,
			connection2: gormPostgresRemoteDB,
		}
	} else {
		//by default use SQLite
		jinzhuSqliteDb, err := jinzhuGorm.Open("sqlite3", "test.db")
		fmt.Println("Created SQLite3 DB")
		if err != nil {
			panic("Failed to connect to SQLite3 DB")
		}
		jinzhuSqliteDb.AutoMigrate(&entity.Video{}, &entity.Person{})
		return &database{
			connection1: jinzhuSqliteDb,
			connection2: nil,
		}
	}
}

func GetConfig() util.Config {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration \n", err)
	}
	return config
}

func GetGormPostgresDBCloudProxy() *gormGorm.DB {
	var config util.Config = GetConfig()
	var dbURI string = fmt.Sprintf("host=%s:%s:%s user=%s dbname=%s password=%s sslmode=%s", config.GcpProject, config.DBRegion, config.DBInstance, config.DBUser, config.DBName, config.DBPassword, config.DBSslmode)
	fmt.Printf("DB_URI: %s \nDB_DRIVER_NAME: %s \n", dbURI, config.DBDriver)

	db, err := gormGorm.Open(postgres.New(postgres.Config{
		DriverName: config.DBDriver,
		DSN:        dbURI,
	}), &gormGorm.Config{})

	if err != nil {
		panic("Failed to connect to Postgres DB with gormGorm")
	}

	fmt.Println("Created Postgres DB Cloud Proxy")
	return db
}

func GetGormPostgresRemoteDB() *gormGorm.DB {
	var config util.Config = GetConfig()
	var dbURI string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSslmode, config.DBTimezone)
	fmt.Printf("DB_URI: %s \n", dbURI)

	db, err := gormGorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gormGorm.Config{})

	if err != nil {
		panic("Failed to connect to Postgre DB")
	}

	fmt.Println("Created Postgres DB")
	return db
}

func (db *database) CloseDB() {
	if db.connection1 != nil {
		err := db.connection1.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}
	if db.connection2 != nil {
		db.CloseDB()
		fmt.Println("Closing Postgres DB Connection")
	}
}
func (db *database) Save(video entity.Video) {
	if db.connection1 != nil {
		db.connection1.Save(&video)
	}
	if db.connection2 != nil {
		db.connection2.Save(&video)
	}
}

func (db *database) Update(video entity.Video) {
	if db.connection1 != nil {
		db.connection1.Save(&video)
	}
	if db.connection2 != nil {
		db.connection2.Save(&video)
	}
}
func (db *database) Delete(video entity.Video) {
	if db.connection1 != nil {
		db.connection1.Delete(&video)
	}
	if db.connection2 != nil {
		db.connection2.Delete(&video)
	}
}
func (db *database) FindAll() []entity.Video {
	var videos []entity.Video
	if db.connection1 != nil {
		db.connection1.Set("gorm:auto_preload", true).Find(&videos)
		return videos
	}
	if db.connection2 != nil {
		db.connection2.Set("gorm:auto_preload", true).Find(&videos)
		return videos
	}
	return nil
}
