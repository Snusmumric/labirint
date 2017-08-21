package gmap

import (
	"../cell"
)

type mapParams struct {
	xSize int
	ySize int
}

type Gmap struct {
	Params mapParams
	Field  [][]cell.Cell
}

func MakeAMap(size int) *Gmap {
	cellarray := [][]cell.Cell{}
	for j := 0; j < size; j++ {
		cellrow := []cell.Cell{}
		for i := 0; i < size; i++ {
			celltoadd := cell.Cell{1, 1}
			cellrow = append(cellrow, celltoadd)
		}
		cellarray = append(cellarray, cellrow)
	}
	return &Gmap{
		Params: mapParams{xSize: size, ySize: size},
		Field:  cellarray,
	}
}

func (m *Gmap) InsertString() string {
	str := "ARRAY["
	for i, row := range m.Field {
		str += "["
		for j, v := range row {
			str += "'" + fmt.Sprintf("%d:%d", v.Kind, v.Hidden) . "o'" 	
			if j < m.Params.xSize - 1 {
				str += ','
			}
		}
		str += ']'
		if i < m.Params.ySize - 1 {
			str += ','
		}
	}
	str += ']'
}
