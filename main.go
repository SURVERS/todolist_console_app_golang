package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	addtodo "todolist-app/src/addtodo"
	"todolist-app/src/api/handle"
	"todolist-app/src/models"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

var allCommands = "Доступные команды:\nadd [текст] - Добавляет новую задачу\nlist - Посмотреть список задач\nchecked [ID-задачи] - Установить статус *Выполнена задача*\ndelete [ID-Задачи] - Удалить задачу по его ID\nexit - Выйти из программы"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки данных из .env")
	}

	connectUrl := os.Getenv("DATABASE_URL")
	server_port := os.Getenv("PORT")

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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		router := mux.NewRouter()
		router.HandleFunc("/add/{todo}", todoAdd)
		router.HandleFunc("/delete/{id:[0-9]+}", todoDelete)
		router.HandleFunc("/checked/{id:[0-9]+}", todoChecked)
		router.HandleFunc("/list", todoList).Methods("POST")
		http.Handle("/", router)

		fmt.Printf("Сервер http://localhost%s был успешно запущен!\n\n", server_port)
		if err := http.ListenAndServe(server_port, nil); err != nil {
			log.Fatalf("сервер не смог запуститься. Ошибка: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()
}

func todoList(w http.ResponseWriter, r *http.Request) {
	handle.HandleList(w)
}

func todoChecked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	handle.HandleChecked(id, pool, w)
}

func todoAdd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todo := vars["todo"]

	handle.HandleAdd(todo, pool, w)
}

func todoDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	handle.HandleDelete(id, pool, w)
}

func processCommand(text string) {
	parts := strings.Fields(text)

	command := parts[0]

	args := parts[1:]
	description := strings.Join(args, " ")

	switch command {
	case "add":
		addtodo.HandleAdd(description, pool)
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
