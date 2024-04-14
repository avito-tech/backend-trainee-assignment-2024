package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func parseQueryInt64Opt(c *gin.Context, name string) *int64 {
	valStr, ok := c.GetQuery(name)
	if !ok {
		return nil
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		logrus.Infof("Cant't parse value to int64: %s", valStr)
		return nil
	}
	return &val
}

func parseQueryBoolOpt(c *gin.Context, name string) *bool {
	valStr, ok := c.GetQuery(name)
	if !ok {
		return nil
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		logrus.Infof("Cant't parse value to bool: %s", valStr)
		return nil
	}
	return &val
}
