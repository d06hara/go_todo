package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

var (
	todos     []Todo
	mutex     sync.Mutex
	idCounter int = 1
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.DELETE("/todos/:id", deleteTodo)
	r.PUT("todos/:id", editTodo)

	r.Run()
}

func getTodos(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(200, todos)
}

func addTodo(c *gin.Context) {
	var newTodo Todo

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	newTodo.ID = idCounter
	idCounter++

	todos = append(todos, newTodo)
	mutex.Unlock()

	c.JSON(200, newTodo)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println("id", id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid todo ID"})
		return
	}

	mutex.Lock()
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.JSON(200, gin.H{"status": "Success"})
			mutex.Unlock()
			return
		}
	}
	mutex.Unlock()

	c.JSON(404, gin.H{"error": "Todo not found"})
}

func editTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid todo ID"})
		return
	}

	var updateTodo Todo
	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = updateTodo
			todos[i].ID = id
			c.JSON(200, todos[i])
			mutex.Unlock()
			return
		}
	}
	mutex.Unlock()
	c.JSON(404, gin.H{"error": "Todo not found"})
}
