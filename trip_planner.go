package main

import (
"fmt"
"log"
"io/ioutil"
"net/http"
"encoding/json"
//"github.com/jasonwinn/geocoder"
//"math/rand"
//"labix.org/v2/mgo"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"github.com/julienschmidt/httprouter"
//"strconv"
)

type Request struct{
	Starting_from_location_id string `json: "Starting_from_location_id"`
	Location_ids string `json: "Location_ids"`
}

type Response1 struct{
	Id int `json: "Id"`
	Status string `json: Status`
	Starting_from_location_id int `json: "Starting_from_location_id"`
	Best_route_location_ids []int `json:"Best_route_location_ids"`
	Total_uber_costs int `json:"Total_uber_costs"`
	Total_uber_duration int `json: "Total_uber_duration"`
	Total_distance int `json: "Total_distance"`
}


type Coordinates struct{
	Lat float64 `json: "Lat"`
	Lng float64 `json: "Lng"`
}

type Response struct{
	Id int `json: "Id"`
	Name string `json: "Name"`
	Address string `json: "Address"`
	City string `json: "City"`
	State string `json: "State"`
	Zip string `json: "Zip"`
	Coordinate Coordinates `json: "Coordinate"`
}

func CreateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var idata Request
	var result Response
	fmt.Println("Inside CreateLocation")
	session :=connectToDb()

	fmt.Println("Connecting to DB")
	c:= session.DB("vinaysh").C("admin")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("Read Error"+err.Error())
	}
	fmt.Println("Input JSON data", string(body))

	if err:= json.Unmarshal(body, &idata);
	err != nil {
		panic(err)
		fmt.Println("Unmarshal failed")
	}

	id:= idata.Starting_from_location_id
	fmt.Println("ID is", id)

	fmt.Println("Find")
    err1 := c.Find(bson.M{"id":id }).One(&result)
    if err1 != nil {
        rw.WriteHeader(404)
        fmt.Println("ID not found", id)
        return
    }
    fmt.Println("Printing lat lng")
    latitude := result.Coordinate.Lat
    longitude := result.Coordinate.Lng

	fmt.Println("Latitude is ", latitude)
	fmt.Println("Longitude is ", longitude)



}


func connectToDb() *mgo.Session{
    session, err := mgo.Dial("mongodb://admin:admin@ds045054.mongolab.com:45054/vinaysh")
        if err != nil {
                panic("Couldn't connect to the database")
        }
           session.SetMode(mgo.Monotonic, true)
           fmt.Println("Session is ",session)
    return session
            
}


func main(){
	 fmt.Println("Inside Main")
     mux := httprouter.New()
     mux.POST("/trips",CreateLocation) 
    // mux.GET("/trips/:trip_id",ReadLocation)
     //mux.PUT("/trips/:trip_id/request",UpdateLocation)
     server := http.Server{
             Addr:        "127.0.0.1:8080",
             Handler: mux,
     }
     server.ListenAndServe()
}