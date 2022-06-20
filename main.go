package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID    string `json: "id"`
	Isbn  string `json: "isbn"`
	Title string `json: "title"`
	//* pointer because the Director struct will be associated with  a Movie struct
	//Director  will have  the same valued defined inside the  director
	Director *Director `json: "director"`
}

type Director struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname`
}

//slice of the type Movie
var movies []Movie

//passing a pointer of the request sent from Postman
//w
func getMovies(w http.ResponseWriter, r *http.Request) {
	//set content type as json for the struct
	w.Header().Set("Content-Type", "application/json")
	//encode the response into json
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//id from the postman as part of the request. Range = forEach
	params := mux.Vars(r)

	for index, item := range movies {
		//accesing the id from the request params sent by me
		if item.ID == params["id"] {
			//whatever id is matching to delete, is  goig to take that index. And directly  append all the other Data it in that place
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//mux creates Routes. Vars allows  us to get the params in Request
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	//sending the entire movie through the body
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	//set json content type
	w.Header().Set("Content-type", "application/json")
	//access to params
	params := mux.Vars(r)
	//loop over the movies
	//delete the movie with the ID  that is sent
	//add a new movie - the movie that  is sent in the body of postman

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}

}

func main() {
	r := mux.NewRouter()

	//append to the above variable
	//&Director because we want the reference object address
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server  at port 8000\n")
	//starting server
	log.Fatal(http.ListenAndServe(":8000", r))
}
