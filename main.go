package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	addtodo "todolist-app/src/addtodo"
	"todolist-app/src/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

type Task struct {
	ID          int
	Description string
	Completed   bool
}

var tasks []Task
var allCommands = "Доступные команды:\nadd [текст] - Добавляет новую задачу\nlist - Посмотреть список задач\nchecked [ID-задачи] - Установить статус *Выполнена задача*\ndelete [ID-Задачи] - Удалить задачу по его ID\nexit - Выйти из программы"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки данных из .env")
	}

	connectUrl := os.Getenv("DATABASE_URL")

	ctx := context.Background()

	var errsql error
	pool, errsql = pgxpool.New(ctx, connectUrl)

	if errsql != nil {
		log.Fatalf("Ошибка при подключение к базе данных. Ошибка: Не удалось создать пул %s", errsql)
	}

	errsql = pool.Ping(ctx)
	if errsql != nil {
		log.Fatalf("Не удалось проверить подключение к базе данных: %s", errsql)
	}

	// Загружаем TodoList в структуру TodoList который находится в models.go
	addtodo.HandleLoad(pool)
	for i, todo := range models.Tasks {
		fmt.Printf("%d %s ID: %d\n", i, todo.Description, todo.ID)
	}

	fmt.Println("📝 ToDo List Console App")
	fmt.Println("---------------------------")
	fmt.Println(allCommands)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		processCommand(input)
	}
}

func processCommand(text string) {
	parts := strings.Fields(text)

	command := parts[0]

	args := parts[1:]

	switch command {
	case "add":
		addtodo.HandleAdd(args, pool)
	case "list":
		addtodo.HandleList()
	case "checked":
		addtodo.HandleChecked(args, pool)
	case "delete":
		addtodo.HandleDelete(args, pool)
	case "help":
		fmt.Println(allCommands)
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("Ошибка: Неизвестная команда! Введите /help для просмотра доступных функций.\n")
	}
}
