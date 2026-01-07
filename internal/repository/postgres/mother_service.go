package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type motherServiceRepository struct {
	db *gorm.DB
}

func NewMotherServiceRepository(db *gorm.DB) repository.MotherServiceRepository {
	log.Println("initializing mother service repository")

	return &motherServiceRepository{
		db: db,
	}
}

func (m *motherServiceRepository) Create(ctx context.Context, motherService *entity.MotherService) error {
	log.Printf("creating mother service record with input: %s", motherService.Name)

	err := m.db.WithContext(ctx).Create(motherService).Error
	if err != nil {
		log.Printf("failed to create mother service record with input '%s': %v", motherService.Name, err)

		return fmt.Errorf("failed to create mother service record: %w", err)
	}

	log.Printf("successfully created mother service record with ID: %d, input: %s", motherService.ID, motherService.Name)

	return nil
}
