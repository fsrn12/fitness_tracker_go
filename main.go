package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/fsrn12/fitness_tracker_go/internal/app"
	"github.com/fsrn12/fitness_tracker_go/internal/routes"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close() // it will run after everything else
	// app.Logger.Printf("Server is running on port :%d", port)

	// http.HandleFunc("/heath", app.HealthCheck)
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is running on port:%d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
