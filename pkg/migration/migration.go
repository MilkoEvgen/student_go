package migration

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"path/filepath"
	"strings"
	"student_go/internal/config"
)

func ApplyMigrations() {
	conf := config.Config.DB
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name, conf.SSLMode)

	projectRoot, _ := os.Getwd()
	for !strings.HasSuffix(projectRoot, "student_go") && projectRoot != filepath.Dir(projectRoot) {
		projectRoot = filepath.Dir(projectRoot)
	}

	path := filepath.Join(projectRoot, "migrations")
	path = strings.ReplaceAll(path, `\`, `/`)

	m, err := migrate.New("file://"+path, dsn)
	if err != nil {
		log.Fatalf("Could not initialize migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
