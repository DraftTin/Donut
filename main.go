package main

import (
	"fmt"
	"math"
	"time"
)

var width, height float64 = 160, 44
var A, B, C float64
var distanceFromCam float64 = 80
var R, r float64
var horizontalOffset float64
var K1 float64 = 40

var zBuffer [160 * 44]float64
var buffer [160 * 44]rune

var backgroundASCIICode rune = ' '

var incrementSpeed float64 = 0.2

func calculateX(i float64, j float64, k float64) float64 {
	return j*math.Sin(A)*math.Sin(B)*math.Cos(C) - k*math.Cos(A)*math.Sin(B)*math.Cos(C) +
		j*math.Cos(A)*math.Sin(C) + k*math.Sin(A)*math.Sin(C) + i*math.Cos(B)*math.Cos(C)
}

func calculateY(i float64, j float64, k float64) float64 {
	return j*math.Cos(A)*math.Cos(C) + k*math.Sin(A)*math.Cos(C) -
		j*math.Sin(A)*math.Sin(B)*math.Sin(C) + k*math.Cos(A)*math.Sin(B)*math.Sin(C) -
		i*math.Cos(B)*math.Sin(C)
}

func calculateZ(i float64, j float64, k float64) float64 {
	return k*math.Cos(A)*math.Cos(B) - j*math.Sin(A)*math.Cos(B) + i*math.Sin(B)
}

func calculateForSurface(donutX, donutY, donutZ float64, ch rune) {
	x := calculateX(donutX, donutY, donutZ)
	y := calculateY(donutX, donutY, donutZ)
	z := calculateZ(donutX, donutY, donutZ) + distanceFromCam

	ooz := 1 / z
	xp := int(width/2 + horizontalOffset + x*K1*2*ooz)
	yp := int(height/2 + y*K1*ooz)

	idx := yp*int(width) + xp
	if idx >= 0 && float64(idx) < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
		}
		buffer[idx] = ch
	}
}

func refresh() {
	for i := 0; i < 160*44; i++ {
		zBuffer[i] = 0
		buffer[i] = backgroundASCIICode
	}
}

func main() {
	R = 20
	r = 10
	horizontalOffset = 10

	for {
		refresh()
		for theta := 0.0; theta < 2*math.Pi; theta += incrementSpeed {
			for sigma := 0.0; sigma < 2*math.Pi; sigma += incrementSpeed {
				donutX := (R + r*math.Cos(theta)) * math.Cos(sigma)
				donutY := (R + r*math.Cos(theta)) * math.Sin(sigma)
				donutZ := r * math.Sin(theta)
				calculateForSurface(donutX, donutY, donutZ, '+')
			}
		}
		w, h := int(width), int(height)
		fmt.Printf("\x1b[H")
		for i := 0; i < w*h; i++ {
			if i%w == 0 {
				fmt.Printf("\n")
			} else {
				fmt.Printf("%c", buffer[i])
			}
		}
		A += 0.05
		B += 0.05
		C += 0.01
		time.Sleep(time.Millisecond * 16)
	}
}
