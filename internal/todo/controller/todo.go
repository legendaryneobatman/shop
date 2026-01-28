package controller

import (
	"github.com/gin-gonic/gin"
	"go-shop/internal/list"
)

type TodoController struct {
	listService list.Service
}

func (todoController *TodoController) CreateTodo(_ *gin.Context) {

}
func (todoController *TodoController) GetAllTodos(_ *gin.Context) {}
func (todoController *TodoController) GetTodoByID(_ *gin.Context) {}
func (todoController *TodoController) UpdateTodo(_ *gin.Context) {

}
func (todoController *TodoController) DeleteTodo(_ *gin.Context) {

}
