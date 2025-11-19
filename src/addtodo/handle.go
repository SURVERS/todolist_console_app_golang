package addtodo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"todolist-app/src/models"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleDelete(args []string, pool *pgxpool.Pool) error {
	if len(args) == 0 {
		fmt.Println("Ошибка: Необходимо указать ID-задачи. Пример: delete 1")
		return nil
	}

	idStr := strings.Join(args, "")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		fmt.Println("Ошибка: Введите ID повторно")
		return nil
	}

	if !containsOnlyDigits(idStr) {
		fmt.Println("Ошибка: Введите ID-задачи!")
		return nil
	}

	for i := range models.Tasks {
		if models.Tasks[i].ID == id {
			query := "DELETE FROM todolist WHERE id = $1"
			_, err := pool.Exec(context.Background(), query, id)
			if err != nil {
				return fmt.Errorf("ошибка при удаление данных. Ошибка: %w", err)
			}
			fmt.Printf("✅ Вы успешно удалили Задачу ID: %d.\n", id)
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			return nil
		}
	}

	fmt.Printf("❌ Задача ID %d не найдена!\n", id)
	return nil
}

func HandleChecked(args []string, pool *pgxpool.Pool) error {
	if len(args) == 0 {
		fmt.Println("Ошибка: Необходимо указать ID-задачи. Пример: checked 1")
		return nil
	}

	idStr := strings.Join(args, " ")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		fmt.Println("Ошибка: Введите ID повторно")
		return nil
	}
	if !containsOnlyDigits(idStr) {
		fmt.Println("Ошибка: Введите ID-задачи!")
		return nil
	}
	for i := range models.Tasks {
		if models.Tasks[i].ID == id {
			query := "UPDATE todolist SET Completed = true WHERE id = $1"
			_, err := pool.Exec(context.Background(), query, id)
			if err != nil {
				return fmt.Errorf("ошибка при обновление данных. Ошибка: %w", err)
			}

			models.Tasks[i].Completed = true
			fmt.Printf("✅ Вы успешно выполнили задачу ID: %d\n", models.Tasks[i].ID)
			return nil
		}
	}
	fmt.Printf("❌ Задача под ID: %d не найдена!\n", id)
	return nil
}

func HandleList() {
	if len(models.Tasks) == 0 {
		fmt.Println("Ошибка: Список задач пуст.")
		return
	}
	fmt.Printf("\n📋 Всего задач: %d шт.\n", len(models.Tasks))
	for i, todo := range models.Tasks {
		status := "[❌ ]"
		if todo.Completed {
			status = "[✅]"
		}
		fmt.Printf("%d) %s Задача ID: %d. %s\n", i+1, status, todo.ID, todo.Description)
	}
}

func HandleLoad(pool *pgxpool.Pool) error {
	query := `SELECT * FROM todolist`
	rows, err := pool.Query(context.Background(), query)

	if err != nil {
		return fmt.Errorf("ошибка при запросе: %w", err)
	}
	defer rows.Close()

	var items []models.Task

	for rows.Next() {
		var item models.Task
		err := rows.Scan(&item.ID, &item.Description, &item.Completed)
		if err != nil {
			return fmt.Errorf("ошибка при сканирование данных: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	if len(items) == 0 {
		return fmt.Errorf("список задач пустой. Код: %d", len(items))
	}

	models.Tasks = append(models.Tasks, items...)
	models.NextID = len(models.Tasks) + 1
	return nil
}

func HandleAdd(args []string, pool *pgxpool.Pool) error {
	if len(args) == 0 {
		return fmt.Errorf("ошибка: Необходимо ввести описание задачи. Пример: add Купить хлеб")
	}

	description := strings.Join(args, " ")
	query := `INSERT INTO todolist (description, completed) VALUES ($1, $2)`
	_, err := pool.Exec(context.Background(), query, description, false)

	if err != nil {
		return fmt.Errorf("ошибка при создании задания. %w", err)
	}

	newTask := models.Task{
		ID:          models.NextID,
		Description: description,
		Completed:   false,
	}

	models.Tasks = append(models.Tasks, newTask)
	models.NextID++

	fmt.Printf("✅ Задача [ID %d] успешно добавлена: %s\n", models.NextID-1, newTask.Description)
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
