package playchar

type position struct {
	Posx int
	Posy int
}

type Playchar struct {
	Healthpoints int
	//position *cell //or i,j int
	// address of cell may not survive save-load process
	Position position
}

func New(hp int, x0 int, y0 int) *Playchar {
	return &Playchar{
		Healthpoints: hp,
		Position: position{
			Posx: x0,
			Posy: y0,
		},
	}
}
