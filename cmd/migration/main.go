package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/canonflow/canonflow-go-ddd/internal/factory"
	"github.com/canonflow/canonflow-go-ddd/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)

	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	m := newMigrate(viperConfig, logrus)
	defer m.Close()

	command := args[0]

	switch command {
	case "up":
		runUp(m, logrus)
	case "down":
		if len(args) > 1 {
			steps, err := strconv.Atoi(args[1])
			if err != nil {
				logrus.Fatalf("Error: invalid steps '%s', must be an integer", args[1])
			}
			runDownSteps(m, logrus, steps)
		} else {
			runDown(m, logrus)
		}
	case "version":
		runVersion(m, logrus)
	case "force":
		if len(args) < 2 {
			logrus.Fatal("Error: force requires a version argument. Usage: go run cmd/migrate/main.go force <version>")
		}
		version, err := strconv.Atoi(args[1])
		if err != nil {
			logrus.Fatalf("Error: invalid version '%s', must be an integer", args[1])
		}
		runForce(m, logrus, version)
	case "drop":
		runDrop(m, logrus)
	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func newMigrate(viper *viper.Viper, logrus *logrus.Logger) *migrate.Migrate {
	driver := strings.ToLower(viper.GetString("DB_DRIVER"))
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetInt("DB_PORT")
	database := viper.GetString("DB_NAME")

	if !utils.SliceContains(config.AVAILABLE_DRIVERS, driver) {
		logrus.Fatal("Unsupported Database Driver")
	}

	databaseDriver := factory.NewDatabaseFactory(driver, username, password, host, port, database)
	if databaseDriver == nil {
		logrus.Fatal("Unsupported Database Driver")
	}

	dsn := databaseDriver.GetDSN()

	if driver == "" {
		logrus.Fatal("Error: DB_DRIVER is not set in .env")
	}
	if dsn == "" {
		logrus.Fatal("Error: DB_DSN is not set in .env")
	}

	logrus.Info(dsn)

	migrationPath := fmt.Sprintf("file://./migrations/%s", driver)

	m, err := migrate.New(migrationPath, fmt.Sprintf("%s://%s", driver, dsn))
	if err != nil {
		logrus.Fatalf("Error creating migrate instance: %v", err)
	}

	return m
}

func runUp(m *migrate.Migrate, logrus *logrus.Logger) {
	logrus.Println("Running migrations up...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error running up: %v", err)
	}
	logrus.Println("Migrations up completed successfully.")
}

func runDown(m *migrate.Migrate, logrus *logrus.Logger) {
	logrus.Println("Running migrations down...")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error running down: %v", err)
	}
	logrus.Println("Migrations down completed successfully.")
}

func runDownSteps(m *migrate.Migrate, logrus *logrus.Logger, steps int) {
	logrus.Printf("Rolling back %d migration(s)...", steps)
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error rolling back %d steps: %v", steps, err)
	}
	logrus.Printf("Rolled back %d migration(s) successfully.", steps)
}

func runVersion(m *migrate.Migrate, logrus *logrus.Logger) {
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logrus.Fatalf("Error getting version: %v", err)
	}
	if err == migrate.ErrNilVersion {
		logrus.Println("No migrations have been run yet.")
		return
	}
	logrus.Printf("Current version: %d | Dirty: %v", version, dirty)
}

func runForce(m *migrate.Migrate, logrus *logrus.Logger, version int) {
	logrus.Printf("Forcing migration to version %d...", version)
	if err := m.Force(version); err != nil {
		logrus.Fatalf("Error forcing version: %v", err)
	}
	logrus.Printf("Forced to version %d successfully.", version)
}

func runDrop(m *migrate.Migrate, logrus *logrus.Logger) {
	logrus.Println("Dropping all migrations...")
	if err := m.Drop(); err != nil {
		logrus.Fatalf("Error dropping migrations: %v", err)
	}
	logrus.Println("All migrations dropped successfully.")
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up               Apply all pending migrations")
	fmt.Println("  down [steps]     Revert all or N applied migrations (e.g. down 2)")
	fmt.Println("  version          Show current migration version")
	fmt.Println("  force <version>  Force set migration version (use after dirty state)")
	fmt.Println("  drop             Drop all migrations")
}
