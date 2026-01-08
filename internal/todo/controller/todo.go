package controller

import (
	"github.com/gin-gonic/gin"
	"go-shop/internal/list/service"
)

type TodoController struct {
	listService service.ListService
}

func (todoController *TodoController) CreateTodo(c *gin.Context) {

}
func (todoController *TodoController) GetAllTodos(c *gin.Context) {}
func (todoController *TodoController) GetTodoById(c *gin.Context) {}
func (todoController *TodoController) UpdateTodo(c *gin.Context) {

}
func (todoController *TodoController) DeleteTodo(c *gin.Context) {

}
