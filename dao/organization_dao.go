package dao

import (
	"testlake/model"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type OrganizationDao struct {
	Limit int
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{Limit: 50}
}

func (dao *OrganizationDao) Create(org *model.Organization) error {
	return Database.Create(org).Error
}

func (dao *OrganizationDao) GetByID(id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	err := Database.First(&org, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (dao *OrganizationDao) GetBySlug(slug string) (*model.Organization, error) {
	var org model.Organization
	err := Database.First(&org, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (dao *OrganizationDao) GetAll(page int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	err := Database.Model(&model.Organization{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Offset(offset).Limit(dao.Limit).Find(&orgs).Error
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (dao *OrganizationDao) GetByCreatedBy(userID uuid.UUID, page int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	query := Database.Model(&model.Organization{}).Where("created_by = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = query.Offset(offset).Limit(dao.Limit).Find(&orgs).Error
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (dao *OrganizationDao) Update(org *model.Organization) error {
	return Database.Save(org).Error
}

func (dao *OrganizationDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.Organization{}, "id = ?", id).Error
}

func (dao *OrganizationDao) UpdateStatus(id uuid.UUID, status model.OrganizationStatus) error {
	return Database.Model(&model.Organization{}).Where("id = ?", id).Update("status", status).Error
}

func (dao *OrganizationDao) SlugExists(slug string) (bool, error) {
	var count int64
	err := Database.Model(&model.Organization{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (dao *OrganizationDao) GetMembers(orgID uuid.UUID) ([]model.OrganizationMember, error) {
	memberDao := NewOrganizationMemberDao()
	return memberDao.GetMembers(orgID)
}
