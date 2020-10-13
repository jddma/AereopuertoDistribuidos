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
}

//Método que el avión usara remotamente para solicitar permiso para aterrizar
func (t *AirTraffic) RequestLanding(avion *client.Plane, response *string) error {

	*response = "Aprobado"

	return nil

}

func (t *AirTraffic) ValidatePlane(plane *client.Plane, response *bool) error {

	for _, airport := range t.airports{
		if airport.Name == plane.CurrentAirport {
			*response = true
			return nil
		}
	}
	*response = false
	return nil

}

func (t *AirTraffic) getAirports() {

	airportsStr, err := ioutil.ReadFile("resources/aiportsInput.json")
	if err != nil{
		log.Fatal("Error al obtener los aerepuertos: ", err)
	}
	json.Unmarshal(airportsStr, &t.airports)

}

func (t *AirTraffic) StartServer(port string) {

	t.getAirports()

	rpc.Register(t)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	fmt.Println("Escuchando el puerto ", port)
	http.Serve(l, nil)

}
