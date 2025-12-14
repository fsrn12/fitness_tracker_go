package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=lol dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v\n", err)
	}

	// migrations
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v\n", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("migrating test db error: %v\n", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push Day",
				Description:     "Upper body strengthening session",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Squats",
						Sets:         5,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(100.0),
						Notes:        "Pushed harder",
						OrderIndex:   1,
					},
					{
						ExerciseName: "Lunges",
						Sets:         4,
						Reps:         IntPtr(12),
						Weight:       FloatPtr(55.0),
						Notes:        "Legs were shaking",
						OrderIndex:   2,
					},
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(12),
						Weight:       FloatPtr(135.5),
						Notes:        "warm up nicely",
						OrderIndex:   3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid payload",
			workout: &Workout{
				Title:           "Full Body",
				Description:     "Total workout",
				DurationMinutes: 60,
				CaloriesBurned:  900,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Weight:       FloatPtr(135.5),
						Notes:        "keep good form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            3,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(185.9),
						Notes:           "Full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			// if no errors are required
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrieved.Entries[i].Reps)
				assert.Equal(t, tt.workout.Entries[i].Notes, retrieved.Entries[i].Notes)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
