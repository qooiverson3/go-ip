package http

import (
	"context"
	"ipkeeper/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pingHandler struct {
	service model.PingService
}

func NewPingHandler(e *gin.Engine, s model.PingService) {
	handler := &pingHandler{
		service: s,
	}

	e.GET("/api/v1/ping/:input", handler.HealthCheck)
	e.GET("/api/v1/ping/mgo/:ip", handler.WriteResult)
	e.GET("/api/v1/ping/get/:amount", handler.GetAvailableIPs)
}

// HealthCheck 測試頁
func (h *pingHandler) HealthCheck(e *gin.Context) {
	context := h.service.HealthCheck(e.Param("input"))

	e.JSON(http.StatusOK, context)
}

func (h *pingHandler) WriteResult(e *gin.Context) {
	data := &model.IP{
		IP:    e.Param("ip"),
		Vlan:  888,
		State: false,
	}
	result, err := h.service.WriteResult(context.TODO(), *data)
	if err != nil {
		e.JSON(http.StatusInternalServerError, "寫入 mongo error")
		return
	}
	e.JSON(http.StatusOK, result)
}

func (h *pingHandler) GetAvailableIPs(e *gin.Context) {
	list := []string{
		"10.249.33",
		"10.249.34",
	}

	amount, err := strconv.Atoi(e.Param("amount"))
	if err != nil {
		e.JSON(http.StatusInternalServerError, err)
	}

	data := h.service.GetAvailableIPs(
		&model.IPResource{
			Networks: list,
			Amount:   amount,
		},
	)
	e.JSON(http.StatusOK, data)
}
