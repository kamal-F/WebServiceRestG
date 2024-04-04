package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"strings"
	"time"

	"rongsokapi/database"
	"rongsokapi/others"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

var jwtKey = []byte("your_secret_key")

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {

	defer database.DB.Close()
	db := database.DB
	http.HandleFunc("/news", others.NewsHandler)

	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Here you should add your database logic to check user credentials
		// For simplicity, we assume the user is authenticated if username and password are not empty
		if user.Username != "" && user.Password != "" {
			tokenString, err := GenerateJWT()
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			idStr := r.URL.Query().Get("id")
			if idStr != "" {
				getProduct(db, w, r)
			} else {
				getProducts(db, w, r)
			}
		case "POST":
			bearerToken := r.Header.Get("Authorization")
			strArr := strings.Split(bearerToken, " ")
			if len(strArr) == 2 {
				isValid, _ := ValidateToken(strArr[1])
				if isValid {
					createProduct(db, w, r)
				} else {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				}
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "PUT":
			updateProduct(db, w, r)
		case "DELETE":
			deleteProduct(db, w, r)
		default:
			http.Error(w, "Unsupported HTTP Method", http.StatusBadRequest)
		}
	})

	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func getProduct(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	//panic("unimplemented")
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id)

	var p Product
	if err := row.Scan(&p.ID, &p.Name, &p.Price); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func getProducts(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", p.Name, p.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProduct(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", p.Name, p.Price, p.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteProduct(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if _, err := db.Exec("DELETE FROM products WHERE id = ?", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
