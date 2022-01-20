package lat

import (
	"fmt"
	"math/rand"
	"os"
)

var kappa [2][]int
var d []int

var ux, uy, u, vx, vy, v int
var link *int

var obs [5]float64
var obs_tot [5]float64

func InitField(restart bool) {
	kappa[0] = make([]int, V)
	kappa[1] = make([]int, V)
	d = make([]int, V)

	if restart {
		read_conf()
	} else {
		for ix := 0; ix < V; ix++ {
			kappa[0][ix] = 0
			kappa[1][ix] = 0
			d[ix] = 0
		}
		ux = rand.Intn(L)
		uy = rand.Intn(L)
		u = ux + uy*L
		vx = ux
		vy = uy
		v = vx + vy*L
		d[u] = 2
	}
}

func read_conf() {
	f, err := os.Open("prova.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = fmt.Fscan(f, &ux, &uy, &vx, &vy)
	if err != nil {
		panic(err)
	}
	u = ux + uy*L
	v = vx + vy*L
	for ix := 0; ix < V; ix++ {
		_, err = fmt.Fscan(f, &kappa[0][ix], &kappa[1][ix])
		if err != nil {
			panic(err)
		}
	}
}
