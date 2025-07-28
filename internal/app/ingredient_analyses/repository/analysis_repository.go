package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/entity"
	"savebite/pkg/log"
)

type AnalysisRepoItf interface {
	Create(analysis *entity.Analysis) error
	GetByUserID(userID uuid.UUID, pagination dto.PaginationRequest) (*[]entity.Analysis, int64, error)
}

type AnalysisRepo struct {
	db *gorm.DB
}

func NewAnalysisRepo(db *gorm.DB) AnalysisRepoItf {
	return &AnalysisRepo{db}
}

func (r *AnalysisRepo) GetByUserID(userID uuid.UUID, pagination dto.PaginationRequest) (*[]entity.Analysis, int64, error) {
	var totalItems int64
	var analyses []entity.Analysis
	err := r.db.Model(&entity.Analysis{}).Count(&totalItems).Preload("Ingredients").Limit(pagination.Limit).Offset(pagination.Offset).Find(&analyses, entity.Analysis{
		UserID: userID,
	}).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisRepo][GetByUserID] Failed to retrieve analyses")
	}

	return &analyses, totalItems, err
}

func (r *AnalysisRepo) Create(analysis *entity.Analysis) error {
	err := r.db.Create(analysis).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisRepo][Create] Failed to create analysis")
	}

	return err
}
