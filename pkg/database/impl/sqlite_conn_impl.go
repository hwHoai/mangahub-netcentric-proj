package impl

import (
	"log"
	"mangahub/pkg/database"
	"mangahub/pkg/models"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteConnImpl struct {}

var _ database.DatabaseConnectionInterface = (*SqliteConnImpl)(nil)

func (s *SqliteConnImpl) InitDB(dbPath string) (*gorm.DB, error) {
	//
	if dbPath == "" {
		dbPath = "../../data/mangahub.db"
	}
	// 1. Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Show heavy queries (>1s)
			LogLevel:                  logger.Info, // Log all queries (Info level)
			IgnoreRecordNotFoundError: true, 		// Ignore "record not found" errors	
			Colorful:                  true,		// Colorful output
		},
	)

	// 2. Open conn
	// _pragma=foreign_keys(1): Enable foreign key constraints, ensure data integrity
	// _pragma=journal_mode(WAL): Use Write-Ahead Logging (WAL) for better concurrency and performance
	db, err := gorm.Open(sqlite.Open(dbPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// 3. Connection pool settings
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)                  // SQLite hoạt động tốt nhất với 1 writer
		sqlDB.SetMaxIdleConns(10)                 // Giữ kết nối chờ
		sqlDB.SetConnMaxLifetime(time.Hour * 1)   // Tránh rò rỉ kết nối
	}

	// 4. Setup Join Tables - MUST DO BEFORE AUTO MIGRATE
	err = db.SetupJoinTable(&models.MangaModel{}, "Followers", &models.MangaFollowerModel{})
	if err != nil {
		return nil, err
	}

	err = db.SetupJoinTable(&models.UserModel{}, "FollowingMangas", &models.MangaFollowerModel{})
	if err != nil {
		return nil, err
	}
    

	// 5. AUTO MIGRATE
	err = db.AutoMigrate(
		&models.UserModel{},
		&models.MangaModel{},
		&models.MangaFollowerModel{},
        &models.ReadingProgressModel{},
		&models.MessageModel{},
		&models.ReviewModel{},
		&models.SessionModel{},
		&models.WishlistModel{},
		// Thêm các model khác vào đây...
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}