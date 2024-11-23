package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func main() {
	router := gin.Default()
	router.POST("/borrow", borrowBook)
	router.GET("/getallbooksread/:user_id", getAllBooksByUserID)
	router.POST("/return", returnBook)
	router.Run("localhost:8080")
}

// borrowBook will add the given book info (JSON body) to the borrowed books under the given user
func borrowBook(c *gin.Context) {
	var newLoanedItem borrowedentity

	if err := c.BindJSON(&newLoanedItem); err != nil {
		return
	}

	// Successful bind, we have the information from the user

	c.IndentedJSON(http.StatusOK, newLoanedItem)
}

// returnBook will add the given book info (JSON body) to the returned books under the given user
func returnBook(c *gin.Context) {
	var newReturnedItem borrowedentity

	if err := c.BindJSON(&newReturnedItem); err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, newReturnedItem)
}

// getAllBooksByUserID returns all the books the given user has read since using the system
func getAllBooksByUserID(c *gin.Context) {
	user_id := c.Param("user_id")

	var itemsreadbyuser borrowedentity

	// Perform basic sanity checks here
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
