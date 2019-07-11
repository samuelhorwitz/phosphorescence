package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var migrationsNamespace = uuid.NewV5(common.PhosphorUUIDV5Namespace, "migrations")
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func migrate(cfg *config) {
	fmt.Println("Running migrations...")
	db, err := sql.Open("postgres", cfg.postgresConnectionString)
	if err != nil {
		log.Fatalf("Could not initialize Postgres: %s", err)
		return
	}
	defer db.Close()
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get executable path: %s", err)
		return
	}
	exPath := filepath.Dir(ex)
	migrationsDirPath := filepath.Join(exPath, "migrations")
	// ioutil.ReadDir is sorted by filename and our filenames all start with Unix time for ordering
	migrationFiles, err := ioutil.ReadDir(migrationsDirPath)
	if err != nil {
		log.Fatalf("Could not read migrations folder: %s", err)
		return
	}
	var migrationsTableExists bool
	err = db.QueryRow("select exists (select 1 from pg_tables where schemaname = 'public' and tablename = '_migrations');").
		Scan(&migrationsTableExists)
	if err != nil {
		log.Fatalf("Unable to check if database is initialized: %s", err)
		return
	}
	for i, migrationFile := range migrationFiles {
		migrationFilename := migrationFile.Name()
		migration, err := ioutil.ReadFile(filepath.Join(migrationsDirPath, migrationFilename))
		if err != nil {
			log.Fatalf("Cannot read migration file %s", migrationFilename)
			return
		}
		migrationID := uuid.NewV5(migrationsNamespace, string(migration))
		log.Printf("About to run migration %s (%s)", migrationFilename, migrationID.String())
		shouldRunMigration := false
		if i == 0 && !migrationsTableExists {
			shouldRunMigration = true
		} else {
			var foundID uuid.UUID
			var foundName string
			err = psql.Select("id", "name").From("_migrations").
				Where(sq.Or{
					sq.Eq{"id": migrationID},
					sq.Eq{"name": migrationFilename},
				}).RunWith(db).QueryRow().Scan(&foundID, &foundName)
			if err != nil {
				if err == sql.ErrNoRows {
					shouldRunMigration = true
				} else {
					log.Fatalf("Unable to check for existing migration: %s", err)
					return
				}
			}
			if (uuid.Equal(migrationID, foundID) && migrationFilename != foundName) || (migrationFilename == foundName && !uuid.Equal(migrationID, foundID)) {
				log.Fatalf("It appears that a migration was likely altered after running, bailing. (File %s, in db as %s, ID %s, in db as %s)", migrationFilename, foundName, migrationID.String(), foundID.String())
				return
			}
		}
		if shouldRunMigration {
			runInTx := true
			if strings.HasSuffix(migrationFilename, ".no-transaction.sql") {
				runInTx = false
			}
			var tx *sql.Tx
			if runInTx {
				tx, err = db.Begin()
				if err != nil {
					rollbackAndFatal(tx, fmt.Errorf("Could not begin transaction: %s", err))
					return
				}
			}
			_, err = db.Exec(string(migration))
			if err != nil {
				err = fmt.Errorf("Could not run migration: %s", err)
				if runInTx {
					rollbackAndFatal(tx, err)
				} else {
					log.Fatal(err)
				}
				return
			}
			ins := psql.Insert("_migrations").Columns("id", "name").
				Values(migrationID, migrationFilename)
			if runInTx {
				_, err = ins.RunWith(tx).Exec()
			} else {
				_, err = ins.RunWith(db).Exec()
			}
			if err != nil {
				err = fmt.Errorf("Could not log migration in database as having run: %s", err)
				if runInTx {
					rollbackAndFatal(tx, err)
				} else {
					log.Fatal(err)
				}
				return
			}
			if runInTx {
				err = tx.Commit()
				if err != nil {
					log.Fatalf("Could not commit transaction: %s", err)
					return
				}
			}
		} else {
			log.Printf("Skipping %s, already ran...", migrationFilename)
		}
	}
}

func rollbackAndFatal(tx *sql.Tx, err error) {
	log.Fatal(common.TryToRollback(tx, err))
}
