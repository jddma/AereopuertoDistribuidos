package server

type planeRegister struct {
	id int
	inFligth bool
	airport string
	enrollment string
}

func NewPlaneRegister(id int, inFligth bool, airport string, enrollment string) *planeRegister {

	return &planeRegister{
		id: id,
		inFligth: inFligth,
		airport: airport,
		enrollment: enrollment,
	}

}
