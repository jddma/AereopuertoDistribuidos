package server

type planeRegister struct {
	id int
	inFligth bool
	airport string
}

func NewPlaneRegister(id int, inFligth bool, airport string) *planeRegister {

	return &planeRegister{
		id: id,
		inFligth: inFligth,
		airport: airport,
	}

}
