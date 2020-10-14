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

//Método a servir al cliente para obtener las posibles rutas que este pueda tomar
func (a *AirTraffic) SearchPossibleRoutes(plane *client.Plane, response *[]client.Routes) error {

	//Intanciar un slice que contendrá las posibles rutas
	posibleRoutes := []client.Routes{}

	//Itera entre todas las rutas buscando una coincidencia para agregarla al slice
	for _, route := range a.Routes {
		if plane.CurrentAirport == route.TakeoffAirport{
			posibleRoutes = append(posibleRoutes, route)
		}
	}

	//Establecer las posibles rutas al cliente
	*response = posibleRoutes
	return nil

}

//Método a servir al cliente para que solicite permiso de aterrizaje o despegue
func (a *AirTraffic) RequestPermission(plane *client.Plane, response *bool) error {

	//Llamar el método para obtener el índice del aereopuerto del cliente
	airportIndex := a.searchAirport(plane.CurrentAirport)

	//Validar que el aereopuerto se encutre libre de cualquier proceso de aterriazaje o despegue
	if a.airports[airportIndex].IsFree {
		*response = true
		a.airports[airportIndex].IsFree = false
	}

	return nil

}

//Método a servir al cliente para que éste confirme la finalización del proceso de aterrizaje o despege
func (a *AirTraffic) ConfirmOperation(plane *client.Plane, response *bool) error {

	airportIndex := a.searchAirport(plane.CurrentAirport)
	a.airports[airportIndex].IsFree = true

	return nil

}

//Método a servir al cliente para que éste valide su avión
func (a *AirTraffic) ValidatePlane(plane *client.Plane, response *bool) error {

	//Iterar y validar el aerepuerto inical del avión con los registrados
	for _, airport := range a.airports{
		if airport.Name == plane.CurrentAirport {
			*response = true
			return nil
		}
	}
	*response = false
	return nil

}

//Método para iniciar el servidor
func (a *AirTraffic) StartServer(port string) {

	//Llamar a los métodos que obtienen y resgistran los aerepuertos y rutas existentes
	a.getAirports()
	a.getRoutes()

	//Configurar el servidor de procesos remotos
	rpc.Register(a)
	rpc.HandleHTTP()

	//Establecer el protocolo y puerto a usar
	l, err := net.Listen("tcp", port)

	//Manejo del errorr
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	//Iniciar el servidor con sus configuraciones
	fmt.Println("Escuchando el puerto ", port)
	http.Serve(l, nil)

}

//Método para obtener y registrar los aereopuertos existentes
func (a *AirTraffic) getAirports() {

	airportsStr, err := ioutil.ReadFile("resources/aiportsInput.json")
	if err != nil{
		log.Fatal("Error al obtener los aerepuertos: ", err)
	}
	json.Unmarshal(airportsStr, &a.airports)

}

//Método para obtener y registrar las rutas existentes
func (a *AirTraffic) getRoutes() {

	routesStr, err := ioutil.ReadFile("resources/routesInput.json")
	if err != nil{
		log.Fatal("Error al obtener los aerepuertos: ", err)
	}
	json.Unmarshal(routesStr, &a.Routes)

}

//Método para buscar el índice de un aereopuerto
func (a *AirTraffic) searchAirport(name string) int {

	result := -1

	for i := 0; i < len(a.airports); i++ {
		if a.airports[i].Name == name{
			return i
		}
	}

	return result

}