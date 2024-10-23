package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct
type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// SaveUser godoc
// @Summary Save a new user
// @Description Creates a new user and saves it to the database.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body User true "User Info"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /save [post]
func SaveUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves a user by their UUID.
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} User
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Server error"
// @Router /{id} [get]
func GetUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/"):]
		var user User

		result := db.First(&user, "id = ?", id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
