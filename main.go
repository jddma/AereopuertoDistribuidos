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
		client := client.NewPlane()
		client.StartClient("127.0.0.1:8082")
		break

	case 2:
		server := new(server.AirTraffic)
		server.StartServer(":8082")
		break

	}

}