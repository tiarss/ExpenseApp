package controllers

import (
	"encoding/json"
	"expense-app-backend/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRequest struct {
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
	SubCategory  []struct {
		Name string `json:"name"`
	} `json:"sub_category"`
}

type Category struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string
	CategoryType string
	SubCategory  []SubCategory `gorm:"foreignKey:CategoryID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SubCategory struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CategoryID uuid.UUID
}

type subCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type categoryResponse struct {
	ID           uuid.UUID             `json:"id"`
	Name         string                `json:"name"`
	CategoryType string                `json:"category_type"`
	SubCategory  []subCategoryResponse `json:"sub_categories"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}

func GetCategories(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		name := r.URL.Query().Get("name")
		categoryType := r.URL.Query().Get("category_type")

		query := db
		if name != "" {
			query = query.Where("name = ?", name)
		}
		if categoryType != "" {
			query = query.Where("category_type = ?", categoryType)
		}

		var categories []Category
		result := query.Preload("SubCategory", "deleted_at IS NULL").Where("deleted_at IS NULL").Limit(10).Find(&categories)
		if result.Error != nil {
			http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
			return
		}

		var response []categoryResponse
		for _, category := range categories {
			subCategoryResponses := make([]subCategoryResponse, len(category.SubCategory))
			for i, subCategory := range category.SubCategory {
				subCategoryResponses[i] = subCategoryResponse{
					ID:   subCategory.ID.String(),
					Name: subCategory.Name,
				}
			}

			response = append(response, categoryResponse{
				ID:           category.ID,
				Name:         category.Name,
				CategoryType: category.CategoryType,
				SubCategory:  subCategoryResponses,
				CreatedAt:    category.CreatedAt.Format(time.RFC3339),
				UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if response == nil {
			response = []categoryResponse{}
		}

		responseJson := struct {
			StatusCode int                `json:"status_code"`
			Data       []categoryResponse `json:"data"`
			Message    string             `json:"message"`
		}{
			StatusCode: http.StatusOK,
			Data:       response,
			Message:    "Categories successfully retrieved",
		}

		if err := json.NewEncoder(w).Encode(responseJson); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetCategoryById(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var id = r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		var category Category
		result := db.Preload("SubCategory", "deleted_at IS NULL").Where("id = ? AND deleted_at IS NULL", id).First(&category)
		if result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusNotFound,
				Message:    "Category not found",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		var response categoryResponse
		subCategoryResponses := make([]subCategoryResponse, len(category.SubCategory))

		for i, subCategory := range category.SubCategory {
			subCategoryResponses[i] = subCategoryResponse{
				ID:   subCategory.ID.String(),
				Name: subCategory.Name,
			}
		}

		response = categoryResponse{
			ID:           category.ID,
			Name:         category.Name,
			CategoryType: category.CategoryType,
			SubCategory:  subCategoryResponses,
			CreatedAt:    category.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
		}

		// var count int64
		// db.Model(&Category{}).Count(&count)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		responseJson := struct {
			StatusCode int              `json:"status_code"`
			Data       categoryResponse `json:"data"`
			Message    string           `json:"message"`
			// Count      int              `json:"count"`
		}{
			StatusCode: http.StatusOK,
			Data:       response,
			Message:    "Categories successfully retrieved",
			// Count:      int(count),
		}

		json.NewEncoder(w).Encode(responseJson)
	}
}

func CreateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var categoryRequest CategoryRequest
		if err := decoder.Decode(&categoryRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		defer r.Body.Close()

		if categoryRequest.Name == "" || categoryRequest.CategoryType == "" {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Name and category type are required",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		var existingCategory models.Category
		if err := db.Where("name = ?", categoryRequest.Name).First(&existingCategory).Error; err == nil {
			w.WriteHeader(http.StatusConflict)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusConflict,
				Message:    "Category already exists",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		category := models.Category{
			ID:            uuid.New(),
			Name:          categoryRequest.Name,
			CategoryType:  categoryRequest.CategoryType,
			SubCategories: []models.SubCategory{},
		}

		for _, subCategory := range categoryRequest.SubCategory {
			category.SubCategories = append(category.SubCategories, models.SubCategory{
				ID:   uuid.New(),
				Name: subCategory.Name,
			})
		}

		if err := db.Create(&category).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create category",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)

		type subCategoryResponse struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		type DataResponse struct {
			ID            string                `json:"id"`
			Name          string                `json:"name"`
			CategoryType  string                `json:"category_type"`
			SubCategories []subCategoryResponse `json:"sub_categories"`
			CreatedAt     string                `json:"created_at"`
			UpdatedAt     string                `json:"updated_at"`
		}

		var subCategories []subCategoryResponse
		for _, subCategory := range category.SubCategories {
			subCategories = append(subCategories, subCategoryResponse{ID: subCategory.ID.String(), Name: subCategory.Name})
		}

		response := struct {
			StatusCode int          `json:"status_code"`
			Data       DataResponse `json:"data"`
			Message    string       `json:"message"`
		}{
			StatusCode: http.StatusCreated,
			Data: DataResponse{
				ID:            category.ID.String(),
				Name:          category.Name,
				CategoryType:  category.CategoryType,
				CreatedAt:     category.CreatedAt.Format(time.RFC3339),
				UpdatedAt:     category.UpdatedAt.Format(time.RFC3339),
				SubCategories: subCategories,
			},
			Message: "Category created successfully",
		}

		json.NewEncoder(w).Encode(response)
	}
}

func UpdateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusBadRequest,
				"message":     "Category ID is required",
			})
			return
		}

		var categoryRequest struct {
			Name         string `json:"name"`
			CategoryType string `json:"category_type"`
			SubCategory  []struct {
				Name string `json:"name"`
			} `json:"sub_category"`
		}

		if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusBadRequest,
				"message":     "Invalid request body",
			})
			return
		}
		defer r.Body.Close()

		if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&models.Category{}).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusNotFound,
				"message":     "Category not found",
			})
			return
		}

		if categoryRequest.Name == "" || categoryRequest.CategoryType == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusBadRequest,
				"message":     "Name and category type are required",
			})
			return
		}

		tx := db.Begin()

		if err := tx.Model(&models.Category{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Updates(map[string]interface{}{
				"name":          categoryRequest.Name,
				"category_type": categoryRequest.CategoryType,
				"updated_at":    time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusInternalServerError,
				"message":     "Failed to update category",
			})
			return
		}

		if err := tx.Where("category_id = ?", id).Delete(&models.SubCategory{}).Error; err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status_code": http.StatusInternalServerError,
				"message":     "Failed to delete existing subcategories",
			})
			return
		}

		var newSubCategories []models.SubCategory
		for _, subCategory := range categoryRequest.SubCategory {
			newSubCategories = append(newSubCategories, models.SubCategory{
				ID:         uuid.New(),
				Name:       subCategory.Name,
				CategoryID: uuid.MustParse(id),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}

		if len(newSubCategories) > 0 {
			if err := tx.Create(&newSubCategories).Error; err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status_code": http.StatusInternalServerError,
					"message":     "Failed to create new subcategories",
				})
				return
			}
		}

		tx.Commit()

		type SubCategoryResponse struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		var responseSubCategories []SubCategoryResponse
		for _, subCategory := range newSubCategories {
			responseSubCategories = append(responseSubCategories, SubCategoryResponse{
				ID:   subCategory.ID.String(),
				Name: subCategory.Name,
			})
		}

		successResponse := struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Data       struct {
				ID            string                `json:"id"`
				Name          string                `json:"name"`
				CategoryType  string                `json:"category_type"`
				SubCategories []SubCategoryResponse `json:"sub_categories"`
				UpdatedAt     string                `json:"updated_at"`
			} `json:"data"`
		}{
			StatusCode: http.StatusOK,
			Message:    "Category and subcategories updated successfully",
			Data: struct {
				ID            string                `json:"id"`
				Name          string                `json:"name"`
				CategoryType  string                `json:"category_type"`
				SubCategories []SubCategoryResponse `json:"sub_categories"`
				UpdatedAt     string                `json:"updated_at"`
			}{
				ID:            id,
				Name:          categoryRequest.Name,
				CategoryType:  categoryRequest.CategoryType,
				SubCategories: responseSubCategories,
				UpdatedAt:     time.Now().Format(time.RFC3339),
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(successResponse)
	}
}

func DeleteCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id = r.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Category ID is required",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		var category models.Category
		if err := db.Preload("SubCategories").Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusNotFound,
				Message:    "Category not found",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		if err := db.Where("category_id = ?", category.ID).Delete(&models.SubCategory{}).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete associated subcategories",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete category",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusOK)
		successResponse := struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
		}{
			StatusCode: http.StatusOK,
			Message:    "Category deleted successfully",
		}
		json.NewEncoder(w).Encode(successResponse)
	}
}
