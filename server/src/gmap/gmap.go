package gmap

import (
	"Cell"
	"fmt"
	"time"
	"math/rand"
)

type MapParams struct {
	Xsize int
	Ysize int
}

type Gmap struct {
	Params MapParams
	Field  [][]Cell.Cell
}

func MakeAMap(size int) *Gmap {
	cellarray := [][]Cell.Cell{}
	for j := 0; j < size; j++ {
		cellrow := []Cell.Cell{}
		for i := 0; i < size; i++ {
			celltoadd := Cell.Cell{1, 1}
			cellrow = append(cellrow, celltoadd)
		}
		cellarray = append(cellarray, cellrow)
	}
	return &Gmap{
		Params: MapParams{Xsize: size, Ysize: size},
		Field:  cellarray,
	}
}

func (gmap *Gmap)MapEventRandomizator(eventsNum int, size int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	field := gmap.Field

	for i:=0; i<size; i++ {
		for j:=0; j<size; j++ {
			if field[i][j].Kind != 0 {
				field[i][j].Kind = r.Perm(eventsNum)[0] + 1
			}
		}
	}
	gmap.Field = field
}

func MakeZeroMap(size int) *Gmap {
	cellarray := [][]Cell.Cell{}

	return &Gmap{
		Params: MapParams{Xsize: size, Ysize: size},
		Field:  cellarray,
	}
}

func (m *Gmap) InsertString() string {
	str := "ARRAY["
	for i, row := range m.Field {
		str += "["
		for j, v := range row {
			str += "'" + fmt.Sprintf("%d:%d", v.Kind, v.Hidden)
			if j < m.Params.Xsize-1 {
				str += "',"
			}
		}
		str += "']"
		if i < m.Params.Ysize-1 {
			str += ","
		}
	}
	str += "]"
	//fmt.Println(str)
	return str
}


