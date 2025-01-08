package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	sql "feedback-io.backend/config"
	"feedback-io.backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSuggestions(c *fiber.Ctx) error {
	var suggestions []models.Suggestion
	var count int64

	offset, err_offset := strconv.Atoi(c.Query("offset", "0"))
	limit, err_limit := strconv.Atoi(c.Query("limit", "10"))

	if err_offset != nil || err_limit != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid offset parameter or limit parameter",
		})
	}

	category, err := strconv.Atoi(c.Query("category", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid category ID format",
		})
	}

	status := c.Query("status", "")

	// Build base query
	query := sql.DB.Model(&suggestions)

	// Apply filters
	if category != 0 {
		query = query.Where("category_id = ?", category)
	}
	if status != "" {
		query = query.Where("status_id = ?", status)
	}

	// Get filtered count
	if err := query.Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch suggestions count",
		})
	}

	// Get filtered and paginated results
	if err := query.Limit(limit).Offset(offset).Find(&suggestions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch suggestions",
		})
	}

	// fetch the comments

	return c.JSON(fiber.Map{
		"success": true,
		"data":    suggestions,
		"count":   count,
	})
}

func GetSuggestion(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid suggestion ID format",
		})
	}

	var suggestion models.Suggestion
	if err := sql.DB.First(&suggestion, &id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Suggestion not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch suggestion",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"data":    suggestion,
	})
}

func VoteSuggestion(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid suggestion ID",
		})
	}

	vote := c.Query("vote", "up")
	if vote != "up" && vote != "down" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vote parameter: must be 'up' or 'down'",
		})
	}

	// Start transaction
	tx := sql.DB.Begin()
	committed := false

	defer func() {
		if r := recover(); r != nil {
			if !committed {
				tx.Rollback()
			}
			panic(r)
		} else if tx.Error != nil && !committed {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to start transaction",
		})
	}

	// Use locking to prevent concurrent votes
	var suggestion models.Suggestion
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&suggestion, id).Error; err != nil {
		fmt.Printf("Error: %v\n", err)
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Suggestion not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch suggestion",
		})
	}

	voteChange := 1
	if vote == "down" {
		voteChange = -1
	}

	if err := tx.Model(&suggestion).
		Where("id = ?", id).
		Update("votes", gorm.Expr("votes + ?", voteChange)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update votes",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to commit transaction",
		})
	}

	// Fetch updated suggestion
	if err := sql.DB.First(&suggestion, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch updated suggestion",
		})
	}
	committed = true
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    suggestion,
	})
}

func CreateSuggestion(c *fiber.Ctx) error {
	type CreateSuggestionInput struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryId uint   `json:"category_id"`
		UserId     uint   `json:"user_id"`
	}

	var input CreateSuggestionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse request body",
		})
	}

	suggestion := models.Suggestion{
		Title:      input.Title,
		Content:    input.Content,
		CategoryId: input.CategoryId,
		UserId:     input.UserId,
		StatusId:   0,
	}

	if err := sql.DB.Create(&suggestion).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create suggestion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    suggestion,
	})
}

func DeleteSuggestion(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid suggestion ID",
		})
	}

	tx := sql.DB.Begin()

	var suggestions models.Suggestion // before we delete suggestion, we need to delete Replies, Comments, and then Suggestion

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.First(&suggestions, id).Error; err != nil {

		tx.Rollback()
		return err
	}

	// we need to fetch comments
	var comments []models.Comment

	if err := sql.DB.Model(&models.Comment{}).First(&comments, suggestions.Id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete all replies
	if len(comments) > 0 {
		for _, comment := range comments {
			if err := tx.Model(&models.Reply{}).Where("comment_id = ?", comment.Id).Update("deleted_at", time.Now()).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Model(&models.Comment{}).Where("suggestion_id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.Suggestion{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Suggestion deleted successfully",
		"data":    suggestions,
	})

}

func UpdateSuggestion(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid suggestion ID",
		})
	}

	type PathInput struct {
		Title      *string `json:"title"`
		Content    *string `json:"content"`
		CategoryId *uint   `json:"category_id"`
		StatusId   *uint   `json:"status_id"`
	}

	var body PathInput
	var suggestion models.Suggestion

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Data Provided",
		})
	}

	if err := sql.DB.First(&models.Suggestion{}, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Suggestion not found",
		})
	}

	updates := make(map[string]interface{})
	if body.Title != nil {
		updates["title"] = body.Title
	}
	if body.Content != nil {
		updates["content"] = body.Content
	}
	if body.CategoryId != nil {
		updates["category_id"] = body.CategoryId
	}
	if body.StatusId != nil {
		updates["status_id"] = body.StatusId
	}

	if err := sql.DB.Model(&models.Suggestion{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update suggestion",
		})
	}

	if err := sql.DB.First(&suggestion, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch updated suggestion",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": false,
		"message": "Suggestion updated successfully",
		"data":    suggestion,
	})
}
