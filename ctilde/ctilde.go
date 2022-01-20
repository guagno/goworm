package ctilde

import (
	//        "fmt"
	"math"
	"runtime"

	"gonum.org/v1/gonum/integrate/quad"
)

const MAXN = 201
const MAXINDEX = 2*MAXN + 1

var Ctil [MAXINDEX]float64

func Init_ctilde(lam float64) {

	startIter0 := 100000
	startIter := startIter0

	f0 := func(x float64) float64 {
		x2 := x * x
		x4 := (x2 - 1.0) * (x2 - 1.0)
		return math.Exp(-x2 - lam*x4)
	}
	f2 := func(x float64) float64 {
		x2 := x * x
		x4 := (x2 - 1.0) * (x2 - 1.0)
		return x2 * math.Exp(-x2-lam*x4)
	}
	f4 := func(x float64) float64 {
		x2 := x * x
		x4 := (x2 - 1.0) * (x2 - 1.0)
		return x2 * x2 * math.Exp(-x2-lam*x4)
	}

	X := func(n float64) float64 {
		return math.Sqrt(n*(n+1.0)) / (n + 0.5)
	}
	Y := func(n float64) float64 {
		gamma := -0.5 * math.Log(2.0*lam)
		return 2.0 * math.Sinh(gamma) / (math.Sqrt(n) * (1.0 + 0.5/n))
	}
	errexp := func(n float64) float64 {
		gamma := -0.5 * math.Log(2.0*lam)
		s := math.Sinh(gamma)
		s2 := s * s
		s4 := s2 * s2
		srn12 := 1.0 / math.Sqrt(n)
		srn32 := srn12 / n
		return 1.0 - s*srn12 + 0.5*s2/n + 0.25*s*srn32 + (1.0-2.0*s4)/(12.0*n*n)
	}
	//	erre := func(n float64, ct float64) float64 {
	//		return math.Sqrt(2.0*lam/n) / ct
	//	}
	iter_forw := func(n float64, rprime float64) float64 {
		return (1.0 / X(n)) * (1.0/rprime - Y(n))
	}
	iter_back := func(n float64, rprime float64) float64 {
		return 1.0 / (rprime*X(n) + Y(n))
	}

	concurrent := runtime.GOMAXPROCS(0)
	ev0 := quad.Fixed(f0, math.Inf(-1), math.Inf(1), 10000, nil, concurrent)
	ev2 := quad.Fixed(f2, math.Inf(-1), math.Inf(1), 10000, nil, concurrent)

	Ctil[2] = ev0 / ev2

	if lam < 0.5 {
		rprime := errexp(float64(startIter))
		for n := startIter; n >= MAXN; n-- {
			rprime = iter_back(float64(n), rprime)
		}
		for n := MAXN - 1; n >= 1; n-- {
			rprime = iter_back(float64(n), rprime)
			index := 2 * (n + 1)
			Ctil[index] = math.Sqrt(2.0*lam/float64(n)) / rprime
		}
	} else {
		ev4 := quad.Fixed(f4, math.Inf(-1), math.Inf(1), 10000, nil, concurrent)
		Ctil[4] = ev2 / ev4
		rprime := math.Sqrt(2.0*lam) / Ctil[4]
		for n := 2; n < MAXN; n++ {
			rprime = iter_forw(float64(n-1), rprime)
			index := 2 * (n + 1)
			Ctil[index] = math.Sqrt(2.0*lam/float64(n)) / rprime
		}
	}
}
