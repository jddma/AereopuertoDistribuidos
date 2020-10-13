package client

import (
	"log"
	"net/rpc"
)

type Plane struct{}

func (a *Plane) StartClient(host string)  {

	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Fatal("Dialing error: ", err)
	}

	var response string
	err = client.Call("AirTraffic.RequestLanding", a, &response)
	if err != nil {
		log.Fatal("Call error: ", err)
	}

}

func NewPlane() *Plane {

	return &Plane{}

}