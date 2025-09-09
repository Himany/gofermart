package worker

import (
	"net/http"
	"time"

	"github.com/Himany/gofermart/internal/accrual"
	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/models"
	"go.uber.org/zap"
)

type Repo interface {
	UpdateOrderStatus(number string, status string, accrual *float64) error
	AddAccrual(userID int, orderNumber string, amount float64) error
}

type Job struct {
	UserID   int
	Number   string
	Attempts int
}

type Pool struct {
	repo    Repo
	cli     *accrual.Client
	workers int
	jobs    chan Job
	quit    chan struct{}
}

func NewPool(repo Repo, cli *accrual.Client, workers int, buffer int) *Pool {
	if workers <= 0 {
		workers = 4
	}
	if buffer <= 0 {
		buffer = 1024
	}
	return &Pool{
		repo:    repo,
		cli:     cli,
		workers: workers,
		jobs:    make(chan Job, buffer),
		quit:    make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		go p.worker()
	}
}

func (p *Pool) Stop() {
	close(p.quit)
}

func (p *Pool) Submit(userID int, number string) {
	select {
	case p.jobs <- Job{UserID: userID, Number: number}:
	default:
		logger.Log.Warn("Submit (worker queue full, drop job)", zap.String("number", number))
	}
}

func (p *Pool) worker() {
	for {
		select {
		case <-p.quit:
			return
		case job := <-p.jobs:
			p.process(job)
		}
	}
}

func (p *Pool) process(job Job) {
	info, code, retryAfter, err := p.cli.GetOrder(job.Number)

	//429 Retry-After
	if code == http.StatusTooManyRequests {
		p.requeue(job, retryAfterOrDefault(retryAfter))
		return
	}
	//204 заказ ещё не зарегистрирован
	if code == http.StatusNoContent {
		p.requeue(job, nextDelay(job.Attempts))
		return
	}
	//Ошибка сети
	if err != nil {
		logger.Log.Warn("accrual request failed",
			zap.String("number", job.Number),
			zap.Int("attempts", job.Attempts),
			zap.Error(err),
		)
		p.requeue(job, nextDelay(job.Attempts))
		return
	}

	status := mapExternal(info.Status)
	var accrualVal *float64
	if info.Accrual != nil {
		accrualVal = info.Accrual
	}

	if err := p.repo.UpdateOrderStatus(job.Number, status, accrualVal); err != nil {
		logger.Log.Error("process (UpdateOrderStatus)", zap.String("number", job.Number), zap.Error(err))
		p.requeue(job, nextDelay(job.Attempts))
		return
	}

	if status == "INVALID" {
		return
	}
	if status == "PROCESSED" {
		if accrualVal != nil && *accrualVal > 0 {
			if err := p.repo.AddAccrual(job.UserID, job.Number, *accrualVal); err != nil {
				logger.Log.Error("process (AddAccrual)", zap.String("number", job.Number), zap.Error(err))
			}
		}
		return
	}

	p.requeue(job, nextDelay(job.Attempts))
}

func (p *Pool) requeue(job Job, delay time.Duration) {
	job.Attempts++
	time.AfterFunc(delay, func() {
		select {
		case p.jobs <- job:
		default:
			logger.Log.Warn("worker queue full on requeue", zap.String("number", job.Number))
		}
	})
}

func nextDelay(attempts int) time.Duration {
	if attempts < 0 {
		attempts = 0
	}

	maxDelay := 60 * time.Second
	delay := time.Duration(5 + (5 * attempts))

	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

func retryAfterOrDefault(delay time.Duration) time.Duration {
	if delay <= 0 {
		return 60 * time.Second
	}
	if delay > 5*time.Minute {
		return 5 * time.Minute
	}
	return delay
}

func mapExternal(status models.StatusOrder) string {
	switch status {
	case models.StatusRegistered:
		return "NEW"
	case models.StatusProcessing:
		return "PROCESSING"
	case models.StatusInvalid:
		return "INVALID"
	case models.StatusProcessed:
		return "PROCESSED"
	default:
		return "NEW"
	}
}
