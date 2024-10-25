package handlers

import (
	"log"
	"net/http"
	"nvrs-gateway/middleware"
	"nvrs-gateway/storage"
	"nvrs-gateway/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Agent represents an agent
type Agent struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

// Task represents a task submitted by an agent
type Task struct {
	AgentID int    `json:"agent_id"`
	Task    string `json:"task"`
}

// HealthCheck - Simple health check endpoint
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

// RegisterAgent handles agent registration
func RegisterAgent(c echo.Context) error {
	agent := new(Agent)
	if err := c.Bind(agent); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid input")
	}

	if agent.Name == "" || agent.Role == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Name and Role are required fields")
	}

	query := "INSERT INTO agents (name, role, status) VALUES (?, ?, ?) RETURNING id"
	err := storage.DB.QueryRow(query, agent.Name, agent.Role, "idle").Scan(&agent.ID)
	if err != nil {
		log.Printf("Error registering agent %s: %v", agent.Name, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not register agent")
	}

	token, err := middleware.GenerateToken(agent.ID, agent.Name)
	if err != nil {
		log.Printf("Error generating token for agent %s: %v", agent.Name, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not generate token")
	}

	log.Printf("Registered agent: ID=%d, Name=%s, Role=%s", agent.ID, agent.Name, agent.Role)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Agent registered successfully",
		"agent_id": agent.ID,
		"token":    token,
	})
}

// SubmitTask handles task submissions
func SubmitTask(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid input")
	}

	if task.AgentID == 0 || task.Task == "" {
		return utils.JSONError(c, http.StatusBadRequest, "AgentID and Task are required fields")
	}

	query := "INSERT INTO tasks (agent_id, task) VALUES (?, ?)"
	_, err := storage.DB.Exec(query, task.AgentID, task.Task)
	if err != nil {
		log.Printf("Error submitting task for agent %d: %v", task.AgentID, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not submit task")
	}

	log.Printf("Task submitted: AgentID=%d, Task=%s", task.AgentID, task.Task)

	return c.JSON(http.StatusOK, map[string]string{"message": "Task submitted successfully"})
}

// UpdateStatus handles status updates
func UpdateStatus(c echo.Context) error {
	var update struct {
		AgentID int    `json:"agent_id"`
		Status  string `json:"status"`
	}
	if err := c.Bind(&update); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid input")
	}

	if update.AgentID == 0 || update.Status == "" {
		return utils.JSONError(c, http.StatusBadRequest, "AgentID and Status are required fields")
	}

	query := "UPDATE agents SET status = ? WHERE id = ?"
	result, err := storage.DB.Exec(query, update.Status, update.AgentID)
	if err != nil {
		log.Printf("Error updating status for agent %d: %v", update.AgentID, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not update status")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Agent not found for status update: AgentID=%d", update.AgentID)
		return utils.JSONError(c, http.StatusNotFound, "Agent not found")
	}

	log.Printf("Status updated: AgentID=%d, NewStatus=%s", update.AgentID, update.Status)

	return c.JSON(http.StatusOK, map[string]string{"message": "Status updated successfully"})
}

// QueryTasks retrieves tasks for an agent
func QueryTasks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, ok := claims["agent_id"].(float64)
	if !ok {
		log.Println("Invalid or missing agent ID in token")
		return utils.JSONError(c, http.StatusBadRequest, "Invalid or missing agent ID in token")
	}

	intAgentID := int(agentID)

	rows, err := storage.DB.Query("SELECT id, task FROM tasks WHERE agent_id = ?", intAgentID)
	if err != nil {
		log.Printf("Failed to fetch tasks for agent %d: %v", intAgentID, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not fetch tasks")
	}
	defer rows.Close()

	tasks := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var task string
		if err := rows.Scan(&id, &task); err != nil {
			log.Printf("Error scanning task for agent %d: %v", intAgentID, err)
			return utils.JSONError(c, http.StatusInternalServerError, "Error processing tasks")
		}
		tasks = append(tasks, map[string]interface{}{"id": id, "task": task})
	}

	return c.JSON(http.StatusOK, tasks)
}

// GetStatus retrieves the current status of the agent
func GetStatus(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, ok := claims["agent_id"].(float64)
	if !ok {
		log.Println("Invalid or missing agent ID in token")
		return utils.JSONError(c, http.StatusBadRequest, "Invalid or missing agent ID in token")
	}

	intAgentID := int(agentID)

	var status string
	err := storage.DB.QueryRow("SELECT status FROM agents WHERE id = ?", intAgentID).Scan(&status)
	if err != nil {
		log.Printf("Could not fetch status for agent %d: %v", intAgentID, err)
		return utils.JSONError(c, http.StatusInternalServerError, "Could not fetch status")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": status})
}
