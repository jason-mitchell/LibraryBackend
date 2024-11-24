package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	_ "LibrarySystem/docs"
	_ "github.com/lib/pq"
)

const (
	host     = "192.168.0.60"
	port     = 5432
	user     = "postgres"
	password = "dbpasswd"
	dbname   = "elibrarytest"
)

type userdetails struct {
	UserID    int
	FirstName string
	LastName  string
}

type Users struct {
	users []userdetails
}

type book struct {
	ID     string `json:"ID"`
	Title  string `json:"Title"`
	Author string `json:"Author"`
}

type borrowedentity struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Books     []book `json:"books"`
}

type titlesInLibrary struct {
	Books []book `json:"books"`
}

func main() {

	getUsers() // test to check DB connectivity

	router := gin.Default()
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.POST("/borrow", borrowBook)
	router.GET("/getallbooksread/:user_id", getAllBooksByUserID)
	router.POST("/return", returnBook)
	router.GET("/getallbooksavailable", getaAllBooksInLibrary)
	router.Run("localhost:8080")
}

// This function is what I used to test that I could reach the PostgreSQL DB and send a very simple query to make sure
// everything works DB side, and I am not going down an inadvertent rabbit hole

func getUsers() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// DB connection OK... send a simple query to get all the entries in the "users" table
	sqlStatement := `SELECT * FROM public."user" ORDER BY user_id ASC `

	rows, _ := db.Query(sqlStatement) // db.Query versus db.QueryRow

	// One or more rows will be returned, scan through all of them...
	for rows.Next() {
		var id int
		var firstName string
		var lastName string
		err = rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			// error handling
			panic(err)
		}

		// Print out all the rows to confirm what we see in pgAdmin is what we see here
		fmt.Println(id, firstName, lastName)
	}
	// Handle any errors during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

}

// -----------------------------------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------------------------------
func borrowBook(c *gin.Context) {
	var newLoanedItem borrowedentity

	if err := c.BindJSON(&newLoanedItem); err != nil {
		return
	}

	// Perform write to DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unable to contact database"})
	}
	defer db.Close()

	// DB is fine

	c.IndentedJSON(http.StatusOK, newLoanedItem)
}

//------------------------------------------------------------------
// getaAllBooksInLibrary will return a list of all available titles
//------------------------------------------------------------------

func getaAllBooksInLibrary(c *gin.Context) {
	var titlesAvailable titlesInLibrary

	// Connect to PostgreSQL DB and get all the books available
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unable to contact database"})
	}
	defer db.Close()

	// DB connection OK... send a simple query to get all the entries in the "users" table
	sqlStatement := `SELECT * FROM public.books ORDER BY id ASC `

	rows, _ := db.Query(sqlStatement) // db.Query versus db.QueryRow

	// One or more rows will be returned, scan through all of them...
	for rows.Next() {
		var id int
		var Title string
		var Author string
		err = rows.Scan(&id, &Title, &Author)
		if err != nil {
			// If something goes wrong with the DB query, return an HTTP error
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "error retrieving records"})
		}

		// Insert each book record into the struct we will return
		//fmt.Println(id, Title, Author)
		titlesAvailable.Books = append(titlesAvailable.Books, book{
			ID:     fmt.Sprintf("%d", id),
			Title:  Title,
			Author: Author,
		})
	}
	// Handle any errors during iteration
	err = rows.Err()
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "error retrieving records"})
	}

	// Everything went well, return the records
	c.IndentedJSON(http.StatusOK, titlesAvailable)

}

// -------------------------------------------------------------------------------------------------
// returnBook will add the given book info (JSON body) to the returned books under the given user
// --------------------------------------------------------------------------------------------------
func returnBook(c *gin.Context) {
	var newReturnedItem borrowedentity

	if err := c.BindJSON(&newReturnedItem); err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, newReturnedItem)
}

// ------------------------------------------------------------------------------------------
// getAllBooksByUserID returns all the books the given user has read since using the system
// ------------------------------------------------------------------------------------------
func getAllBooksByUserID(c *gin.Context) {
	user_id := c.Param("user_id")

	var itemsreadbyuser borrowedentity

	// Static data used to return canned responses for checking the API
	if user_id != "ABCDEF" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "records not found"})
	} else {
		itemsreadbyuser.UserID = "ABCDEF"
		itemsreadbyuser.FirstName = "Joe"
		itemsreadbyuser.LastName = "Bloggs"

		// Let's say our database returns only 2 items, we need to then make sure we can assign
		itemsreadbyuser.Books = append(itemsreadbyuser.Books, book{
			ID:     "8-99787-01",
			Title:  "A Brief History of Time",
			Author: "Stephen Hawking",
		})
		c.IndentedJSON(http.StatusOK, itemsreadbyuser)
	}

}
