package main

import (
	"testing"

	"github.com/casen/snakegame/model"
)

func TestDetectCylcles(t *testing.T) {
	type testCases struct {
		visited  []model.Point
		startIdx int
		period   int
		hasCycle bool
	}

	cases := []testCases{
		{[]model.Point{}, 0, 0, false},
		{[]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}, 0, 0, false},
		{[]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}}, 0, 0, false},
		{[]model.Point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 0},
		}, 0, 2, true},
		{[]model.Point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
			{X: 0, Y: 4},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
			{X: 0, Y: 4},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
			{X: 0, Y: 4},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
		}, 2, 3, true},
	}

	for _, c := range cases {
		gotStartIdx, gotPeriod, gotHasCycle := DetectCycles(c.visited)
		if gotStartIdx != c.startIdx || gotPeriod != c.period || gotHasCycle != c.hasCycle {
			t.Errorf("DetectCycles(%v) = %v, %v, %v; want %v, %v %v", c.visited, gotStartIdx, gotPeriod, gotHasCycle, c.startIdx, c.period, c.hasCycle)
		}
	}

}
