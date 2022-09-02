package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	model TodoRepository
}

type TodoID struct {
	ID int `uri:"id" binding:"required"`
}

func (t todoHandler) getTodo(c *gin.Context) {
	var uri TodoID
	if c.BindUri(&uri) != nil {
		return
	}
	todo, found, err := t.model.FindOne(c.Request.Context(), uri.ID)
	if err != nil {
		internalServerError(c, err)
		return
	}
	if !found {
		todoNotFound(c)
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (t todoHandler) getTodos(c *gin.Context) {
	todos, err := t.model.FindAll(c.Request.Context())
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (t todoHandler) createTodo(c *gin.Context) {
	var todo Todo
	err := c.Bind(&todo)
	if err != nil {
		return
	}
	todo, err = t.model.Create(c.Request.Context(), todo)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (t todoHandler) editTodo(c *gin.Context) {
	var updateTodo UpdateTodo
	err := c.Bind(&updateTodo)
	if err != nil {
		return
	}
	var uri TodoID
	if c.BindUri(&uri) != nil {
		return
	}
	todo, err := t.model.Edit(c.Request.Context(), uri.ID, updateTodo)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, todo)
	return
}

func (t todoHandler) deleteTodo(c *gin.Context) {
	var uri TodoID
	if c.BindUri(&uri) != nil {
		return
	}
	err := t.model.Delete(c.Request.Context(), uri.ID)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"id": uri.ID})
	return
}

func todoNotFound(c *gin.Context) {
	_ = c.AbortWithError(http.StatusNotFound, &gin.Error{
		Err:  errors.New("todo not found"),
		Type: gin.ErrorTypePublic,
	})
}

func internalServerError(c *gin.Context, err error) {
	_ = c.AbortWithError(http.StatusInternalServerError, err)
}
