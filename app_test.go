package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M){
	err := a.Initialise(DbUser, DbPassword, "test")
	if err!= nil{
		log.Fatal("Error occured while initializing the database")
	}
	createTable()
	m.Run()

}

func createTable(){
	createTableQuery := `create table if not exists products(
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY (id)
	);`
	
	//make a check if any error occur while creating the table
	_, err := a.DB.Exec(createTableQuery);
	if err!= nil{
		log.Fatal(err)
	}
}

//before testing we need to clear the content of the table
func clearTable(){
	a.DB.Exec("DELETE from products")
	a.DB.Exec("alter table products auto_increment=1")
	log.Println("clear table")
}

//adding a product to the table 
func addProduct(name string, quantity int, price float64){
	query := fmt.Sprintf("insert into products(name, quantity, price) VALUES('%v',%v,%v)", name, quantity, price)
	_ , err := a.DB.Exec(query)
	if err!= nil{
		log.Println(err)
	}
}


//first we are going to clear table then we will add one product and then we will make an api call
//to fetch the product and see if the responce is correct or not 
func TestGetProduct(t *testing.T){
	clearTable()
	addProduct("keyboard", 100, 100)
	//httptest package
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int){
	if expectedStatusCode != actualStatusCode{
		t.Errorf("Expected status: %v, Recieved: %v", expectedStatusCode, actualStatusCode)
	}
}


func sendRequest(request *http.Request) *httptest.ResponseRecorder{
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}




func TestCreateProduct(t *testing.T) {
	clearTable()
	//sending the payload in the form of bytes buffer to "/product" in the way as app/json states 
	var product = []byte(`{"name":"chair", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}                //created the map m which takes input as string but can have output as anything
	json.Unmarshal(response.Body.Bytes(), &m)   //stores everything in the response body into the map m

	if m["name"] != "chair"{
		t.Errorf("Expected name: %v, Got: %v", "chair" ,m["name"] )
	}
	if m["quantity"]!=1.0{
		t.Errorf("Expected quantity: %v, Got : %v", 1.0 , m["quantity"])
	}
}



//first we are going to clear the table and then add somthing to the table and use GET api to see if we are able to 
//fetch the product or not and then we are going to execute DELETE api and again try and get that product
func TestDeleteProduct(t *testing.T){
	clearTable()
	addProduct("connector", 10, 10)
	//getting the product
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
	//deleting the product
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
	//if we again get the product it should not send the status ok
	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code) //StatusOk is 200 and StatusNotFound is 404

}



//first we are going to add somthing to the table and then use GET api to fetch that product and then we are going to
//update that product using out PUT endpoint, we are again going to save that response and then we are going to compare
//those two against each other
// func TestUpdatProduct(t *testing.T){
// 	clearTable()
// 	addProduct("connector", 10, 10)
// 	//using GET to fetch product
// 	req, _ := http.NewRequest("GET", "/product/1", nil)
// 	response := sendRequest(req)	
// 	//storing it 
// 	var oldValue map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &oldValue)
// 	//using PUT and changing values
// 	var product = []byte(`{"name":"connector", "quantity":1, "price":10}`)
// 	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
// 	req.Header.Set("Content-Type", "application/json")
	
// 	// panic: runtime error: invalid memory address or nil pointer dereference [recovered]
// 	response = sendRequest(req)
// 	//storing it again
// 	var newValue map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &newValue)
// 	//comparing both

// 	if oldValue["id"]!= newValue["id"]{
// 		t.Errorf("Expected id: %v, Got: %v",newValue["id"], oldValue["id"])
// 	}

// 	if oldValue["name"]!= newValue["name"]{
// 		t.Errorf("Expected id: %v, Got: %v",newValue["name"], oldValue["name"])
// 	}

// 	if oldValue["quantity"]!= newValue["quantity"]{
// 		t.Errorf("Expected id: %v, Got: %v",newValue["quantity"], oldValue["quantity"])
// 	}

// 	if oldValue["price"]!= newValue["price"]{
// 		t.Errorf("Expected id: %v, Got: %v",newValue["price"], oldValue["price"])
// 	}
// }
