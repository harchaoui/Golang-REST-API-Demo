/*
* this demo shows how to implement a full RESTful web-app
* with data Persistence "using postgresSQL as our data-store"
**/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var (
	dbConnection *sql.DB
)

func init() {

	var err error

	//Create the Connection Pool
	dbConnection, err = sql.Open("postgres", "user=gopher password=1111aaaa  dbname=reset_demo sslmode=disable")

	//check for err
	if err != nil {
		log.Fatal(err)
		return
	}

	//Check for the database Connectivity, because the sql.Open() doesn't perform it
	if err = dbConnection.Ping(); err != nil {
		log.Fatal(err)
		return
	} else {
		fmt.Printf("Database Connection Established!\n")
	}

}

type Customer struct {
	uid   string
	name  string
	email string
}

func main() {

	//Create Routes and related actions!
	http.HandleFunc("/customers", customers)
	http.HandleFunc("/customers/show", customersShow) //  -> link like  localhost:3000/customers/show?uid= 29929
	http.HandleFunc("/customers/create", customersCreate)

	//Start The Server
	http.ListenAndServe(":3000", nil)
}

//Handle the third route /customers/create
func customersCreate(w http.ResponseWriter, r *http.Request) {
	//check if the Method is Post
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		log.Fatal(w)
		return
	}

	// get params from the post Request
	uid := r.FormValue("uid")
	name := r.FormValue("name")
	email := r.FormValue("email")

	//check if they are empty
	if uid == "" || name == "" || email == "" {
		//Bad Request ... we need a uid !
		http.Error(w, http.StatusText(400), 400)
		log.Fatal(w)
		return
	}

	//Build the Query String
	query := "INSERT INTO reset_demo.customer VALUES($1,$2,$3)"

	//execute the insert Query
	result, err := dbConnection.Exec(query, uid, name, email)

	//check if any error
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Fatal(err)
		return
	}

	//check how many rows affected (added or modified in case of update!!!)
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Fatal(err)
		return
	}

	// output the result
	fmt.Fprintf(w, " Customer %s created Succesfully (%d rows affected! ) \n", uid, rowsAffected)

}

// handle the second route /customers/show
func customersShow(w http.ResponseWriter, r *http.Request) {
	//check the request Method again
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		log.Fatal(w)
		return
	}

	/// get the UID from the request
	uid := r.FormValue("uid")

	//check if the uid isn't empty
	if uid == "" {
		//Bad Request ... we need a uid !
		http.Error(w, http.StatusText(400), 400)
		log.Fatal(w)
		return
	}

	//build the Query String
	query := "SELECT * FROM reset_demo.customer WHERE uid=$1"

	//Execute the query
	row := dbConnection.QueryRow(query, uid)

	//create a new Customer object to hold data
	ctm := new(Customer)

	//scan data from the row
	err := row.Scan(&ctm.uid, &ctm.name, &ctm.email)

	//check if any error !
	//First we need to check if no row is found !

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		log.Fatal(err)
		return
	} else if err != nil {
		// internal error
		http.Error(w, http.StatusText(500), 500)
		log.Fatal(err)
		return
	}

	// output the result
	fmt.Fprintf(w, "uid:%s name:%s email:%s", ctm.uid, ctm.name, ctm.email)

}

// handle the First route /customers
func customers(w http.ResponseWriter, r *http.Request) {

	// Check the Method sent
	// We are looking for a GET method or the app will die/ stoped
	if r.Method != "GET" {
		//any other method is not allowed !
		http.Error(w, http.StatusText(405), 405)
		//log the error to the console
		log.Fatal(w)
		return

	}

	// create the query string

	query := "SELECT * FROM reset_demo.customer"

	//Execute the Query
	rows, err := dbConnection.Query(query)

	if err != nil {
		//Internal error
		http.Error(w, http.StatusText(500), 500)
		log.Fatal(err)
		return
	}
	defer rows.Close() // Always close opened any resources

	//create a new slice of Customer Object to hold all retrieved data
	ctms := make([]*Customer, 0)

	//Loop Over the rows
	for rows.Next() {
		//create a Customer to hold one row
		ctm := new(Customer)

		//Scan For data and fill the customer fields
		err := rows.Scan(&ctm.uid, &ctm.name, &ctm.email)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			log.Fatal(err)
			return
		}

		// append ctm to ctms slice
		ctms = append(ctms, ctm)

	}
	// check for any errors occured during the retrieving process!!! might happened

	if err := rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Fatal(err)
		return
	}

	//Send data to the output
	fmt.Fprintf(w, "Customers List \n")
	//Loopr Over Ctms to print the customers list

	for _, ctm := range ctms {
		fmt.Fprintf(w, "uid:%s name:%s email:%s", ctm.uid, ctm.name, ctm.email)
	}

}
