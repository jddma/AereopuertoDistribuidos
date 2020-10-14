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

//Método para manejar los errores al momento de usar procesos remotos
func (p *Plane) handleCallError(err error)  {

	if err != nil {
		log.Fatal("Call error: ", err)
	}

}

//Método para mostrar al usuario los posibles destinos
func (p *Plane) showPosibleRoutes(posibleRoutes []Routes)  {

	i := 1
	for _, route := range posibleRoutes {
		fmt.Printf("\t*** Destino #%d ***\n", i)
		fmt.Println("\t   -",route.DestinationAirport)
		i++
	}

}

//Método para obtener la rutas posibles deacuerdo al aerepuerto actual
func (p *Plane) getRoute(client *rpc.Client)  {

	//Solicitar al servidor las rutas posibles
	var posibleRoutes []Routes
	err := client.Call("AirTraffic.SearchPossibleRoutes", p, &posibleRoutes)
	p.handleCallError(err)

	//Validar que existan rutas posibles
	if len(posibleRoutes) == 0 {
		log.Fatal("No existe una ruta posible para el aereopuerto actual")
		return
	}

	//Obtener el aereopuerto de destino por parte del usuario
	var destinationAirportIndex int
	p.showPosibleRoutes(posibleRoutes)
	fmt.Print("Digite el numero del destino: ")
	fmt.Scanln(&destinationAirportIndex)

	//Establecer el aereopuerto de destino
	p.Destination = posibleRoutes[destinationAirportIndex - 1].DestinationAirport

}

//Método para inciar el proceso de despegue
func (p *Plane) trigerTakeoff(client *rpc.Client)  {

	//Obtener la rutas posibles deacuerdo al aereopueto actual en el que se encuentra el avión
	p.getRoute(client)

	//Solicitar permiso al servidor para despegar
	var permission bool
	err := client.Call("AirTraffic.RequestPermission", p, &permission)
	p.handleCallError(err)

	//Validar que el permiso fue concedido
	if permission {

		//Esperar al usuario para finalizar el despege
		fmt.Print("Proceso de despegue autorizado presione enter para indicar la culminación del despegue: ")
		fmt.Scanln()

		//Confirmar la finalización del proceso al servidor
		err = client.Call("AirTraffic.ConfirmOperation", p, &permission)
		p.handleCallError(err)

		//Cambiar el aereopuerto actual
		p.CurrentAirport = p.Destination
		p.Destination = ""

		//Cambiar el estado de vuelo
		p.InFligth = true
	} else {
		fmt.Println("Permiso denegado espere un momento y vuelva a intentarlo")
	}

}

//Método para iniciar el proceso de aterrizaje
func (p *Plane) trigerLanding(client *rpc.Client)  {

	//Solicitar al servidor permiso para aterrizar
	var permission bool
	err := client.Call("AirTraffic.RequestPermission", p, &permission)
	p.handleCallError(err)

	//Validar si el permiso fue concedido
	if permission {

		//Esperar al usuario para finalizar el proceso de aterrizaje
		fmt.Print("Aterrizaje autorizado presione enter para finalizar el proceso: ")
		fmt.Scanln()

		//Confimar la finalización del aterrizaje al servidor
		err = client.Call("AirTraffic.ConfirmOperation", p, &permission)
		p.handleCallError(err)

		//Cambiar el estado de vuelo del avión
		p.InFligth = false
	} else {
		fmt.Println("Permiso denegado espere un momento e intentelo nuevamente")
	}

}

//Método que inicia el cliente
func (p *Plane) StartClient(host string)  {

	//Establecer la conexión con el servidor
	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Fatal("Dialing error: ", err)
	}

	//Solicitar al servidor la validación del avión
	var planeIsValid bool
	err = client.Call("AirTraffic.ValidatePlane", p, &planeIsValid)
	p.handleCallError(err)

	//En caso de que al avión no sea valido el proceso finalizará
	if ! planeIsValid {
		log.Fatal("El aereopuerto ingresado no se encuentra registrado")
		return
	}

	//Bucle para leer las opciones del usuario
	for {

		//Validar si el siguiente proceso es aterrizar o despeguar
		if p.InFligth{
			fmt.Print("Avión en vuelo para solicitar permiso de aterrizaje al destino presione enter: ")
			fmt.Scanln()

			//Llamar al método que inicia el despeque
			p.trigerLanding(client)
		}else {

			//Validar si el usuario desea finalzar con la operación o solicitar despegue
			var option string
			fmt.Print("Presione enter para solicitar permiso de despeque, si desea finalizar la operación digite cualquier cosa y luego enter: ")
			fmt.Scanln(&option)
			if option != "" {
				break
			}

			//Llamar al método que inicia el despegue
			p.trigerTakeoff(client)
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