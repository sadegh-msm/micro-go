package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"shortnerApp/data"
	"strconv"
	"time"
)

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"customShort"`
	ExpireTime     time.Duration `json:"expireTime"`
	XRAteRemain    int           `json:"xRAteRemain"`
	XRestLimitRest time.Duration `json:"xRestLimitRest"`
}

func ShortenURL(body data.Request, realIp string) (int, interface{}) {
	r2 := data.CreateClients(1)
	defer r2.Close()
	val, err := r2.Get(data.Context, realIp).Result()

	if err == redis.Nil {
		_ = r2.Set(data.Context, realIp, 10, 30*time.Minute).Err()
	} else {
		val, _ = r2.Get(data.Context, realIp).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(data.Context, realIp).Result()

			return 503, echo.Map{
				"error":           "rate limit exceeded",
				"limit_time_left": limit,
			}

		}
	}

	if !govalidator.IsURL(body.URL) {
		return 400, echo.Map{"error": "URL is not correct"}
	}

	if !RemoveDomainError(body.URL) {
		return 503, echo.Map{"error": "error in finding domain"}
	}
	body.URL = EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	rdb := data.CreateClients(0)
	defer rdb.Close()

	val, _ = rdb.Get(data.Context, id).Result()

	if val != "" {
		return 403, echo.Map{
			"error": "url is already used",
		}
	}

	if body.ExpireTime == 0 {
		body.ExpireTime = 24
	}

	er := rdb.Set(data.Context, id, body.URL, body.ExpireTime*60*time.Minute)
	if er != nil {
		return 500, echo.Map{
			"error": "unable to connect to server",
		}
	}

	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		ExpireTime:     body.ExpireTime,
		XRAteRemain:    10,
		XRestLimitRest: 30,
	}

	r2.Decr(data.Context, realIp)

	val, _ = r2.Get(data.Context, realIp).Result()
	resp.XRAteRemain, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(data.Context, realIp).Result()
	resp.XRestLimitRest = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = "urlshortner-service" + "/" + id

	return 200, resp
}

func ShortenUrlEcho(c echo.Context) error {
	body := data.Request{}

	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "cannot pars JSON"})
	}

	status, res := ShortenURL(body, c.RealIP())

	switch status {
	case 500:
		return c.JSON(status, res)

	case 503:
		return c.JSON(status, res)

	case 400:
		return c.JSON(status, res)

	case 403:
		return c.JSON(status, res)

	}

	return c.JSON(status, res)
}
