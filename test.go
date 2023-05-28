package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todos struct {
	ID        string
	text      string
	idCounter int
}

var (
	todoss []Todo
)

func test() {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/todos", get)
}

func get(c *gin.Context) {
	// 排他制御
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(200, todoss)
}

func add(c *gin.Context) {
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
