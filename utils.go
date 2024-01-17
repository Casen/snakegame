package main

import (
	"github.com/casen/snakegame/model"
)

/*
 * This is using Floyd's cycle-finding algorithm
 */
func DetectCycles(visited []model.Point) (startIdx int, period int, hasCycle bool) {
	t := 1 // Tortoise
	h := 2 // Hare

	// Min snake length is 4, so we need at least 12 points to have a cycle. Let's pad that a bit, and go with 32
	if len(visited) < 32 {
		return 0, 0, false
	}

	for visited[t] != visited[h] && h < len(visited)-2 {
		t++
		h += 2
	}

	if visited[t] != visited[h] || t < 3 {
		return 0, 0, false
	}

	period = t

	t = 0
	for visited[t] != visited[h] && h < len(visited)-1 {
		t++
		h++
	}

	if visited[t] != visited[h] {
		return 0, 0, false
	}

	return t, period, true
}
