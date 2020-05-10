package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
	"runtime"
)

func ExecuteMigration(db *sql.DB) {
	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	//executeCmd()

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")
}

func executeCmd() {
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {
		// here we perform the pwd command.
		// we can store the output of this in our out variable
		// and catch any errors in err
		out, err := exec.Command("ls").Output()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			fmt.Printf("%s", err)
		}
		// as the out variable defined above is of type []byte we need to convert
		// this to a string or else we will see garbage printed out in our console
		// this is how we convert it to a string
		fmt.Println("Command Successfully Executed")
		output := string(out[:])
		fmt.Println(output)

		// let's try the pwd command herer
		out, err = exec.Command("ls").Output()
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Println("Command Successfully Executed")
		output = string(out[:])
		fmt.Println(output)
	}
}
