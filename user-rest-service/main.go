package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"       db:"id"       validate:"required,min=1,max=10,excludesall= "`
	Name     string `json:"name"     db:"name"     validate:"required,min=1,max=50,excludesall= "`
	Email    string `json:"email"    db:"email"    validate:"required,email,min=5,max=50,excludesall= "`
	Password string `json:"password" db:"password" validate:"required,min=8,max=50,excludesall= "`
}

type ErrorMsg struct {
	ID       string `json:"error_id"`
	Name     string `json:"error_name"`
	Email    string `json:"error_email"`
	Password string `json:"error_password"`
}

type sqlHandler struct {
	db *sqlx.DB
}

func main() {
	db := InitDB()
	h := NewSqlHandler(db)

	router := mux.NewRouter()
	router.HandleFunc("/user", h.SignUp).Methods("POST")
	log.Print(http.ListenAndServe(":8080", router))
}

func NewSqlHandler(db *sqlx.DB) *sqlHandler {
	return &sqlHandler{db: db}
}

func InitDB() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load the .env: %v", err)
	}
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Printf("failed to open the DB: %v", err)
	}
	return db
}

func UserValidate(user *User) *ErrorMsg {
	var errorMsg ErrorMsg
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			switch fieldName {
			case "ID":
				errorMsg.ID = "ユーザーIDが正しくありません"
			case "Name":
				errorMsg.Name = "ユーザーネームが正しくありません"
			case "Email":
				errorMsg.Email = "ユーザーメールが正しくありません"
			case "Password":
				errorMsg.Password = "パスワードが正しくありません"
			}
		}
		return &errorMsg
	}
	return nil
}

func checkForUniqueID(h *sqlHandler, user *User) *ErrorMsg {
	var errorMsg ErrorMsg
	var dbID string
	if err := h.db.QueryRowx("SELECT id FROM users WHERE id = ?", user.ID).Scan(&dbID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Printf("failed to scan the row: %v", err)
		}
	}
	errorMsg.ID = "このユーザーIDは登録できません"
	return &errorMsg
}

func responseByJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *sqlHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("failed to Decode the user: %v", err)
	}
	if errorMsg := UserValidate(&user); errorMsg != nil {
		responseByJson(w, http.StatusBadRequest, errorMsg)
		return
	}
	if errorMsg := checkForUniqueID(h, &user); errorMsg != nil {
		responseByJson(w, http.StatusConflict, errorMsg)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)
	if _, err := h.db.Exec("INSERT INTO users(id, name, email, password) VALUES(?,?,?,?)", user.ID, user.Name, user.Email, user.Password); err != nil {
		log.Printf("failed to insert the users: %v", err)
	}
	user.Password = ""
	responseByJson(w, http.StatusOK, user)
}
