package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/guagno/goworm/ctilde"
	"github.com/guagno/goworm/lat"
	"github.com/guagno/goworm/wdb"
)

var ell int
var lam, beta float64
var restart bool
var pid string
var nMeas, nSweepPerMeas, nDisc int

func read_param() {
	flag.IntVar(&ell, "L", -1, "A value for the lattice size")
	flag.Float64Var(&lam, "lam", -1.0, "A value for lambda")
	flag.Float64Var(&beta, "beta", -1.0, "A value for beta")
	flag.IntVar(&nMeas, "nMeas", 100, "N Measures")
	flag.IntVar(&nSweepPerMeas, "nSweep", 10, "N Sweep per meas.")
	flag.IntVar(&nDisc, "nDisc", 100, "N Discarded sweeps")
	flag.BoolVar(&restart, "restart", false, "Restarting boolean (default false)")
	flag.StringVar(&pid, "pid", "", "Process ID")

	flag.Parse()
	if ell < 0 {
		panic("L cannot be negative!")
	}
	if lam < 0 {
		panic("Lambda cannot be negative!")
	}
	if beta < 0 {
		panic("Beta cannot be negative!")
	}
	if pid == "" {
		panic("pid cannot be empty!")
	}

	fmt.Printf("L       = %v\n", ell)
	fmt.Printf("lambda  = %v\n", lam)
	fmt.Printf("beta    = %v\n", beta)
	fmt.Printf("pid     = %v\n", pid)
	fmt.Printf("restart = %v\n", restart)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	read_param()
	wdb.Init_DB()
	ctilde.Init_ctilde(lam)
	lat.InitGeometry(ell)
	lat.InitField(restart)

	fmt.Println("Starting thermalization")
	for i := 1; i <= nDisc; i++ {
		lat.Sweep(beta)
	}
	fmt.Println("Starting measure sweeps")
	for iMeas := 1; iMeas <= nMeas; iMeas++ {
		lat.ResetObs()
		for i := 0; i < nSweepPerMeas; i++ {
			lat.Sweep(beta)
		}
		lat.Measure(iMeas, nSweepPerMeas)
	}
}
