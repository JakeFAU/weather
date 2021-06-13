package main

import "strconv"

func toFahrenheit(c float64) float64 {
	return (c * 9 / 5.0) + 32
}

func toMercury(m float64) float64 {
	return m * 0.02953
}

func convert32(value float32) string {
	return strconv.FormatFloat(float64(value), 'e', 2, 32)
}

func convert64(value float64) string {
	return strconv.FormatFloat(float64(value), 'e', 2, 32)
}
