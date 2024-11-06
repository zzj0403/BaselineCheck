package repository

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	DBHost          string        `mapstructure:"db_host" json:"db_host"`
	DBPort          string        `mapstructure:"db_port" json:"db_port"`
	DBName          string        `mapstructure:"db_name" json:"db_name"`
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_life_time" json:"conn_max_life_time"` // second
}

func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		DBHost:          "127.0.0.1",
		DBPort:          "3306",
		DBName:          "awesome_template",
		Username:        "root",
		Password:        "123456",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 60 * time.Second,
	}
}

func (cfg *DBConfig) ToDSN() string {
	params := url.Values{}
	params.Add("charset", "utf8mb4")
	params.Add("collation", "utf8mb4_unicode_ci")
	params.Add("parseTime", "true")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.Username, cfg.Password, cfg.DBHost, cfg.DBPort, cfg.DBName, params.Encode())
}

type MyDB struct {
	*gorm.DB
}

func NewDb(cfg *DBConfig) (MyDB, error) {
	if cfg == nil {
		cfg = DefaultDBConfig()
	}

	connectionString := cfg.ToDSN()
	// log.Printf("DSN: %s", connectionString)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connectionString,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return MyDB{}, fmt.Errorf("数据库连接错误: %w", err)
	}

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return MyDB{}, fmt.Errorf("获取数据库实例失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("成功连接到数据库:", cfg.DBName)
	return MyDB{db}, nil
}
