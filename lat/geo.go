package lat

import (
	"math"
)

var L, V, Vhalf int
var Neigh [4][]int
var Coseno []float64

func InitGeometry(ell int) {
	L = ell
	V = L * L
	Vhalf = V / 2

	for i := 0; i < 4; i++ {
		Neigh[i] = make([]int, V)
	}

	for iy := 0; iy < L; iy++ {
		for ix := 0; ix < L; ix++ {
			iz := ix + L*iy
			if iy == 0 {
				Neigh[0][iz] = iz + L*(L-1)
			} else {
				Neigh[0][iz] = iz - L
			}
			if ix == L-1 {
				Neigh[1][iz] = iz - (L - 1)
			} else {
				Neigh[1][iz] = iz + 1
			}
			if iy == L-1 {
				Neigh[2][iz] = iz - L*(L-1)
			} else {
				Neigh[2][iz] = iz + L
			}
			if ix == 0 {
				Neigh[3][iz] = iz + (L - 1)
			} else {
				Neigh[3][iz] = iz - 1
			}
		}
	}

	Coseno = make([]float64, L/2+1)
	pstar := 2.0 * math.Pi / float64(L)
	for i := 0; i <= L/2; i++ {
		Coseno[i] = math.Cos(pstar * float64(i))
	}

}
