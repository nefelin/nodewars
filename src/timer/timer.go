package timer

import (
	"fmt"
	"time"
)

// Timer is used to maintain a collection of functions that should be called at varying intervals, regulated by a central clock that can be stopped and started again.
type Timer struct {
	status     timerStatus
	jobs       map[jobID]*timerJob
	resolution time.Duration
}

type jobID = string
type timerStatus int
type tickableFunc = func(time.Duration)

const (
	running timerStatus = iota
	stopped
)

type timerJob struct {
	fn      tickableFunc
	freq    time.Duration
	dormant time.Duration
}

// NewTimer Returns a new Timer object
func NewTimer() *Timer {
	return &Timer{
		status:     stopped,
		jobs:       make(map[jobID]*timerJob),
		resolution: 100 * time.Millisecond,
	}
}

// Start Begins the Timer's internal clock and job monitoring
func (gt *Timer) Start() *Timer {
	gt.status = running
	go timerClock(gt)
	return gt
}

// Stop Stops the Timer's internal clock and job monitoring
func (gt *Timer) Stop() *Timer {
	gt.status = stopped
	return gt
}

// SetRes Sets the interval of the Timer's internal clock, determining how frequently job dormancy is checked
func (gt *Timer) SetRes(r time.Duration) *Timer {
	gt.resolution = r
	return gt
}

// AddJob adds a job to the Timer's watched jobs. Job's function is called on every tick
func (gt *Timer) AddJob(j jobID, t tickableFunc) (*Timer, error) {
	if _, ok := gt.jobs[j]; ok {
		err := fmt.Errorf("duplicate jobID: '%s'", j)
		return gt, err
	}

	gt.jobs[j] = &timerJob{t, 0, 0}
	return gt, nil
}

// AddScheduledJob behaves just like AddJob execpt job's function is only called after specified duration
func (gt *Timer) AddScheduledJob(j jobID, t tickableFunc, d time.Duration) (*Timer, error) {
	if _, ok := gt.jobs[j]; ok {
		err := fmt.Errorf("duplicate jobID: '%s'", j)
		return gt, err
	}

	gt.jobs[j] = &timerJob{t, d, 0}
	return gt, nil
}

// KillJob removes a job from the watched job list
func (gt *Timer) KillJob(j jobID) (*Timer, error) {
	if _, ok := gt.jobs[j]; !ok {
		err := fmt.Errorf("jobID not found: '%s'", j)
		return gt, err
	}

	delete(gt.jobs, j)
	return gt, nil
}

func timerClock(gt *Timer) {
	for gt.status == running {
		start := time.Now()

		<-time.After(gt.resolution)

		for _, j := range gt.jobs {
			processJob(j, time.Since(start))
		}
	}
}

func processJob(j *timerJob, e time.Duration) {
	// fmt.Printf("Processing job: %+v\n", j)
	j.dormant += e
	if j.dormant >= j.freq {
		j.fn(j.dormant)
		j.dormant = 0
	}
}
