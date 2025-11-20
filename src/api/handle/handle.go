package handle

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"todolist-app/src/models"
	"todolist-app/src/response.go"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleDelete(idStr string, pool *pgxpool.Pool, w http.ResponseWriter) error {
	if len(idStr) == 0 {
		response.WriteJSONResponse(w, false, "Ошибка: Необходимо указать ID-задачи. Пример: /delete/1", http.StatusBadRequest)
		return nil
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.WriteJSONResponse(w, false, "Ошибка: Введите ID повторно", http.StatusBadRequest)
		return nil
	}

	if !containsOnlyDigits(idStr) {
		response.WriteJSONResponse(w, false, "Ошибка: Введите ID-задачи!", http.StatusBadRequest)
		return nil
	}

	for i := range models.Tasks {
		if models.Tasks[i].ID == id {
			query := "DELETE FROM todolist WHERE id = $1"
			_, err := pool.Exec(context.Background(), query, id)
			if err != nil {
				return fmt.Errorf("ошибка при удаление данных. Ошибка: %w", err)
			}
			response.WriteJSONResponse(w, true, fmt.Sprintf("Вы успешно удалили Задачу ID: %d.", id), http.StatusOK)
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			return nil
		}
	}

	response.WriteJSONResponse(w, false, fmt.Sprintf("Задача ID %d не найдена!", id), http.StatusBadRequest)
	return nil
}

func HandleChecked(idStr string, pool *pgxpool.Pool, w http.ResponseWriter) error {
	if len(idStr) == 0 {
		response.WriteJSONResponse(w, false, "Ошибка: Необходимо указать ID-задачи. Пример: /checked/1", http.StatusBadRequest)
		return nil
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.WriteJSONResponse(w, false, "Ошибка: Введите ID повторно", http.StatusBadRequest)
		return nil
	}
	if !containsOnlyDigits(idStr) {
		response.WriteJSONResponse(w, false, "Ошибка: Введите ID-задачи!", http.StatusBadRequest)
		return nil
	}
	for i := range models.Tasks {
		if models.Tasks[i].ID == id {
			query := "UPDATE todolist SET Completed = true WHERE id = $1"
			_, err := pool.Exec(context.Background(), query, id)
			if err != nil {
				response.WriteJSONResponse(w, false, fmt.Sprintf("ошибка при обновление данных. Ошибка: %s", err), http.StatusBadRequest)
				return nil
			}

			models.Tasks[i].Completed = true
			response.WriteJSONResponse(w, true, fmt.Sprintf("Вы успешно выполнили задачу ID: %d\n", models.Tasks[i].ID), http.StatusOK)
			return nil
		}
	}
	response.WriteJSONResponse(w, false, fmt.Sprintf("Задача под ID: %d не найдена!\n", id), http.StatusBadRequest)
	return nil
}

func HandleList(w http.ResponseWriter) {
	var allTodoList string

	if len(models.Tasks) == 0 {
		fmt.Println()
		response.WriteJSONResponse(w, false, "Ошибка: Список задач пуст.", http.StatusBadRequest)
		return
	}
	allTodoList = fmt.Sprintf("\n📋 Всего задач: %d шт.\n", len(models.Tasks))

	for i, todo := range models.Tasks {
		status := "[❌ ]"
		if todo.Completed {
			status = "[✅]"
		}
		allTodoList += fmt.Sprintf("%d) %s Задача ID: %d. %s\n", i+1, status, todo.ID, todo.Description)
	}

	response.WriteJSONResponse(w, true, allTodoList, http.StatusOK)
}

func HandleAdd(args string, pool *pgxpool.Pool, w http.ResponseWriter) error {
	if len(args) == 0 {
		response.WriteJSONResponse(w, false, "ошибка: Необходимо ввести описание задачи.", http.StatusBadRequest)
		return nil
	}
	query := `INSERT INTO todolist (description, completed) VALUES ($1, $2)`
	_, err := pool.Exec(context.Background(), query, args, false)

	if err != nil {
		response.WriteJSONResponse(w, false, fmt.Sprintf("ошибка при создании задания. %s", err), http.StatusBadRequest)
		return nil
	}

	newTask := models.Task{
		ID:          models.NextID,
		Description: args,
		Completed:   false,
	}

	models.Tasks = append(models.Tasks, newTask)
	models.NextID++

	response.WriteJSONResponse(w, true, fmt.Sprintf("Задача [ID %d] успешно добавлена: %s\n", models.NextID-1, newTask.Description), http.StatusOK)
	return nil
}

func containsOnlyDigits(s string) bool {
	if s == "" {
		return false
	}

	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
