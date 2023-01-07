package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))
	
	listCategory, err := c.categoryService.GetCategories(r.Context(), idLogin)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return		
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listCategory)	
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	// TODO: answer here
	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return		
	}

	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))

	categoryStore, errStore := c.categoryService.StoreCategory(r.Context(), &entity.Category{
		UserID: idLogin,
		Type: category.Type,
	})

	if errStore != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))		
		return
	}

	resultCategory := entity.ResponseCategory{
		UserId: idLogin,
		CategoryId: categoryStore.ID,
		Message: "success create new category",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultCategory)			
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))

	cId := r.URL.Query().Get("category_id")

	categoryId, _ := strconv.Atoi(cId)

	errDelete := c.categoryService.DeleteCategory(r.Context(), categoryId)

	if errDelete != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error internal server")		
		return
	}

	resultCategory := entity.ResponseCategory{
		UserId: idLogin,
		CategoryId: categoryId,
		Message: "success delete category",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultCategory)
	// TODO: answer here
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
