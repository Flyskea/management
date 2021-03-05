package middlewares

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/andres-erbsen/clock"
	"github.com/gin-gonic/gin"
)

var defaultNegelectMethods = []string{}
var defaultTimeLimitPerAct = 5
var defaultPer perOption = perOption(time.Second)
var defaultMaxSlack slackOption = 10

// Note: This file is inspired by:
// https://github.com/prashantv/go-bench/blob/master/ratelimit

// Limiter is used to rate-limit some process, possibly across goroutines.
// The process is expected to call Take() before every iteration, which
// may block to throttle the goroutine.
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
}

// Clock is the minimum necessary interface to instantiate a rate limiter with
// a clock or mock clock, compatible with clocks created using
// github.com/andres-erbsen/clock.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// config configures a limiter.
type config struct {
	clock    Clock
	maxSlack time.Duration
	per      time.Duration
}

// New returns a Limiter that will limit to the given RPS.
func New(rate int, opts ...Option) Limiter {
	return newAtomicBased(rate, opts...)
}

// buildConfig combines defaults with options.
func buildConfig(opts []Option) config {
	c := config{
		clock:    clock.New(),
		maxSlack: 10,
		per:      time.Second,
	}
	for _, opt := range opts {
		opt.apply(&c)
	}
	return c
}

// Option configures a Limiter.
type Option interface {
	apply(*config)
}

type clockOption struct {
	clock Clock
}

func (o clockOption) apply(c *config) {
	c.clock = o.clock
}

// WithClock returns an option for ratelimit.New that provides an alternate
// Clock implementation, typically a mock Clock for testing.
func WithClock(clock Clock) Option {
	return clockOption{clock: clock}
}

type slackOption int

func (o slackOption) apply(c *config) {
	c.maxSlack = time.Duration(o)
}

// WithoutSlack is an Option for ratelimit.New that initializes the limiter
// without any initial tolerance for bursts of traffic.
var WithoutSlack Option = slackOption(0)

type perOption time.Duration

func (p perOption) apply(c *config) {
	c.per = time.Duration(p)
}

// Per allows configuring limits for different time windows.
//
// The default window is one second, so New(100) produces a one hundred per
// second (100 Hz) rate limiter.
//
// New(2, Per(60*time.Second)) creates a 2 per minute rate limiter.
func Per(per time.Duration) Option {
	return perOption(per)
}

type unlimited struct{}

// NewUnlimited returns a RateLimiter that is not limited.
func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() time.Time {
	return time.Now()
}

type state struct {
	last     time.Time
	sleepFor time.Duration
}

type atomicLimiter struct {
	state unsafe.Pointer
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	padding [56]byte // cache line size - state pointer size = 64 - 8; created to avoid false sharing.

	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

// newAtomicBased returns a new atomic based limiter.
func newAtomicBased(rate int, opts ...Option) *atomicLimiter {
	config := buildConfig(opts)
	l := &atomicLimiter{
		perRequest: config.per / time.Duration(rate),
		maxSlack:   -1 * config.maxSlack * time.Second / time.Duration(rate),
		clock:      config.clock,
	}

	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *atomicLimiter) Take() time.Time {
	newState := state{}
	taken := false
	for !taken {
		now := t.clock.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)

		newState = state{}
		newState.last = now

		// If this is our first request, then we allow it.
		if oldState.last.IsZero() {
			taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}

		// sleepFor calculates how much time we should sleep based on
		// the perRequest budget and how long the last request took.
		// Since the request may take longer than the budget, this number
		// can get negative, and is summed across requests.
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		// We shouldn't allow sleepFor to get too negative, since it would mean that
		// a service that slowed down a lot for a short period of time would get
		// a much higher RPS following that.
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			newState.last = newState.last.Add(newState.sleepFor)
		}
		taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
	t.clock.Sleep(newState.sleepFor)
	return newState.last
}

// LimitMap limit by ip
type LimitMap struct {
	Lm map[string]Limiter
	mu sync.RWMutex
}

// NewLimitMap new limitmap
func NewLimitMap() *LimitMap {
	return &LimitMap{
		Lm: make(map[string]Limiter, 16),
	}
}

func (l *LimitMap) add(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Lm[ip] = New(defaultTimeLimitPerAct, defaultMaxSlack, defaultPer)
}

func (l *LimitMap) get(ip string) Limiter {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.Lm[ip]
	if ok {
		return v
	}
	return nil
}

//Config config for Limit func
type Config struct {
	// TimeLimitPerAct is User can send one requests per timelimitperact
	TimeLimitPerAct int
	// 单位时间
	Per time.Duration
	// IgnoreMethods skip middleware when the request method in this
	IgnoreMethods []string
	// avoid  a much high RPS
	MaxSlack time.Duration
}

// Limit the middleware of request limit
func Limit(config *Config, limiterMap *LimitMap) gin.HandlerFunc {
	ignoreMethods := config.IgnoreMethods
	timeLimitPerAct := config.TimeLimitPerAct
	per := perOption(config.Per)
	maxSlack := slackOption(config.MaxSlack)
	if ignoreMethods != nil {
		defaultNegelectMethods = ignoreMethods
	}
	if timeLimitPerAct != 0 {
		defaultTimeLimitPerAct = timeLimitPerAct
	}
	if per != 0 {
		defaultPer = per
	}
	if maxSlack != 0 {
		defaultMaxSlack = maxSlack
	}
	return func(c *gin.Context) {
		limiter := limiterMap.get(c.ClientIP())
		if limiter != nil {
			limiter.Take()
		} else {
			limiterMap.add(c.ClientIP())
		}
		c.Next()
	}
}
