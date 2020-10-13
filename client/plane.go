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
	Destination string
}

func (p *Plane) handleCallError(err error)  {

	if err != nil {
		log.Fatal("Call error: ", err)
	}

}

func (p *Plane) showPosibleRoutes(posibleRoutes []Routes)  {

	i := 1
	for _, route := range posibleRoutes {
		fmt.Printf("\t*** Destino #%d ***\n", i)
		fmt.Println("\t   -",route.DestinationAirport)
		i++
	}

}

func (p *Plane) getRoute(client *rpc.Client)  {

	var posibleRoutes []Routes
	err := client.Call("AirTraffic.SearchPossibleRoutes", p, &posibleRoutes)
	p.handleCallError(err)

	if len(posibleRoutes) == 0 {
		log.Fatal("No existe una ruta posible para el aereopuerto actual")
		return
	}

	var destinationAirportIndex int
	p.showPosibleRoutes(posibleRoutes)
	fmt.Print("Digite el numero del destino: ")
	fmt.Scanln(&destinationAirportIndex)
	p.Destination = posibleRoutes[destinationAirportIndex - 1].DestinationAirport

}

func (p *Plane) StartClient(host string)  {

	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Fatal("Dialing error: ", err)
	}

	var planeIsValid bool
	err = client.Call("AirTraffic.ValidatePlane", p, &planeIsValid)
	p.handleCallError(err)

	if ! planeIsValid {
		log.Fatal("El aereopuerto ingresado no se encuentra registrado")
		return
	}



	for {

		//Validar si el siguiente proceso es aterrizar o despeguar
		if p.InFligth{
			fmt.Print("Avión en vuelo para solicitar permiso de aterrizaje al destino presione enter: ")
			fmt.Scanln()

			var permission bool
			err = client.Call("AirTraffic.RequestPermission", p, &permission)
			p.handleCallError(err)

			if permission {
				fmt.Print("Aterrizaje autorizado presione enter para finalizar el proceso: ")
				fmt.Scanln()

				err = client.Call("AirTraffic.ConfirmOperation", p, &permission)
				p.handleCallError(err)
				p.InFligth = false
			} else {
				fmt.Println("Permiso denegado espere un momento e intentelo nuevamente")
			}

		}else {
			p.getRoute(client)
			var option string
			fmt.Print("Presione enter para solicitar permiso de despeque, si desea finalizar la operación digite cualquier cosa y luego enter: ")
			fmt.Scanln(&option)
			if option != "" {
				break
			}

			var permission bool
			err = client.Call("AirTraffic.RequestPermission", p, &permission)
			p.handleCallError(err)

			if permission {
				fmt.Print("Proceso de despegue autorizado presione enter para indicar la culminación del despegue: ")
				fmt.Scanln()
				err = client.Call("AirTraffic.ConfirmOperation", p, &permission)
				p.handleCallError(err)

				p.CurrentAirport = p.Destination
				p.Destination = ""
				p.InFligth = true
			} else {
				fmt.Println("Permiso denegado espere un momento y vuelva a intentarlo")
			}
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