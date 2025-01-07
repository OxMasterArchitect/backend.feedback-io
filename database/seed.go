package database

import (
	"fmt"

	"feedback-io.backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed function to populate the database with initial data
func Seed(db *gorm.DB) error {
	// Clear existing data
	if err := cleanDatabase(db); err != nil {
		return fmt.Errorf("error cleaning database: %v", err)
	}

	// Create users
	users, err := createUsers(db)
	if err != nil {
		return fmt.Errorf("error creating users: %v", err)
	}

	// Create suggestions
	suggestions, err := createSuggestions(db, users)
	if err != nil {
		return fmt.Errorf("error creating suggestions: %v", err)
	}

	// Create comments
	comments, err := createComments(db, users, suggestions)
	if err != nil {
		return fmt.Errorf("error creating comments: %v", err)
	}

	// Create replies
	if err := createReplies(db, users, comments); err != nil {
		return fmt.Errorf("error creating replies: %v", err)
	}

	// Create categories
	if _, err := createCategories(db); err != nil {
		return fmt.Errorf("error creating categories: %v", err)
	}

	fmt.Println("Database seeded successfully!")
	return nil
}

func cleanDatabase(db *gorm.DB) error {
	// Drop tables in correct order to avoid foreign key constraints
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Reply{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Comment{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Suggestion{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Category{}).Error; err != nil {
		return err
	}

	return db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
}

func createUsers(db *gorm.DB) ([]models.User, error) {
	// Create test users
	users := []models.User{
		{
			Username:  "john_doe",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  hashPassword("password123"),
			Avatar:    nil,
		},
		{
			Username:  "jane_smith",
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
			Password:  hashPassword("password123"),
			Avatar:    nil,
		},
		{
			Username:  "bob_wilson",
			FirstName: "Bob",
			LastName:  "Wilson",
			Email:     "bob@example.com",
			Password:  hashPassword("password123"),
			Avatar:    nil,
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

func createSuggestions(db *gorm.DB, users []models.User) ([]models.Suggestion, error) {
	suggestions := []models.Suggestion{
		{
			Title:      "Improve Website Performance",
			Content:    "We should optimize our website loading times by implementing caching and reducing image sizes.",
			Votes:      5,
			Status:     "in-progress",
			CategoryId: 0,
			UserId:     users[0].Id,
		},
		{
			Title:      "Add Dark Mode",
			Content:    "Implement a dark mode theme for better user experience during night time usage.",
			Votes:      10,
			Status:     "in-progress",
			CategoryId: 1,
			UserId:     users[1].Id,
		},
		{
			Title:      "Mobile App Development",
			Content:    "We should create a mobile app version of our platform for better accessibility.",
			Votes:      8,
			Status:     "in-progress",
			CategoryId: 2,
			UserId:     users[2].Id,
		},
	}

	for i := range suggestions {
		if err := db.Create(&suggestions[i]).Error; err != nil {
			return nil, err
		}
	}

	return suggestions, nil
}
func createCategories(db *gorm.DB) ([]models.Category, error) {
	categories := []models.Category{
		{
			Name: "Performance",
		},
		{
			Name: "Design",
		},
		{
			Name: "Feature",
		},
		{
			Name: "Bug",
		},
		{
			Name: "Enhancement",
		},
		{
			Name: "UI",
		},
		{
			Name: "UX",
		},
		{
			Name: "Security",
		},
		{
			Name: "Accessibility",
		},
		{
			Name: "Other",
		},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func createComments(db *gorm.DB, users []models.User, suggestions []models.Suggestion) ([]models.Comment, error) {
	comments := []models.Comment{
		{
			Content:      "Great idea! This would definitely improve user experience.",
			UserId:       users[1].Id,
			SuggestionId: suggestions[0].Id,
		},
		{
			Content:      "I've been waiting for dark mode for a long time!",
			UserId:       users[2].Id,
			SuggestionId: suggestions[1].Id,
		},
		{
			Content:      "Mobile app would be amazing for on-the-go access.",
			UserId:       users[0].Id,
			SuggestionId: suggestions[2].Id,
		},
	}

	for i := range comments {
		if err := db.Create(&comments[i]).Error; err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func createReplies(db *gorm.DB, users []models.User, comments []models.Comment) error {
	replies := []models.Reply{
		{
			Content:   "Agreed! Perhaps we could start with basic caching.",
			UserId:    users[2].Id,
			CommentId: comments[0].Id,
		},
		{
			Content:   "Would love to help test the dark mode when it's ready!",
			UserId:    users[0].Id,
			CommentId: comments[1].Id,
		},
		{
			Content:   "What platforms should we target first? iOS or Android?",
			UserId:    users[1].Id,
			CommentId: comments[2].Id,
		},
	}

	for i := range replies {
		if err := db.Create(&replies[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func hashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// In a real application, you'd want to handle this error appropriately
		return ""
	}
	return string(hashedBytes)
}
