package usecase

import (
	"errors"

	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

// BridgeUsecase 定义了中转配置用例的接口
type BridgeUsecase interface {
	CreateBridge(bridge *models.Bridge) error
	GetAllBridges() ([]models.Bridge, error)
	GetBridgeByID(id string) (*models.Bridge, error)
	UpdateBridge(id string, bridgeUpdates *models.Bridge) (*models.Bridge, error)
	DeleteBridge(id string) error
	// TODO: 考虑添加ToggleActive(id string) (bool, error) 方法
}

type bridgeUsecase struct {
	bridgeRepo  repository.BridgeRepository
	channelRepo repository.ChannelRepository // 用于验证 ChannelID 是否存在
}

// NewBridgeUsecase 创建一个新的中转配置用例实例
func NewBridgeUsecase(bridgeRepo repository.BridgeRepository, channelRepo repository.ChannelRepository) BridgeUsecase {
	return &bridgeUsecase{
		bridgeRepo:  bridgeRepo,
		channelRepo: channelRepo,
	}
}

// validateChannelExists 验证Channel是否存在
func (u *bridgeUsecase) validateChannelExists(channelID string) error {
	if channelID == "" {
		return errors.New("渠道ID不能为空")
	}
	_, err := u.channelRepo.GetByID(channelID)
	if err != nil {
		return errors.New("渠道ID " + channelID + " 不存在或无效")
	}
	return nil
}

// CreateBridge 创建一个新的中转配置
func (u *bridgeUsecase) CreateBridge(bridge *models.Bridge) error {
	if bridge.SourceChannelID == bridge.TargetChannelID {
		return errors.New("源渠道和目标渠道不能相同")
	}
	if err := u.validateChannelExists(bridge.SourceChannelID); err != nil {
		return errors.New("源渠道验证失败: " + err.Error())
	}
	if err := u.validateChannelExists(bridge.TargetChannelID); err != nil {
		return errors.New("目标渠道验证失败: " + err.Error())
	}
	return u.bridgeRepo.Create(bridge)
}

// GetAllBridges 获取所有中转配置
func (u *bridgeUsecase) GetAllBridges() ([]models.Bridge, error) {
	return u.bridgeRepo.GetAll()
}

// GetBridgeByID 根据ID获取一个中转配置
func (u *bridgeUsecase) GetBridgeByID(id string) (*models.Bridge, error) {
	return u.bridgeRepo.GetByID(id)
}

// UpdateBridge 更新一个中转配置
func (u *bridgeUsecase) UpdateBridge(id string, bridgeUpdates *models.Bridge) (*models.Bridge, error) {
	existingBridge, err := u.bridgeRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("中转配置未找到")
	}

	// 更新字段，仅当在请求中提供了这些字段时
	if bridgeUpdates.Name != "" {
		existingBridge.Name = bridgeUpdates.Name
	}
	if bridgeUpdates.SourceChannelID != "" {
		if err := u.validateChannelExists(bridgeUpdates.SourceChannelID); err != nil {
			return nil, errors.New("新的源渠道验证失败: " + err.Error())
		}
		existingBridge.SourceChannelID = bridgeUpdates.SourceChannelID
	}
	if bridgeUpdates.TargetChannelID != "" {
		if err := u.validateChannelExists(bridgeUpdates.TargetChannelID); err != nil {
			return nil, errors.New("新的目标渠道验证失败: " + err.Error())
		}
		existingBridge.TargetChannelID = bridgeUpdates.TargetChannelID
	}

	if existingBridge.SourceChannelID == existingBridge.TargetChannelID {
		return nil, errors.New("更新后源渠道和目标渠道不能相同")
	}

	// IsActive 的处理: 通常 bool 类型会有一个明确的值，这里假设如果传递了 IsActive，就使用传递的值
	// 如果API设计允许部分更新且不包含 IsActive，则不更新它。
	// 为简单起见，这里直接赋值，如果需要更复杂的逻辑（例如，仅当显式提供时才更新），则需要调整。
	// 暂时我们让Update请求体中必须包含IsActive字段。
	existingBridge.IsActive = bridgeUpdates.IsActive

	if err := u.bridgeRepo.Update(existingBridge); err != nil {
		return nil, err
	}
	return existingBridge, nil
}

// DeleteBridge 根据ID删除一个中转配置
func (u *bridgeUsecase) DeleteBridge(id string) error {
	// 可以在删除前添加检查逻辑，例如检查中转配置是否仍在使用等
	_, err := u.bridgeRepo.GetByID(id) // 确认存在
	if err != nil {
		return errors.New("中转配置未找到，无法删除")
	}
	return u.bridgeRepo.Delete(id)
}
