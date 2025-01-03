package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DateTime struct {
	time.Time
}

func (t *DateTime) Scan(value interface{}) error {
	if value == nil {
		*t = DateTime{time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = DateTime{parsedTime}
	case time.Time:
		*t = DateTime{v}
	default:
		return fmt.Errorf("cannot scan type %T into DateTime", value)
	}
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format("2006-01-02 15:04:05"), nil
}

type Suggestion struct {
	Id       uint      `json:"id" gorm:"column:id;type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Title    string    `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Content  string    `json:"content" gorm:"column:content;type:text;not null"`
	Votes    int       `json:"votes" gorm:"column:votes;default:0"`
	Comments []Comment `json:"comments" gorm:"foreignKey:SuggestionId;references:Id"`
	Category string    `json:"category" gorm:"column:category;type:varchar(20);not null"`
	Status   string    `json:"status" gorm:"column:status;type:varchar(20);not null"`
	UserId   uint      `json:"user_id" gorm:"column:user_id;type:INT(10) UNSIGNED NOT NULL;index"`
	// User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"` we can use user_id to get user so we don't need to load user data
	CreatedAt DateTime `json:"created_at" gorm:"column:created_at;type:DATETIME"`
	UpdatedAt DateTime `json:"updated_at" gorm:"column:updated_at;type:DATETIME"`
	// DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

type Comment struct {
	Id           uint        `json:"id" gorm:"column:id;type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Content      string      `json:"content" gorm:"column:content;type:text;not null"`
	UserId       uint        `json:"user_id" gorm:"column:user_id;type:INT(10) UNSIGNED NOT NULL;index"`
	User         User        `json:"user" gorm:"foreignKey:UserId;references:Id"`
	SuggestionId uint        `json:"suggestion_id" gorm:"column:suggestion_id;type:INT(10) UNSIGNED NOT NULL;index"`
	Suggestion   *Suggestion `json:"suggestion" gorm:"foreignKey:SuggestionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Replies      []Reply     `json:"replies" gorm:"foreignKey:CommentId;references:Id"`
	CreatedAt    DateTime    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    DateTime    `json:"updated_at" gorm:"column:updated_at"`
	// DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

type User struct {
	Id          uint           `json:"id" gorm:"column:id;type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Username    string         `json:"username" gorm:"column:username;type:varchar(255);uniqueIndex;not null"`
	FirstName   string         `json:"firstName" gorm:"column:first_name;type:varchar(255)"`
	LastName    string         `json:"lastName" gorm:"column:last_name;type:varchar(255)"`
	Email       string         `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Avatar      *string        `json:"avatar,omitempty" gorm:"column:avatar;type:varchar(255)"` // Made nullable
	Password    string         `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Suggestions []Suggestion   `json:"suggestions" gorm:"foreignKey:UserId;references:Id"`
	Comments    []Comment      `json:"comments" gorm:"foreignKey:UserId;references:Id"`
	Replies     []Reply        `json:"replies" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt   DateTime       `json:"created_at" gorm:"column:created_at;type:DATETIME"`
	UpdatedAt   DateTime       `json:"updated_at" gorm:"column:updated_at;type:DATETIME"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

type Reply struct {
	Id        uint           `json:"id" gorm:"column:id;type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Content   string         `json:"content" gorm:"column:content;type:text;not null"`
	CommentId uint           `json:"comment_id" gorm:"column:comment_id;type:INT(10) UNSIGNED NOT NULL;index"`
	Comment   Comment        `json:"comment" gorm:"foreignKey:CommentId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserId    uint           `json:"user_id" gorm:"column:user_id;type:INT(10) UNSIGNED NOT NULL;index"`
	User      User           `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt DateTime       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt DateTime       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}
