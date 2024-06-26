package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp-positions-main-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type GetPositionsResult struct {
	Positions []models.Position `json:"positions"`
	Count     int               `json:"count"`
}
type GetInteviewPosition struct {
	Interviews []*models.Interview `json:"interviews"`
	Count      int                 `json:"count"`
}
type skillsReq struct {
	Skills []string `json:"skills"`
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

func (h *handler) GetPositionInterviews(c *gin.Context) {
	publicID := c.Param("position_public_id")

	err := h.service.PositionService.Exists(publicID)
	if err != nil {
		if errors.Is(err, models.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrPositionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}
	res, count, err := h.service.PositionService.GetPositionInterviews(publicID, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, GetInteviewPosition{
		Interviews: res,
		Count:      count,
	}, nil))
}

func (h *handler) GetPosition(c *gin.Context) {
	publicID := c.Param("position_public_id")

	err := h.service.PositionService.Exists(publicID)
	if err != nil {
		if errors.Is(err, models.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrPositionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	res, err := h.service.GetPosition(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, res, nil))
}

func (h *handler) CreatePosition(c *gin.Context) {
	publicID := c.GetString("public_id")

	req := &models.Position{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("failed to parse request body when creating position. %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	req.RecruiterPublicID = &publicID
	a := 0
	req.Status = &a
	res, err := h.service.PositionService.CreatePosition(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, res, nil))
}
func (h *handler) AddSkillsToPosition(c *gin.Context) {
	req := &skillsReq{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("Failed to parse request body when adding skills to position: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.Param("position_public_id") // Assuming the position public ID is in the URL path
	if err := h.service.PositionService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	err := h.service.PositionService.CreateSkillsForPosition(publicID, req.Skills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
}

func (h *handler) DeleteSkillsFromPosition(c *gin.Context) {
	req := &skillsReq{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("Failed to parse request body when deleting skills from position: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.Param("position_public_id") // Assuming the position public ID is in the URL path
	if err := h.service.PositionService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	err := h.service.DeleteSkillsFromPosition(publicID, req.Skills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}
func (h *handler) GetPositionsByCompany(c *gin.Context) {
	companyID := c.Param("company_public_id")
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}
	positions, count, err := h.service.PositionService.GetPositionsByCompany(companyID, pageNum, pageSize, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, GetPositionsResult{
		Positions: positions,
		Count:     count,
	}, nil))
}

func (h *handler) GetPositionsByRecruiter(c *gin.Context) {
	id := c.Param("recruiter_public_id")
	if err := uuid.Validate(id); err != nil {
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
	}
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}
	positions, count, err := h.service.PositionService.GetPositionsByRecruiter(id, pageNum, pageSize, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, GetPositionsResult{
		Positions: positions,
		Count:     count,
	}, nil))
}

type Questions struct {
	PositionPublicID string             `json:"position_public_id"`
	Questions        []*models.Question `json:"questions"`
}

func (h *handler) AddQuestionsToPosition(c *gin.Context) {
	id := c.Param("position_public_id")
	req := &Questions{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("Failed to parse request body when deleting skills from position: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	res, err := h.service.AddQuestionsToPosition(id, req.Questions)
	if err != nil {
		if errors.Is(err, models.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrPositionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, Questions{
		PositionPublicID: id,
		Questions:        res,
	}, nil))
}

func (h *handler) GetQuestionsToPosition(c *gin.Context) {
	id := c.Param("position_public_id")
	err := h.service.PositionService.Exists(id)
	if err != nil {
		if errors.Is(err, models.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrPositionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	res, err := h.service.GetPositionQuestions(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, Questions{
		PositionPublicID: id,
		Questions:        res,
	}, nil))
}

type PublicIDResponse struct {
	PublicID string `json:"public_id"`
}

func (h *handler) CreateInterview(c *gin.Context) {
	candidatePublicID := c.GetString("public_id")
	if c.GetString("role") != "candidate" {
		c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
		return
	}

	positionPublicID := c.Param("position_public_id")

	publicID, err := h.service.PositionService.CreateInterview(positionPublicID, candidatePublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, PublicIDResponse{
		PublicID: publicID,
	}, nil))
}

func (h *handler) DeleteQuestion(c *gin.Context) {
	publicID := c.Param("question_public_id")

	err := h.service.PositionService.DeleteQuestion(publicID)

	if err != nil {
		if errors.Is(err, models.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrQuestionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) UpdateQuestion(c *gin.Context) {
	publicID := c.Param("question_public_id")
	req := &models.Question{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("Failed to parse request body when deleting skills from position: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}
	req.PublicID = publicID
	res, err := h.service.PositionService.UpdateQuestion(req)
	if err != nil {
		if errors.Is(err, models.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrQuestionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, res, nil))
}
