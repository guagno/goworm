package lat

import (
	"fmt"
	"math/rand"

	"github.com/guagno/goworm/ctilde"
)

var tot [5]float64

func periodic(x int) int {
	if x < 0 {
		x = -x
	}
	if x > L/2 {
		x = L - x
	}
	return x
}

func ResetObs() {
	for i := 0; i < 5; i++ {
		obs_tot[i] = 0.0
	}
}

func Measure(iMeas int, nSweepPerMeas int) {
	fMeas := float64(iMeas)
	for i := 0; i < 5; i++ {
		obs_tot[i] /= float64(nSweepPerMeas)
		tot[i] += obs_tot[i]
	}
	fmt.Printf("%v %v %v %v\n", iMeas, tot[0]/fMeas, tot[1]/fMeas, tot[2]/fMeas)
}

func Sweep(beta float64) {

	accept := 0.0

	get_link := func(w int) (int, int, int, *int) {
		var link *int
		dir := rand.Intn(4) // una delle quattro direzioni
		wn := Neigh[dir][w]
		wny := wn / L
		wnx := wn - wny*L
		switch dir {
		case 0:
			link = &kappa[1][wn]
		case 1:
			link = &kappa[0][w]
		case 2:
			link = &kappa[1][w]
		case 3:
			link = &kappa[0][wn]
		}
		return wn, wnx, wny, link
	}

	inc_or_dec_link := func(w *int, wx *int, wy *int, wn int, wnx int, wny int) {
		incdec := rand.Intn(2) // incremento o decremento kappa?
		r := rand.Float64()
		if incdec == 0 { // incremento, posso sempre provarci
			prob := beta / (float64(*link+1) * (ctilde.Ctil[d[wn]+2]))
			if r < prob {
				// se accetto, prima di tutto updato d(un); d(u) invariato
				d[wn] = d[wn] + 2
				*w = wn
				*wx = wnx
				*wy = wny
				*link = *link + 1

				// anche obs(4) (cioè K) va incrementato...
				// KK = KK + 1
				accept++
			}
		} else { // decremento, posso provarci solo se kappa > 0
			if *link > 0 {
				prob := float64(*link) * ctilde.Ctil[d[*w]] / beta
				if r < prob {
					// se accetto, prima di tutto updato d(u); d(un) invariato
					d[*w] = d[*w] - 2
					*w = wn
					*wx = wnx
					*wy = wny
					*link = *link - 1

					// anche obs(4) (cioè K) va decrementato...
					//	KK = KK - 1

					accept++
				}
			}
		}
	}

	new_worm := func() {
		r := rand.Float64()
		un := rand.Intn(V)
		prob := ctilde.Ctil[d[u]] / ctilde.Ctil[d[un]+2]
		if r < prob {
			// sposto u e v, quindi decremento d
			d[u] = d[u] - 2
			uny := un / L
			unx := un - uny*L
			u = un
			ux = unx
			uy = uny
			v = un
			vx = unx
			vy = uny
			d[u] = d[u] + 2
		}
	}

	cumulate_obs := func() {
		if u == v {
			obs[0] = obs[0] + 1.0
			obs[1] = obs[1] + ctilde.Ctil[d[u]]
		}

		dx := periodic(ux - vx)
		dy := periodic(uy - vy)
		cc := 0.5 * (Coseno[dx] + Coseno[dy])
		obs[2] = obs[2] + cc
		//		obs[3] = obs[3] + KK
		//		obs[4] = obs[4] + cc*KK
	}

	for i := 0; i < 5; i++ {
		obs[i] = 0.0
	}
	for s := 0; s < Vhalf; s++ {
		var un, unx, uny, vn, vnx, vny int

		// Move Worm Head
		un, unx, uny, link = get_link(u)
		inc_or_dec_link(&u, &ux, &uy, un, unx, uny)

		// il secondo step lo si fa solo in questo caso
		if u == v {
			new_worm()
		}

		cumulate_obs()

		// Move Worm Tail
		vn, vnx, vny, link = get_link(v)
		inc_or_dec_link(&v, &vx, &vy, vn, vnx, vny)

		// il secondo step lo si fa solo in questo caso
		if u == v {
			new_worm()
		}

		cumulate_obs()

	} // end for on sites

	for i := 0; i < 5; i++ {
		obs_tot[i] += (obs[i] / float64(V))
	}
	accept /= float64(V)

}
