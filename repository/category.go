package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	var listCategory []entity.Category

	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("user_id = ?", id).Find(&listCategory).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []entity.Category{}, nil
	}

	if err != nil {
		return nil, err
	}

	return listCategory, nil // TODO: replace this
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	err = r.db.WithContext(ctx).Model(&entity.Category{}).Create(&category).Error
	if err != nil {
		return 0, err		
	}

	return category.ID, nil // TODO: replace this
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	return r.db.WithContext(ctx).Model(&entity.Category{}).CreateInBatches(&categories, len(categories)).Error 
	// TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var categoryWithId entity.Category

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&categoryWithId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Category{}, nil
	}

	if err != nil {
		return entity.Category{}, err
	}

	return categoryWithId, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	return 	r.db.WithContext(ctx).Model(&entity.Category{}).Where("id = ?", category.ID).Updates(category).Error 
	// TODO: replace this
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	return 	r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error 
	 // TODO: replace this
}