package controller

import (
	"github.com/gin-gonic/gin"
)

type IListService interface{}

type TodoController struct {
	listService IListService
}

func (todoController *TodoController) CreateTodo(_ *gin.Context) {

}
func (todoController *TodoController) GetAllTodos(_ *gin.Context) {}
func (todoController *TodoController) GetTodoByID(_ *gin.Context) {}
func (todoController *TodoController) UpdateTodo(_ *gin.Context) {

}
func (todoController *TodoController) DeleteTodo(_ *gin.Context) {

}
