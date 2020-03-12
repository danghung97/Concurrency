package runner

import(
	"os"
	"time"
	"os/signal"
	"errors"
)

type Runner struct {
	interrupt chan os.Signal
	complete chan error
	timeout <-chan time.Time
	tasks []func(int)
}

var ErrTimeout = errors.New("Receive time out")
var ErrInterrupt = errors.New("receive interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		// buffered channel with a buffer of 1, this guarantees at least one os.Signal value
		// is received from the runtime. The runtime sends this event in a nonblocking away
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

func (r *Runner) Add(task ...func(int)) {
	r.tasks = append(r.tasks, task...)
}

func (r *Runner) Start() error {
	// we want to receive all interrupt based signal
	signal.Notify(r.interrupt, os.Interrupt)
	
	go func() {
		r.complete <- r.run()
	}()
	select {
	case err := <-r.complete:
		return err
	case <- r.timeout:
		return ErrTimeout
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		// check for an interrupt signal from OS
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	// signaled when an interrupt event is sent
	case <-r.interrupt:
		signal.Stop(r.interrupt) //stop receiving any signals
		return true
	default:
		return false
	}
}
