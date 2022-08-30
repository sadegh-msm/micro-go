package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	ExpireTime  time.Duration `json:"expireTime"`
}

type firstRequest struct {
	URL         string `json:"url"`
	CustomShort string `json:"customShort"`
	ExpireTime  int64  `json:"expireTime"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"customShort"`
	ExpireTime     time.Duration `json:"expireTime"`
	XRAteRemain    int           `json:"xRAteRemain"`
	XRestLimitRest time.Duration `json:"xRestLimitRest"`
}

func ShortenURL(c echo.Context) error {
	firstBody := new(firstRequest)
	body := new(request)

	if err := c.Bind(&firstBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "cannot pars JSON"})
	}

	body.URL = firstBody.URL
	body.CustomShort = firstBody.CustomShort
	body.ExpireTime = time.Duration(firstBody.ExpireTime)

	r2 := CreateClients(1)
	defer r2.Close()
	val, err := r2.Get(Context, c.RealIP()).Result()

	if err == redis.Nil {
		_ = r2.Set(Context, c.RealIP(), 10, 30*time.Minute).Err()
	} else {
		val, _ = r2.Get(Context, c.RealIP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(Context, c.RealIP()).Result()

			return c.JSON(http.StatusServiceUnavailable, echo.Map{
				"error":           "rate limit exceeded",
				"limit_time_left": limit,
			})

		}
	}

	if !govalidator.IsURL(body.URL) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL is not correct"})
	}

	if !RemoveDomainError(body.URL) {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{"error": "error in finding domain"})
	}
	body.URL = EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	rdb := CreateClients(0)
	defer rdb.Close()

	val, _ = rdb.Get(Context, id).Result()

	if val != "" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "url is already used",
		})
	}

	if body.ExpireTime == 0 {
		body.ExpireTime = 24
	}

	er := rdb.Set(Context, id, body.URL, body.ExpireTime*60*time.Minute)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "unable to connect to server",
		})
	}

	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		ExpireTime:     body.ExpireTime,
		XRAteRemain:    10,
		XRestLimitRest: 30,
	}

	r2.Decr(Context, c.RealIP())

	val, _ = r2.Get(Context, c.RealIP()).Result()
	resp.XRAteRemain, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(Context, c.RealIP()).Result()
	resp.XRestLimitRest = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = "urlshortner-service" + "/" + id

	return c.JSON(http.StatusOK, resp)
}
