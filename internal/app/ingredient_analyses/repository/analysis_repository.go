package repository

import (
	"gorm.io/gorm"
	"savebite/internal/domain/entity"
	"savebite/pkg/log"
)

type AnalysisRepoItf interface {
	CreateWithIngredients(analysis *entity.Analysis) error
}

type AnalysisRepo struct {
	db *gorm.DB
}

func NewAnalysisRepo(db *gorm.DB) AnalysisRepoItf {
	return &AnalysisRepo{db}
}

func (r *AnalysisRepo) CreateWithIngredients(analysis *entity.Analysis) error {
	err := r.db.Create(analysis).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisRepo][CreateWithIngredients] Failed to create analysis")
	}

	return err
}
