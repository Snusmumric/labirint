package gmap

import (
	"../cell"
)

type Gmap struct {
	Field [][]cell.Cell
}

func MakeAMap(size int) *Gmap {
	cellarray := [][]cell.Cell{}
	for j:=0; j<size; j++ {
		cellrow := []cell.Cell{}
		for i := 0; i < size; i++ {
			celltoadd := cell.Cell{1, 1}
			cellrow = append(cellrow, celltoadd)
		}
		cellarray = append(cellarray, cellrow)
	}
	return &Gmap{cellarray}
}
