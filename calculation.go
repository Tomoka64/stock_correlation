package main

import (
	"fmt"
	"log"
	"math"
)

func standardDeviation(xs []float64) float64 {
	return math.Sqrt(variance(xs))
}

func variance(vector []float64) float64 {
	n := 0.0
	mean := 0.0
	S := 0.0
	delta := 0.0

	for _, v := range vector {
		n++
		delta = v - mean
		mean = mean + (delta / n)
		S += delta * (v - mean)
	}

	return S / (n - 1)
}

func relativize(data []float64) []float64 {
	nv := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		nv[i-1] = (data[i] - data[i-1]) / data[i-1]
	}
	return nv
}

func Correlation(xs, ys []float64) (float64, float64, float64, float64) {
	coviance := covariance(xs, ys)
	StanX := standardDeviation(xs)
	StanY := standardDeviation(ys)
	Corre := covariance(xs, ys) / (standardDeviation(xs) * standardDeviation(ys))
	return coviance, StanX, StanY, Corre
}

func covariance(x, y []float64) float64 {
	if len(x) != len(y) {
		fmt.Println(len(x), len(y))
		log.Fatalln("length not match")
	}

	n := len(x)
	sum, xsum, xmean, ysum, ymean := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < n; i++ {
		xsum += x[i]
		ysum += y[i]
	}
	xmean = xsum / float64(n)
	ymean = ysum / float64(n)

	for i := 0; i < n; i++ {
		sum += (x[i] - xmean) * (y[i] - ymean)
	}

	return sum / float64(n-1)
}
