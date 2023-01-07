package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var task []entity.Task

	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("user_id = ?", id).Find(&task).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []entity.Task{}, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	err = r.db.WithContext(ctx).Model(&entity.Task{}).Create(&task).Error
	if err != nil {
		return 0, err
	}

	return task.ID, nil // TODO: replace this
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var taskWithId entity.Task

	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).First(&taskWithId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Task{}, nil
	}

	if err != nil {
		return entity.Task{}, err
	}

	return taskWithId, nil // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	var listTaskWithCatId []entity.Task

	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("category_id = ?", catId).Find(&listTaskWithCatId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []entity.Task{}, nil
	}

	if err != nil {
		return nil, err
	}

	return listTaskWithCatId, nil // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	return 	r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", task.ID).Updates(&task).Error
	// TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {	
	return 	r.db.WithContext(ctx).Delete(&entity.Task{}, id).Error
	// TODO: replace this
}
