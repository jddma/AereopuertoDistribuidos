package client

import (
	"fmt"
	"log"
	"net/rpc"
)

type Plane struct{
	CurrentAirport string
	Model string
	InFligth bool
}

func (p *Plane) StartClient(host string)  {

	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Fatal("Dialing error: ", err)
	}

	var planeIsValid bool
	err = client.Call("AirTraffic.ValidatePlane", p, &planeIsValid)
	if err != nil {
		log.Fatal("Call error: ", err)
		return
	}

	if ! planeIsValid {
		log.Fatal("El aereopuerto ingresado no se encuentra registrado")
		return
	}

	for  {
		if p.InFligth{
			//
		}else {
			var option string
			fmt.Print("Presione enter para solicitar permiso de despeque, si desea finalizar la operación digite cualquier cosa y luego enter: ")
			fmt.Scanln(&option)
			if option != "" {
				break
			}

			//solicitar permiso
		}
	}

}

func NewPlane() *Plane {

	var startAirport string
	fmt.Print("Digite el aereopuerto de incio: ")
	fmt.Scanln(&startAirport)

	var model string
	fmt.Print("Digite el modelo del avión: ")
	fmt.Scanln(&model)

	return &Plane{
		CurrentAirport: startAirport,
		Model: model,
		InFligth: false,
	}

}