package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Task struct {
	Name  string   `json:"name"`
	Date  string   `json:"date"`
	Tasks []string `json:"tasks"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping() // Check if the connection to the database works
	if err != nil {
		log.Fatal("Database connection failed:", err)
	} else {
		log.Println("Successfully connected to the database.")
	}

	// Serve the to_do.html page at the root path (/)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the to_do.html page as the landing page
		http.ServeFile(w, r, "./static/to_do.html")
	})

	// Serve static files (e.g., CSS, JS) from the /static/ route
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			var task Task
			err := json.NewDecoder(r.Body).Decode(&task)
			defer r.Body.Close()

			// Check if decoding the body failed
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Log received task for debugging
			log.Printf("Received task: %+v\n", task)

			subtasksStr := strings.Join(task.Tasks, ",")

			// Insert the main task along with subtasks as a string
			_, err = db.Exec("INSERT INTO tasks (name, date, task) VALUES ($1, $2, $3)",
				task.Name, task.Date, subtasksStr)

			// Check if there was an error during insertion
			if err != nil {
				log.Printf("Error inserting task: %v", err)
				http.Error(w, "Failed to insert task", http.StatusInternalServerError)
				return
			}

			// Return the task with a 200 OK status after insertion
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)

		}
	})

	// Handler for fetching tasks
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handling OPTIONS request separately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodGet {
			name := r.URL.Query().Get("name")
			date := r.URL.Query().Get("date")

			// Log query parameters for debugging
			log.Printf("Fetching tasks for name: %s, date: %s\n", name, date)

			if name == "" || date == "" {
				http.Error(w, "Missing name or date", http.StatusBadRequest)
				return
			}

			rows, err := db.Query("SELECT task FROM tasks WHERE name = $1 AND date = $2", name, date)
			if err != nil {
				log.Println("Error executing query:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var tasks []string
			for rows.Next() {
				var task string
				if err := rows.Scan(&task); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Split the tasks by commas and append to the tasks slice
				subtasks := strings.Split(task, ",")
				tasks = append(tasks, subtasks...)
			}

			// Check for errors after the loop
			if err := rows.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Log received task for debugging
			log.Printf("Fetched task: %+v\n", tasks)

			// Return the fetched tasks as JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	//handler to delete
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodDelete {
			name := r.URL.Query().Get("name")
			date := r.URL.Query().Get("date")

			if name == "" || date == "" {
				http.Error(w, "Missing name or date", http.StatusBadRequest)
				return
			}

			// Perform deletion
			_, err := db.Exec("DELETE FROM tasks WHERE name = $1 AND date = $2", name, date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Task deleted successfully"))
			log.Printf("Deleted task: Name=%s, Date=%s\n", name, date)

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	// Handler for updating tasks
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPut {
			var task Task
			err := json.NewDecoder(r.Body).Decode(&task)
			defer r.Body.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Log received task for debugging
			log.Printf("Updated task: %+v\n", task)

			subtasksStr := strings.Join(task.Tasks, ",")

			// Update the task in the database
			_, err = db.Exec("UPDATE tasks SET task = $1 WHERE name = $2 AND date = $3", subtasksStr, task.Name, task.Date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func getConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
}
