package usecase

import (
	"errors"
	"strconv"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

// BridgeUsecase 定义了中转配置用例的接口
type BridgeUsecase interface {
	CreateBridge(bridge *models.Bridge) error
	GetAllBridges() ([]models.Bridge, error)
	GetBridgeByID(id int) (*models.Bridge, error)
	GetBridgeByTicket(ticket string) (*models.Bridge, error)
	UpdateBridge(id int, bridgeUpdates *models.Bridge) (*models.Bridge, error)
	DeleteBridge(id int) error
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
func (u *bridgeUsecase) validateChannelExists(channelID int) error {
	if channelID == 0 {
		return errors.New("渠道ID不能为空")
	}
	_, err := u.channelRepo.GetByID(channelID)
	if err != nil {
		return errors.New("渠道ID " + strconv.FormatUint(uint64(channelID), 10) + " 不存在或无效")
	}
	return nil
}

// CreateBridge 创建一个新的中转配置
func (u *bridgeUsecase) CreateBridge(bridge *models.Bridge) error {
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
func (u *bridgeUsecase) GetBridgeByID(id int) (*models.Bridge, error) {
	return u.bridgeRepo.GetByID(id)
}

// GetBridgeByTicket 根据Ticket获取一个中转配置
func (u *bridgeUsecase) GetBridgeByTicket(ticket string) (*models.Bridge, error) {
	return u.bridgeRepo.GetByTicket(ticket)
}

// UpdateBridge 更新一个中转配置
func (u *bridgeUsecase) UpdateBridge(id int, bridgeUpdates *models.Bridge) (*models.Bridge, error) {
	existingBridge, err := u.bridgeRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("中转配置未找到")
	}

	// 更新字段，仅当在请求中提供了这些字段时
	if bridgeUpdates.Name != "" {
		existingBridge.Name = bridgeUpdates.Name
	}

	if err := u.validateChannelExists(bridgeUpdates.TargetChannelID); err != nil {
		return nil, errors.New("新的目标渠道验证失败: " + err.Error())
	}

	existingBridge.IsActive = bridgeUpdates.IsActive

	if err := u.bridgeRepo.Update(existingBridge); err != nil {
		return nil, err
	}
	return existingBridge, nil
}

// DeleteBridge 根据ID删除一个中转配置
func (u *bridgeUsecase) DeleteBridge(id int) error {
	// 可以在删除前添加检查逻辑，例如检查中转配置是否仍在使用等
	_, err := u.bridgeRepo.GetByID(id) // 确认存在
	if err != nil {
		return errors.New("中转配置未找到，无法删除")
	}
	return u.bridgeRepo.Delete(id)
}
