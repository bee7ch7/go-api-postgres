package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "x"
	password = "x"
	dbname   = "x"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type NewAccount struct {
	Owner    string `json:"owner" binding: "required"`
	Currency string `json:"currency" binding: "required", oneof=USD UAH EUR`
}

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Account `json:"data"`
	Message string    `json:"message"`
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetAccounts(c *gin.Context) {
	db := setupDB()

	printMessage("Getting accounts...")

	rows, err := db.Query("SELECT * FROM accounts")

	// check errors
	checkErr(err)

	var accounts []Account

	for rows.Next() {
		var id int64
		var owner string
		var balance int64
		var currency string
		var created_at time.Time

		err = rows.Scan(&id, &owner, &balance, &currency, &created_at)

		// check errors
		checkErr(err)

		accounts = append(accounts, Account{ID: id, Owner: owner, Balance: balance, Currency: currency, CreatedAt: created_at})
	}

	var response = JsonResponse{Type: "success", Data: accounts}
	c.IndentedJSON(http.StatusOK, response)

}

func CreateAccount(c *gin.Context) {

	var req NewAccount
	var account []Account

	if err := c.BindJSON(&req); err != nil {
		return
	}

	db := setupDB()

	var id int64
	var owner string
	var balance int64 = 0
	var currency string
	var created_at time.Time

	err := db.QueryRow("INSERT INTO accounts(owner, currency, balance) VALUES($1, $2, $3) returning id, owner, balance, currency, created_at;",
		req.Owner, req.Currency, balance).Scan(&id, &owner, &balance, &currency, &created_at)

	checkErr(err)
	account = append(account, Account{ID: id, Owner: owner, Balance: balance, Currency: currency, CreatedAt: created_at})

	var response = JsonResponse{Type: "success", Data: account}
	c.IndentedJSON(http.StatusOK, response)
}

func main() {

	router := gin.Default()

	router.GET("/accounts", GetAccounts)
	router.POST("/account", CreateAccount)

	fmt.Println("Server at 8080")
	log.Fatal(router.Run("0.0.0.0:8080"))
}
