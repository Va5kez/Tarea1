package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type Address struct {
	Origen  string
	Destino string
}

type Restaurants struct {
	Location string
}

type Direccion struct {
	Lat float64 "json: \"lat\""
	Lon float64 "json: \"lon\""
}

func main() {
	http.HandleFunc("/ejercicio1", handler1)
	http.HandleFunc("/ejercicio2", handler2)
	http.ListenAndServe(":8080", nil)
}

func handler1(w http.ResponseWriter, r *http.Request) {
	var dir Address
	if r.Body == nil {
		http.Error(w, "Porfavor ", 400)
		return
	}
	er := json.NewDecoder(r.Body).Decode(&dir)
	if er != nil {
		http.Error(w, er.Error(), 400)
		return
	}
	fmt.Println(dir)
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyAzzrnc71pLvEvOdY322DQwwbUsFQZT7Vg"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	direccion := regexp.MustCompile("[a-zA-Z\\d ]*, [a-zA-Z\\d ]*")
	dir1 := direccion.FindString(dir.Origen)
	dir2 := direccion.FindString(dir.Destino)
	fmt.Println(direccion)
	a := &maps.DirectionsRequest{
		Origin:      dir1,
		Destination: dir2,
	}
	resp, _, b := c.Directions(context.Background(), a)
	if b != nil {
		log.Fatalf("Error: %s", b)
	}
	fmt.Println(resp)
	buffer := new(bytes.Buffer)
	buffer.WriteString("{\"ruta\":[")
	json.NewDecoder(r.Body).Decode(&resp)
	for x := 0; x < len(resp[0].Legs[0].Steps); x++ {
		buffer.WriteString("{\"lat\":")
		buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lat, 'f', 5, 64))
		buffer.WriteString(", ")
		buffer.WriteString("\"lon\":")
		buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lng, 'f', 5, 64))
		buffer.WriteString("}, ")
		if x == (len(resp[0].Legs[0].Steps) - 1) {
			buffer.WriteString("{\"lat\":")
			buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].EndLocation.Lat, 'f', 5, 64))
			buffer.WriteString(", ")
			buffer.WriteString("\"lon\":")
			buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].EndLocation.Lng, 'f', 5, 64))
			buffer.WriteString("} ")
		}
	}
	buffer.WriteString("]}")
	fmt.Println(buffer.String())
	fmt.Fprintf(w, buffer.String())
}

func handler2(w http.ResponseWriter, r *http.Request) {
	var cercanos Restaurants
	if r.Body == nil {
		http.Error(w, "No se especifico origen", 400)
		return
	}
	json.NewDecoder(r.Body).Decode(&cercanos)
	gg, err := maps.NewClient(maps.WithAPIKey("AIzaSyDrNPltnuTgkKRVIiiAoHrunzjBIPDqDvY"))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	d := &maps.DirectionsRequest{
		Origin:      cercanos.Location,
		Destination: cercanos.Location,
	}
	resp, _, b := gg.Directions(context.Background(), d)
	if b != nil {
		log.Fatalf("Error: %s", b)
	}
	json.NewDecoder(r.Body).Decode(&resp)
	fmt.Println(resp)
	c, _ = maps.NewClient(maps.WithAPIKey("AIzaSyDrNPltnuTgkKRVIiiAoHrunzjBIPDqDvY"))
	t := &maps.NearbySearchRequest{
		Location: &maps.LatLng{resp[0].Legs[0].Steps[0].StartLocation.Lat, resp[0].Legs[0].Steps[0].StartLocation.Lng},
		Radius:   800,
		Type:     "restaurant",
	}
	response, _ := c.NearbySearch(context.Background(), t)
	json.NewDecoder(r.Body).Decode(&response)
	buffer := new(bytes.Buffer)
	buffer.WriteString("{\"Restaurantes\":[")
	for x := 0; x < len(response.Results); x++ {
		buffer.WriteString("{\"Nombre\":\"")
		buffer.WriteString(response.Results[x].Name)
		buffer.WriteString("\", ")
		buffer.WriteString("\"lat\":")
		buffer.WriteString(strconv.FormatFloat(response.Results[x].Geometry.Location.Lat, 'f', 5, 64))
		buffer.WriteString(", ")
		buffer.WriteString("\"lon\":")
		buffer.WriteString(strconv.FormatFloat(response.Results[x].Geometry.Location.Lng, 'f', 5, 64))
		if x == (len(response.Results) - 1) {
			buffer.WriteString("}")
		} else {
			buffer.WriteString("}, ")
		}
	}
	buffer.WriteString("]}")
	fmt.Fprintf(w, buffer.String())
}
