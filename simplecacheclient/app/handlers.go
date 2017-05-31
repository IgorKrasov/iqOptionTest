package app

import (
	"github.com/labstack/echo"
	"github.com/iqOptionTest/simplecacheclient/cacheclient"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

func (a *App) Set(c echo.Context) error {
	request := new(cacheclient.SetBody)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result, err := a.cacheClient.Set(context.Background(), request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, result)
}

func (a *App) Get(c echo.Context) error {
	key := c.Param("key")
	result, err := a.cacheClient.Get(context.Background(), key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Unset(c echo.Context) error {
	key := c.Param("key")
	result, err := a.cacheClient.Unset(context.Background(), key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Keys(c echo.Context) error {
	result, err := a.cacheClient.Keys(context.Background())

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Rpush(c echo.Context) error {
	request := new(cacheclient.SetBody)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result, err := a.cacheClient.Rpush(context.Background(), request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, result)
}

func (a *App) Pop(c echo.Context) error {
	key := c.Param("key")
	result, err := a.cacheClient.Pop(context.Background(), key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Lgetall(c echo.Context) error {
	key := c.Param("key")
	result, err := a.cacheClient.Lgetall(context.Background(), key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Lget(c echo.Context) error {
	key := c.Param("key")
	i := c.Param("i")

	iInt, err := strconv.Atoi(i)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := a.cacheClient.Lget(context.Background(), key, iInt)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Hset(c echo.Context) error {
	request := new(cacheclient.SetBody)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result, err := a.cacheClient.Hset(context.Background(), request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, result)
}

func (a *App) Hgetall(c echo.Context) error {
	key := c.Param("key")
	result, err := a.cacheClient.Hgetall(context.Background(), key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) Hget(c echo.Context) error {
	key := c.Param("key")
	i := c.Param("i")

	result, err := a.cacheClient.Hget(context.Background(), key, i)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (a *App) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
