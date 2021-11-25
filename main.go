package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// For listen port (mandatory)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env is required")
	}

	instanceID := os.Getenv("INSTANCE_ID")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Response handler
		if r.Method != "GET" {
			http.Error(w, "the requested http method is not allowed", http.StatusMethodNotAllowed)
		}
		// Response body
		text := "hello world! tag: latest"
		if instanceID != "" {
			text = text + ". from " + instanceID
		}
		w.Write([]byte(text))
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllUsersHandler(w, r)
		case "POST":
			createUserHandler(w, r)
		default:
			http.Error(w, "http method not allowed", http.StatusBadRequest)
		}
	})
	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":" + port
	log.Println("web server is starting at", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := conn()
	if err != nil {
		// to write if there are an errors
		writreError(w, err)
		return
	}
	defer conn.Close()

	qry, err := conn.QueryContext(context.Background(), "SELECT * FROM users")
	if err != nil {
		// to write if there are an errors
		writreError(w, err)
		return
	}
	result := make([]User, 0)
	for qry.Next() {
		var id sql.NullInt32
		var firstName sql.NullString
		var lastName sql.NullString
		var birth sql.NullTime
		err = qry.Scan(&id, &firstName, &lastName, &birth)

		if err != nil {
			writreError(w, err)
			return
		}

		user := User{}
		user.ID = int(id.Int32)
		user.FirstName = firstName.String
		user.LastName = lastName.String
		user.Birth = birth.Time
		result = append(result, user)
	}

	writeData(w, result)

}
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := new(User)

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		writreError(w, err)
		return
	}

	conn, err := conn()
	if err != nil {
		writreError(w, err)
		return
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(context.Background(), "INSERT INTO users(first_name, last_name,birth) VALUES (?,?,?)")
	if err != nil {
		writreError(w, err)
		return
	}

	// Insert value to finding variable (?,?,?)
	stmtRes, err := stmt.ExecContext(context.Background(), payload.FirstName, payload.LastName, payload.Birth)
	if err != nil {
		writreError(w, err)
		return
	}

	id, _ := stmtRes.LastInsertId()
	result := map[string]interface{}{"lastInsertID": id}

	writeData(w, result)
}

// To run (for windows)
// set MYSQL_CONN_STRING=root@tcp(localhost:3306)/database_name?parseTime=true
// set PORT=8083
// go run main.go lib.go model.go

// (for linux)
// export MYSQL_CONN_STRING=root@tcp(localhost:3306)/database_name?parseTime=true
// PORT=8083 go run main.go lib.go model.go
