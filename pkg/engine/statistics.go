package engine

import (
	"time"

	"golife/pkg/universe"
)

// Statistics holds simulation statistics
type Statistics struct {
	Generation    int
	LivingCells   int
	Births        int
	Deaths        int
	StartTime     time.Time
	LastFrameTime time.Time
	FPS           float64
}

// NewStatistics creates a new Statistics instance
func NewStatistics(initialPopulation int) *Statistics {
	now := time.Now()
	return &Statistics{
		Generation:    0,
		LivingCells:   initialPopulation,
		StartTime:     now,
		LastFrameTime: now,
		FPS:           0,
	}
}

// Update updates the statistics for the current generation
func (s *Statistics) Update(u *universe.Universe2D) {
	prevLivingCells := s.LivingCells
	s.Generation++
	s.LivingCells = u.CountLiving()

	// Calculate births and deaths
	diff := s.LivingCells - prevLivingCells
	if diff > 0 {
		s.Births = diff
		s.Deaths = 0
	} else if diff < 0 {
		s.Births = 0
		s.Deaths = -diff
	} else {
		s.Births = 0
		s.Deaths = 0
	}

	// Calculate FPS
	now := time.Now()
	if !s.LastFrameTime.IsZero() {
		frameDuration := now.Sub(s.LastFrameTime).Seconds()
		if frameDuration > 0 {
			s.FPS = 1.0 / frameDuration
		}
	}
	s.LastFrameTime = now
}

// Reset resets statistics to initial state
func (s *Statistics) Reset(initialPopulation int) {
	now := time.Now()
	s.Generation = 0
	s.LivingCells = initialPopulation
	s.Births = 0
	s.Deaths = 0
	s.StartTime = now
	s.LastFrameTime = now
	s.FPS = 0
}
