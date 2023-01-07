package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	// TODO: answer here

	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))		
		return
	}

	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))


	taskStore, errStore := t.taskService.StoreTask(r.Context(), &entity.Task{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		CategoryID: task.CategoryID,
		UserID: idLogin,
	})

	if errStore != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("errir internal server"))
		return
	}

	resultTask := entity.ResponseTask{
		UserId: taskStore.UserID,
		TaskId: taskStore.ID,
		Message: "success create new task",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultTask)	
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))

	tId := r.URL.Query().Get("task_id")
	
	if tId == "" {
		taksUserId, errById := t.taskService.GetTasks(r.Context(), idLogin)
		if errById != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taksUserId)		
		return
	}

	taskID, _ := strconv.Atoi(tId)

	taskById, errId := t.taskService.GetTaskByID(r.Context(), taskID)
	if errId != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskById)			
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))


	taksID := r.URL.Query().Get("task_id")

	taskById, _ := strconv.Atoi(taksID)

	errDelete := t.taskService.DeleteTask(r.Context(), taskById)
	if errDelete != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return		
	}

	resultTask := entity.ResponseTask{
		UserId: idLogin,
		TaskId: taskById,
		Message: "success delete task",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultTask)	
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here
	userId := r.Context().Value("id")

	if userId.(string) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(userId.(string))

	taskUpdate, errUpdate := t.taskService.UpdateTask(r.Context(), &entity.Task{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		CategoryID: task.CategoryID,
		UserID: idLogin,
	})

	if errUpdate != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	resultTask := entity.ResponseTask{
		UserId: taskUpdate.UserID,
		TaskId: taskUpdate.ID,
		Message: "success update task",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultTask)
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
