package repository

import (
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/gorm"
)

// BridgeRepository 定义了中转配置的接口
type BridgeRepository interface {
	Create(bridge *models.Bridge) error
	GetAll() ([]models.Bridge, error)
	GetByID(id string) (*models.Bridge, error)
	Update(bridge *models.Bridge) error
	Delete(id string) error
}

type bridgeRepository struct {
	db *gorm.DB
}

// NewBridgeRepository 创建一个新的中转配置仓库实例
func NewBridgeRepository(db *gorm.DB) BridgeRepository {
	return &bridgeRepository{db: db}
}

// Create 创建一个新的中转配置记录
func (r *bridgeRepository) Create(bridge *models.Bridge) error {
	return r.db.Create(bridge).Error
}

// GetAll 获取所有中转配置记录
func (r *bridgeRepository) GetAll() ([]models.Bridge, error) {
	var bridges []models.Bridge
	// 预加载关联的 SourceChannel 和 TargetChannel
	err := r.db.Preload("SourceChannel").Preload("TargetChannel").Find(&bridges).Error
	return bridges, err
}

// GetByID 根据ID获取一个中转配置记录
func (r *bridgeRepository) GetByID(id string) (*models.Bridge, error) {
	var bridge models.Bridge
	err := r.db.Preload("SourceChannel").Preload("TargetChannel").First(&bridge, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bridge, nil
}

// Update 更新一个已存在的中转配置记录
func (r *bridgeRepository) Update(bridge *models.Bridge) error {
	return r.db.Save(bridge).Error
}

// Delete 根据ID删除一个中转配置记录
func (r *bridgeRepository) Delete(id string) error {
	return r.db.Delete(&models.Bridge{}, "id = ?", id).Error
}
