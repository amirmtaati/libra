package main

import (
	"flag"
	"fmt"
	"github.com/amirmtaati/libra/internal/app"
	"github.com/amirmtaati/libra/internal/handlers"
	"os"
)

func main() {
	var (
		dbPath  = flag.String("db", ".libra", "Database path")
		port    = flag.String("port", "8080", "Server port")
		libPath = flag.String("libpath", "", "Library path for scanning")
	)
	flag.Parse()

	config := &app.Config{
		DBPath:  *dbPath,
		LibPath: *libPath,
		Port:    *port,
	}

	application := app.NewApp(config)
	if err := application.Init(); err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
		os.Exit(1)
	}
	defer application.Shutdown()

	server := handlers.NewServer(application)
	fmt.Printf("🚀 API Server starting on port %s\n", *port)
	if err := server.Start(":" + *port); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
