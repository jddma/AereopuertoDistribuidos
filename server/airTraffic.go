package server

import (
	"../client"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type AirTraffic struct {
	airports []airport
	Routes []client.Routes
}

func (a *AirTraffic) SearchPossibleRoutes(plane *client.Plane, response *[]client.Routes) error {

	posibleRoutes := []client.Routes{}

	for _, route := range a.Routes {
		if plane.CurrentAirport == route.TakeoffAirport{
			posibleRoutes = append(posibleRoutes, route)
		}
	}

	*response = posibleRoutes
	return nil

}

func (a *AirTraffic) RequestPermission(plane *client.Plane, response *bool) error {

	airportIndex := a.searchAirport(plane.CurrentAirport)
	if a.airports[airportIndex].IsFree {
		*response = true
		a.airports[airportIndex].IsFree = false
	}

	return nil

}

func (a *AirTraffic) ConfirmOperation(plane *client.Plane, response *bool) error {

	airportIndex := a.searchAirport(plane.CurrentAirport)
	a.airports[airportIndex].IsFree = true

	return nil

}


func (a *AirTraffic) ValidatePlane(plane *client.Plane, response *bool) error {

	for _, airport := range a.airports{
		if airport.Name == plane.CurrentAirport {
			*response = true
			return nil
		}
	}
	*response = false
	return nil

}

func (a *AirTraffic) StartServer(port string) {

	a.getAirports()
	a.getRoutes()

	rpc.Register(a)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	fmt.Println("Escuchando el puerto ", port)
	http.Serve(l, nil)

}

func (a *AirTraffic) getAirports() {

	airportsStr, err := ioutil.ReadFile("resources/aiportsInput.json")
	if err != nil{
		log.Fatal("Error al obtener los aerepuertos: ", err)
	}
	json.Unmarshal(airportsStr, &a.airports)

}

func (a *AirTraffic) getRoutes() {

	routesStr, err := ioutil.ReadFile("resources/routesInput.json")
	if err != nil{
		log.Fatal("Error al obtener los aerepuertos: ", err)
	}
	json.Unmarshal(routesStr, &a.Routes)

}

func (a *AirTraffic) searchAirport(name string) int {

	result := -1

	for i := 0; i < len(a.airports); i++ {
		if a.airports[i].Name == name{
			return i
		}
	}

	return result

}