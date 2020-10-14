package main

import (
	"./client"
	"./server"
	"fmt"
)

func main()  {

	//Leer el modo de ejecución
	fmt.Println("Digite si desea ejecutar como:\n\t*Avión [1]\n\t*Controlador [2]")
	var option int
	fmt.Scanf("%d", &option)

	switch option {

	case 1:

		//Obtener el host a usar
		var host string
		fmt.Print("Digite el host del servidor: ")
		fmt.Scanln(&host)

		client := client.NewPlane()
		client.StartClient(host)
		break

	case 2:

		//Obtener el puerto a escuchar
		var port string
		fmt.Print("Digite el puerto a usar: ")
		fmt.Scanln(&port)

		server := new(server.AirTraffic)
		server.StartServer(":8082")
		break

	}

}