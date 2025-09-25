package utils

import (
	"log"

	"gorm.io/gorm"
)

// MigrationUtils 数据库迁移工具
type MigrationUtils struct {
	db *gorm.DB
}

// NewMigrationUtils 创建迁移工具
func NewMigrationUtils(db *gorm.DB) *MigrationUtils {
	return &MigrationUtils{db: db}
}

// RunMigrations 运行所有迁移
func (m *MigrationUtils) RunMigrations(models ...interface{}) error {
	log.Println("Starting database migrations...")

	// 迁移所有模型
	err := m.db.AutoMigrate(models...)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateIndexes 创建索引
func (m *MigrationUtils) CreateIndexes() error {
	log.Println("Creating database indexes...")

	// 为 Thing 创建索引
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_things_type ON things(type)").Error; err != nil {
		return err
	}

	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_things_created_at ON things(created_at)").Error; err != nil {
		return err
	}

	// 为 Relationship 创建索引
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_relationships_source_id ON relationships(source_id)").Error; err != nil {
		return err
	}

	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_relationships_target_id ON relationships(target_id)").Error; err != nil {
		return err
	}

	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_relationships_type ON relationships(type)").Error; err != nil {
		return err
	}

	// 为 ThingType 创建索引
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_thing_types_category ON thing_types(category)").Error; err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}

// SeedData 种子数据
func (m *MigrationUtils) SeedData() error {
	log.Println("Seeding initial data...")

	// 这里可以添加一些初始数据
	// 例如：默认的关系类型、系统配置等

	log.Println("Initial data seeded successfully")
	return nil
}
