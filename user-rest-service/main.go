package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ID       string `json:"id"                 db:"id"       validate:"required,min=1,max=10,excludesall= "`
	Name     string `json:"name"               db:"name"     validate:"required,min=1,max=50,excludesall= "`
	Email    string `json:"email"              db:"email"    validate:"required,email,min=5,max=50,excludesall= "`
	Password string `json:"password,omitempty" db:"password" validate:"required,min=8,max=50,excludesall= "`
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
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run() error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	h := NewSqlHandler(db)
	router := mux.NewRouter()
	router.HandleFunc("/user", h.SignUp).Methods("POST")
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}
	return nil
}

func NewSqlHandler(db *sqlx.DB) *sqlHandler {
	return &sqlHandler{db: db}
}

func InitDB() (*sqlx.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func UserValidate(user *User) *ErrorMsg {
	var errorMsg ErrorMsg
	validate := validator.New()
	err := validate.Struct(user)
	if err == nil {
		return nil
	}
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

func checkForUniqueID(h *sqlHandler, user *User) (*ErrorMsg, error) {
	var errorMsg ErrorMsg
	var dbID string
	if err := h.db.QueryRowx("SELECT id FROM users WHERE id = ?", user.ID).Scan(&dbID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
	}
	errorMsg.ID = "このユーザーIDは登録できません"
	return &errorMsg, nil
}

func responseByJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *sqlHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if errorMsg := UserValidate(&user); errorMsg != nil {
		responseByJson(w, http.StatusBadRequest, errorMsg)
		return
	}
	errorMsg, err := checkForUniqueID(h, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if errorMsg != nil {
		responseByJson(w, http.StatusConflict, errorMsg)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)
	if _, err := h.db.Exec("INSERT INTO users(id, name, email, password) VALUES(?,?,?,?)", user.ID, user.Name, user.Email, user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = ""
	responseByJson(w, http.StatusOK, user)
}
