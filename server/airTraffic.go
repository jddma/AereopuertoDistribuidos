package server

import (
	"../client"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type AirTraffic struct {}

//Método que el avión usara remotamente para solicitar permiso para aterrizar
func (t *AirTraffic) RequestLanding(avion *client.Plane, response *string) error {

	*response = "Aprobado"

	return nil

}

func (t *AirTraffic) StartServer(port string) {

	rpc.Register(t)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	fmt.Println("Escuchando el puerto ", port)
	http.Serve(l, nil)

}
