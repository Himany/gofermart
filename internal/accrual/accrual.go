package accrual

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Himany/gofermart/internal/models"
	"github.com/Himany/gofermart/internal/utils"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	baseURL string
	HTTP    *resty.Client
}

func New(baseURL string, timeout time.Duration) *Client {
	c := resty.New().
		SetTimeout(timeout).
		SetHeader("Accept", "application/json")

	return &Client{
		baseURL: baseURL,
		HTTP:    c,
	}
}

func (c *Client) GetOrder(number string) (models.OrderInfo, int, time.Duration, error) {
	var out models.OrderInfo
	resp, err := c.HTTP.R().
		SetResult(&out).
		Get(fmt.Sprintf("%s/api/orders/%s", c.baseURL, number))

	if err != nil {
		return models.OrderInfo{}, 0, 0, err
	}

	code := resp.StatusCode()
	switch code {
	case http.StatusOK:
		return out, code, 0, nil
	case http.StatusNoContent:
		return models.OrderInfo{}, code, 0, nil
	case http.StatusTooManyRequests:
		ra := utils.ParseRetryAfter(resp.Header().Get("Retry-After"))
		if ra == 0 {
			ra = 60 * time.Second
		}
		return models.OrderInfo{}, code, ra, fmt.Errorf("too many requests")
	default:
		return models.OrderInfo{}, code, 0, fmt.Errorf("unexpected status: %d", code)
	}
}
