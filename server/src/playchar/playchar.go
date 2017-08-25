package playchar

import "fmt"

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

func (pl *Playchar)ToString() string{
	var resultstr string
	resultstr = fmt.Sprintf("'%d:%d:%d'",pl.Healthpoints,pl.Position.Posx, pl.Position.Posy)
	return resultstr
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
