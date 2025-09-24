package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Todo represents a todo item in the database
type Todo struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	CreateTime string    `json:"create_time" db:"create_time"`
	Importance int       `json:"importance" db:"importance"`
	Flag       bool      `json:"flag" db:"flag"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTodoRequest represents the request body for creating a new todo
type CreateTodoRequest struct {
	Name       string `json:"name" binding:"required,max=500"`
	CreateTime string `json:"create_time,omitempty"`
	Importance int    `json:"importance,omitempty"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Name       *string `json:"name,omitempty"`
	CreateTime *string `json:"create_time,omitempty"`
	Importance *int    `json:"importance,omitempty"`
	Flag       *bool   `json:"flag,omitempty"`
}

// ReorderRequest represents the request body for reordering todos
type ReorderRequest struct {
	TodoIDs []int `json:"todo_ids" binding:"required"`
}

type TodoHandler struct {
	db *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

// GetTodos retrieves all todos ordered by importance (highest first), then by create_time
func (h *TodoHandler) GetTodos(c *gin.Context) {
	query := `
		SELECT id, name, create_time, importance, flag, created_at, updated_at 
		FROM todos 
		ORDER BY importance DESC, create_time ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.CreateTime, &todo.Importance,
			&todo.Flag, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan todo"})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default create_time to today if not provided
	if req.CreateTime == "" {
		req.CreateTime = time.Now().Format("2006-01-02")
	}

	// Get the highest current importance and add 1 for new todo
	var maxImportance int
	err := h.db.QueryRow("SELECT COALESCE(MAX(importance), 0) FROM todos").Scan(&maxImportance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get max importance"})
		return
	}

	if req.Importance == 0 {
		req.Importance = maxImportance + 1
	}

	query := `
		INSERT INTO todos (name, create_time, importance, flag) 
		VALUES ($1, $2, $3, false) 
		RETURNING id, name, create_time, importance, flag, created_at, updated_at
	`

	var todo Todo
	err = h.db.QueryRow(query, req.Name, req.CreateTime, req.Importance).Scan(
		&todo.ID, &todo.Name, &todo.CreateTime, &todo.Importance,
		&todo.Flag, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": todo})
}

// UpdateTodo updates an existing todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		setParts = append(setParts, "name = $"+strconv.Itoa(argCount))
		args = append(args, *req.Name)
		argCount++
	}
	if req.CreateTime != nil {
		setParts = append(setParts, "create_time = $"+strconv.Itoa(argCount))
		args = append(args, *req.CreateTime)
		argCount++
	}
	if req.Importance != nil {
		setParts = append(setParts, "importance = $"+strconv.Itoa(argCount))
		args = append(args, *req.Importance)
		argCount++
	}
	if req.Flag != nil {
		setParts = append(setParts, "flag = $"+strconv.Itoa(argCount))
		args = append(args, *req.Flag)
		argCount++
	}

	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := "UPDATE todos SET " + strings.Join(setParts, ", ") +
		", updated_at = NOW() WHERE id = $" + strconv.Itoa(argCount) +
		" RETURNING id, name, create_time, importance, flag, created_at, updated_at"

	args = append(args, id)

	var todo Todo
	err = h.db.QueryRow(query, args...).Scan(
		&todo.ID, &todo.Name, &todo.CreateTime, &todo.Importance,
		&todo.Flag, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

// DeleteTodo deletes a todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	result, err := h.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// ReorderTodos handles drag-and-drop reordering by updating importance values
func (h *TodoHandler) ReorderTodos(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update importance based on array position (reverse order for highest first)
	for i, todoID := range req.TodoIDs {
		importance := len(req.TodoIDs) - i // Highest importance for first item
		_, err := tx.Exec("UPDATE todos SET importance = $1, updated_at = NOW() WHERE id = $2",
			importance, todoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder todos"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todos reordered successfully"})
}
