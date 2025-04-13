package routes

import (
	"github.com/fsrn12/fitness_tracker_go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Get("/users/{id}", app.UserHandler.HandleGetUserByID)
	r.Post("/users", app.UserHandler.HandleCreateUser)
	return r
}
