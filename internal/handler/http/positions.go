package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/gin-gonic/gin"
)

type GetPositionsResult struct {
	Positions []models.Position `json:"positions"`
	Count     int               `json:"count"`
}

func (h *handler) GetPositions(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}

	res, count, err := h.service.GetAllPositions(c.Query("search"), pageNum, pageSize)
	if err != nil {
		var errMsg error
		var code int
		switch {
		case errors.Is(err, models.ErrUsernameExists):
			errMsg = models.ErrUsernameExists
			code = http.StatusBadRequest
		default:
			errMsg = models.ErrInternalServer
			code = http.StatusInternalServerError
		}
		c.JSON(code, sendResponse(-1, nil, errMsg))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, GetPositionsResult{
		Positions: res,
		Count:     count,
	}, nil))
}
