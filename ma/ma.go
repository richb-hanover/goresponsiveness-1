package ma

import (
	"github.com/hawkinsw/goresponsiveness/saturating"
	"github.com/hawkinsw/goresponsiveness/utilities"
)

// Convert this to a Type Parameterized interface when they are available
// in Go (1.18).
type MovingAverage struct {
	intervals int
	instants  []float64
	index     int
	divisor   *saturating.SaturatingInt
}

func NewMovingAverage(intervals int) *MovingAverage {
	return &MovingAverage{instants: make([]float64, intervals), intervals: intervals, divisor: saturating.NewSaturatingInt(intervals)}
}

func (ma *MovingAverage) AddMeasurement(measurement float64) {
	ma.instants[ma.index] = measurement
	ma.divisor.Add(1)
	// Invariant: ma.index always points to the oldest measurement
	ma.index = (ma.index + 1) % ma.intervals
}

func (ma *MovingAverage) CalculateAverage() float64 {
	total := float64(0)
	for i := 0; i < ma.intervals; i++ {
		total += ma.instants[i]
	}
	return float64(total) / float64(ma.divisor.Value())
}

func (ma *MovingAverage) AllSequentialIncreasesLessThan(limit float64) bool {

	// If we have not yet accumulated a complete set of intervals,
	// this is false.
	if ma.divisor.Value() != ma.intervals {
		return false
	}

	// Invariant: ma.index always points to the oldest (see AddMeasurement above)
	oldestIndex := ma.index
	previous := ma.instants[oldestIndex]
	for i := 1; i < ma.intervals; i++ {
		currentIndex := (oldestIndex + i) % ma.intervals
		current := ma.instants[currentIndex]
		percentChange := utilities.SignedPercentDifference(current, previous)
		previous = current
		if percentChange > limit {
			return false
		}
	}
	return true
}
