package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type ProductInfo struct {
	Id       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
	Price    float64 `json:"price,omitempty"`
}

func GetsqlConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/products?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetProducts(w http.ResponseWriter, r *http.Request) {

	db = GetsqlConnection()
	defer db.Close()

	s := ProductInfo{}
	ss := []ProductInfo{}
	rows, err := db.Query("select * from product")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&s.Id, &s.Name, &s.Quantity, &s.Price)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
	}

}

func addProducts(w http.ResponseWriter, r *http.Request) {

	db = GetsqlConnection()
	defer db.Close()

	s := ProductInfo{}
	json.NewDecoder(r.Body).Decode(&s)
	result, err := db.Exec("insert into product(id,name,quantity,price) values(?,?,?,?)", s.Id, s.Name, s.Quantity, s.Price)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode("{error: orecord not inserted}")
		} else {
			json.NewEncoder(w).Encode(s)
		}

	}

}
func updateProducts(w http.ResponseWriter, r *http.Request) {
	db = GetsqlConnection()
	defer db.Close()

	s := ProductInfo{}
	json.NewDecoder(r.Body).Decode(&s)

	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("update product set name=?, Quantity=?, price=? where id = ?", s.Name, s.Quantity, s.Price, id)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{error: no record is updated}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}

}
func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	db = GetsqlConnection()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("delete from product where id = ?", id)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{error: no record is updated}")
		} else {
			json.NewEncoder(w).Encode(rowsAffected)
		}
	}
}

func main() {
	Router := mux.NewRouter()
	Router.HandleFunc("/Product", GetProducts).Methods("GET")
	Router.HandleFunc("/Product/add", addProducts).Methods("POST")
	Router.HandleFunc("/Product/{id}", updateProducts).Methods("PUT")
	Router.HandleFunc("/Product/{id}", DeleteProducts).Methods("DELETE")
	http.ListenAndServe(":8080", Router)
}
