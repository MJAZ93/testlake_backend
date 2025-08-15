package dao

import (
	"testlake/model"

	"github.com/google/uuid"
)

type PlanDao struct {
	Limit int
}

func NewPlanDao() *PlanDao {
	return &PlanDao{Limit: 50}
}

func (dao *PlanDao) Create(plan *model.Plan) error {
	return Database.Create(plan).Error
}

func (dao *PlanDao) GetByID(id uuid.UUID) (*model.Plan, error) {
	var plan model.Plan
	err := Database.First(&plan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (dao *PlanDao) GetBySlug(slug string) (*model.Plan, error) {
	var plan model.Plan
	err := Database.First(&plan, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (dao *PlanDao) GetAll() ([]model.Plan, error) {
	var plans []model.Plan
	err := Database.Where("is_active = ?", true).Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (dao *PlanDao) GetAllWithPagination(page int) ([]model.Plan, int64, error) {
	var plans []model.Plan
	var total int64

	err := Database.Model(&model.Plan{}).Where("is_active = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Where("is_active = ?", true).Offset(offset).Limit(dao.Limit).Find(&plans).Error
	if err != nil {
		return nil, 0, err
	}

	return plans, total, nil
}

func (dao *PlanDao) Update(plan *model.Plan) error {
	return Database.Save(plan).Error
}

func (dao *PlanDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.Plan{}, "id = ?", id).Error
}

func (dao *PlanDao) SetActive(id uuid.UUID, active bool) error {
	return Database.Model(&model.Plan{}).Where("id = ?", id).Update("is_active", active).Error
}

func (dao *PlanDao) UpdatePayPalPlanIDs(id uuid.UUID, monthlyPlanID, yearlyPlanID string) error {
	updates := map[string]interface{}{
		"paypal_monthly_plan_id": monthlyPlanID,
		"paypal_yearly_plan_id":  yearlyPlanID,
	}
	return Database.Model(&model.Plan{}).Where("id = ?", id).Updates(updates).Error
}
