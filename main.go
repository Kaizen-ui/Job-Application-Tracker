package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable(db)

	for {
		fmt.Println("========================================")
		fmt.Println("|| Welcome to Job Application Tracker ||")
		fmt.Println("========================================")
		fmt.Println("1. Display Applications\n2. Search Applications\n3. Add Application\n4. Update Application\n5. Delete Application\n6. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			displayapplication(db)
		case 2:
			searchapplication(db)
		case 3:
			addapplication(db)
		case 4:
			updateapplication(db)
		case 5:
			deleteapplication(db)
		case 6:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		var continueChoice string
		fmt.Print("Do you want to perform another operation? (y/n): ")
		fmt.Scanln(&continueChoice)

		continueChoice = strings.ToLower(strings.TrimSpace(continueChoice))

		if continueChoice != "y" && continueChoice != "yes" {
			fmt.Println("Exiting...")
			break
		}
	}
}

func displayapplication(db *sql.DB) {
	query := "SELECT id, company_name, job_title, description, salary, location, applied FROM application"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var companyName, jobTitle, description, location string
		var salary float64
		var applied string

		err := rows.Scan(&id, &companyName, &jobTitle, &description, &salary, &location, &applied)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, Company: %s, Job Title: %s, Description: %s, Salary: %.2f, Location: %s, Applied: %s\n",
			id, companyName, jobTitle, description, salary, location, applied)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func searchapplication(db *sql.DB) {
	var searchTerm string
	fmt.Print("Enter search term (company name/job title): ")
	fmt.Scanln(&searchTerm)

	query := "SELECT id, company_name, job_title, description, salary, location, applied FROM application WHERE company_name ILIKE $1 OR job_title ILIKE $1"
	rows, err := db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Flag to check if any rows were found
	found := false

	for rows.Next() {
		var id int
		var companyName, jobTitle, description, location string
		var salary float64
		var applied string

		err := rows.Scan(&id, &companyName, &jobTitle, &description, &salary, &location, &applied)
		if err != nil {
			log.Fatal(err)
		}

		// Display each record
		fmt.Printf("ID: %d, Company: %s, Job Title: %s, Description: %s, Salary: %.2f, Location: %s, Applied: %s\n",
			id, companyName, jobTitle, description, salary, location, applied)

		// Set the flag to true if at least one row is found
		found = true
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// If no rows were found, display a message
	if !found {
		fmt.Println("No matching applications found.")
	}
}

func addapplication(db *sql.DB) {
	var companyName, jobTitle, description, location string
	var salary float64

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter company name: ")
	companyName, _ = reader.ReadString('\n')
	companyName = strings.TrimSpace(companyName)

	fmt.Print("Enter job title: ")
	jobTitle, _ = reader.ReadString('\n')
	jobTitle = strings.TrimSpace(jobTitle)

	fmt.Print("Enter job description: ")
	description, _ = reader.ReadString('\n')
	description = strings.TrimSpace(description)

	for {
		fmt.Print("Enter salary: ")

		salaryInput, _ := reader.ReadString('\n')
		salaryInput = strings.TrimSpace(salaryInput)

		salaryInput = strings.ReplaceAll(salaryInput, ",", "")

		var err error
		salary, err = strconv.ParseFloat(salaryInput, 64)
		if err != nil {

			fmt.Println("Invalid salary input. Please enter a valid number.")
		} else if salary < 0 {

			fmt.Println("Salary cannot be negative. Please enter a valid amount.")
		} else {

			break
		}
	}

	fmt.Print("Enter location: ")
	location, _ = reader.ReadString('\n')
	location = strings.TrimSpace(location)

	query := `INSERT INTO application (company_name, job_title, description, salary, location)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, companyName, jobTitle, description, salary, location)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application added successfully!")
}
func updateapplication(db *sql.DB) {
	var id int
	var companyName, jobTitle, description, location string
	var salary float64

	fmt.Print("Enter the ID of the application you want to update: ")
	fmt.Scanln(&id)
	fmt.Print("Enter new company name: ")
	fmt.Scanln(&companyName)
	fmt.Print("Enter new job title: ")
	fmt.Scanln(&jobTitle)
	fmt.Print("Enter new job description: ")
	fmt.Scanln(&description)
	fmt.Print("Enter new salary: ")
	fmt.Scanln(&salary)
	fmt.Print("Enter new location: ")
	fmt.Scanln(&location)

	query := `UPDATE application SET company_name = $1, job_title = $2, description = $3, salary = $4, location = $5 WHERE id = $6`
	_, err := db.Exec(query, companyName, jobTitle, description, salary, location, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application updated successfully!")
}

func deleteapplication(db *sql.DB) {
	var id int
	fmt.Print("Enter the ID of the application you want to delete: ")
	fmt.Scanln(&id)

	query := `DELETE FROM application WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application deleted successfully!")
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS application (
		id SERIAL PRIMARY KEY,
		company_name VARCHAR (255) NOT NULL,
		job_title VARCHAR (255) NOT NULL,
		description VARCHAR (255) NOT NULL,
		salary Numeric (15,2) NOT NULL,
		location VARCHAR (255) NOT NULL,
		applied timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
