package routes

import (
	"expense-app-backend/controllers"
	"expense-app-backend/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", controllers.Register(db)).Methods("POST")
	router.HandleFunc("/api/login", controllers.Login(db)).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/categories", controllers.GetCategories(db)).Methods("GET")
	protected.HandleFunc("/categories", controllers.CreateCategory(db)).Methods("POST")
	// protected.HandleFunc("/categories/{id}", controllers.GetCategoryById(db)).Methods("GET")
	// protected.HandleFunc("/categories/{id}", controllers.UpdateCategory(db)).Methods("PUT")
	// protected.HandleFunc("/categories/{id}", controllers.DeleteCategory(db)).Methods("DELETE")
	// 	protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
	// 		w.Write([]byte("Protected route"))
	// }).Methods("GET")

	return router
	// protected.HandleFunc("/expenses", controllers.GetExpenses(db)).Methods("GET")
	// protected.HandleFunc("/expenses", controllers.CreateExpense(db)).Methods("POST")
	// protected.HandleFunc("/expenses/{id}", controllers.GetExpense(db)).Methods("GET")
	// protected.HandleFunc("/expenses/{id}", controllers.UpdateExpense(db)).Methods("PUT")
	// protected.HandleFunc("/expenses/{id}", controllers.DeleteExpense(db)).Methods("DELETE")
}
