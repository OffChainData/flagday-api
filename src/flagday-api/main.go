package main

import (
	"fmt"
	"flagday"
	"net/http"
	"log"
	"flag"
	"strconv"
	"encoding/json"
	"time"
)

type Holiday struct {
	Date time.Time `json:"date"` 
	Name string `json:"name"`
	Type flagday.HolidayKind `json:"type"`
}

type RequestEnv struct {}

type RequestHandler struct {
	*RequestEnv
	Handler func(p *RequestEnv, w http.ResponseWriter, r *http.Request) error
}

type Error struct {
	Error string `json:"error"` 
}

func main() {
	port := flag.String("port", "3001", "Port number to server the adapter on")
	flag.Parse()

	p := &RequestEnv{}
	http.Handle("/", RequestHandler{p, ProcessRequest});
	log.Println(fmt.Sprintf("Listening on port %s", *port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

/**
* Return all holidays for the specified years
*/
func ProcessRequest(p *RequestEnv, w http.ResponseWriter, r *http.Request) error {

	startYear, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		return err
	}
	endYear, err := strconv.Atoi(r.URL.Query().Get("end"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//Loop through requested years and return all holidays
	holidays := []Holiday{}
	for i := startYear; i <= endYear; i++ {
		dates := flagday.InYear(i)
		for _, td := range dates {
			n := Holiday{
				Name: td.Name(), 
				Date: td.Time(), 
				Type: td.Kind(), 
			}
			holidays = append(holidays, n)
		}	
	}

	//Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(holidays)
    if err != nil {
		return err
    }
	fmt.Fprint(w, string(b))
	return nil
}

func (fh RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fh.Handler(fh.RequestEnv, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, err := json.Marshal(Error{Error: err.Error()})
		if err != nil {
			log.Printf("Error %s", err)
		} else {
			log.Printf("Response %s", b)
			fmt.Fprint(w, string(b))
		}
	}
}