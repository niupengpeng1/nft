package handler

import (
	"net/http"
	"strconv"

	"dapp/internal/models"
	"dapp/internal/repository"
	"dapp/internal/service"
	"dapp/pkg/logger"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService *service.EventListenerService
	eventRepo    *repository.EventRepository
}

func NewEventHandler(eventService *service.EventListenerService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		eventRepo:    repository.NewEventRepository(),
	}
}

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Status returns the status of the event listener
// @Summary Get service status
// @Description Returns the current status of the event listener service
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/status [get]
func (h *EventHandler) Status(c *gin.Context) {
	status := h.eventService.GetStatus()
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "Success",
		Data:    status,
	})
}

// GetEvents returns all events with optional filtering
// @Summary Get contract events
// @Description Query contract events with filters (block_number, contract_address, or tx_hash)
// @Tags Events
// @Accept json
// @Produce json
// @Param block_number query int false "Block number"
// @Param contract_address query string false "Contract address"
// @Param tx_hash query string false "Transaction hash"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	blockNumberStr := c.Query("block_number")
	contractAddr := c.Query("contract_address")
	txHash := c.Query("tx_hash")

	var events interface{}
	var err error

	if txHash != "" {
		events, err = h.eventRepo.GetEventByTxHash(txHash)
	} else if blockNumberStr != "" {
		blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "Invalid block number format",
			})
			return
		}
		events, err = h.eventRepo.GetEventsByBlockNumber(blockNumber)
	} else if contractAddr != "" {
		events, err = h.eventRepo.GetEventsByContractAddress(contractAddr)
	} else {
		// Return all events (could add pagination here)
		logger.Warn("Fetching all events without filter - consider adding pagination")
		events = []*models.ContractEvent{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to fetch events: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "Success",
		Data:    events,
	})
}

// StartListening starts the event listener
// @Summary Start event listener
// @Description Start the blockchain event listening service
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/start [post]
func (h *EventHandler) StartListening(c *gin.Context) {
	go func() {
		if err := h.eventService.Start(); err != nil {
			logger.Error("Failed to start event listener: %v", err)
		}
	}()

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "Event listener starting",
	})
}

// StopListening stops the event listener
// @Summary Stop event listener
// @Description Stop the blockchain event listening service
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/stop [post]
func (h *EventHandler) StopListening(c *gin.Context) {
	h.eventService.Stop()

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "Event listener stopped",
	})
}

// GetEventByID gets a single event by database ID
// @Summary Get event by ID
// @Description Get a specific event by database ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /api/events/:id [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Event ID is required",
		})
		return
	}

	// Note: This would require implementing GetEventByID in repository
	c.JSON(http.StatusNotImplemented, Response{
		Code:    501,
		Message: "Not implemented yet",
	})
}

func (h *EventHandler) QueryAllLogs(c *gin.Context) {

	queryDTO := models.LogDTO{}
	if err := c.ShouldBind(&queryDTO); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	height, err := h.eventRepo.GetBlockHeight()
	logger.Info("Current block height: %d", height)

	resukt, err := h.eventRepo.QueryAllLogs(&queryDTO)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to fetch events: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: resukt,
	})
}
