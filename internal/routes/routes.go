package routes

import (
	"github.com/fsrn12/fitness_tracker_go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
		r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
		r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
		r.Get("/users/{id}", app.UserHandler.HandleGetUserByID)
		r.Put("/users/{id}", app.UserHandler.HandleUpdateUserByID)
		r.Delete("/users/{id}", app.UserHandler.HandleDeleteUserByID)
	})

	r.Get("/health", app.HealthCheck)
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	// r.Post("/users", app.UserHandler.HandleCreateUser)
	return r
}
