package queue

import "gorm.io/gorm"

type QueueRepository struct {
	DB *gorm.DB
}

func NewQueueRepository(db *gorm.DB) *QueueRepository {
	return &QueueRepository{
		DB: db,
	}
}

func (repo *QueueRepository) Create(db *gorm.DB, payload *QueueRecord) error {
	return db.Create(payload).Error
}

func (repo *QueueRepository) UpdateStatus(uniqueID string, status Status) error {
	var record QueueRecord
	err := repo.DB.Where("unique_id = ?", uniqueID).First(&record).Error
	if err != nil {
		return err
	}

	record.Status = status
	return repo.DB.Save(&record).Error
}

func (repo *QueueRepository) FindByUniqueID(uniqueID string) (*QueueRecord, error) {
	var record QueueRecord
	err := repo.DB.Where("unique_id = ?", uniqueID).First(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}
