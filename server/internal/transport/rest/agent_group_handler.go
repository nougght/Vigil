package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/service/agent_groups"
	"github.com/nougght/monitoring-system/server/internal/transport/dto/mapper"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
	sharedUtil "github.com/nougght/monitoring-system/shared/go/util"
)

type AgentGroupHandler struct {
	groupService *agent_groups.AgentGroupsService
}

func newAgentGroupHandler(groupService *agent_groups.AgentGroupsService) *AgentGroupHandler {
	if groupService == nil {
		log.Panicf("agent group handler params required")
	}
	return &AgentGroupHandler{
		groupService: groupService,
	}
}

func (h *AgentGroupHandler) RegisterRoutes(r *gin.RouterGroup) {
	group := r.Group("/groups")

	group.POST("", h.CreateGroup)
	group.GET("/:id", h.GetGroupByID)
	group.GET("", h.GetAllGroups)
	group.PATCH("/:id", h.UpdateGroup)
	group.DELETE("/:id", h.DeleteGroup)

	group.POST("/:id/agents", h.MoveAgentsToGroup)
	r.PUT("/agents/:id/group", h.MoveAgentToGroup)
}

// CreateGroup godoc
// @Id createAgentGroup
// @Tags agent-groups
// @Summary Create new agent group
// @Accept json
// @Produce json
// @Param request body dto.CreateAgentGroupBody true "Create agent group body"
// @Success 200 {object} dto.AgentGroupDTO
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups [post]
func (h *AgentGroupHandler) CreateGroup(c *gin.Context) {
	ctx := c.Request.Context()
	var body dto.CreateAgentGroupBody

	err := c.ShouldBindJSON(&body)
	if err != nil {
		handleInvalidRequestBody(c, err)
		return
	}

	res, err := h.groupService.CreateGroup(ctx, mapper.CreateAgentGroupInputFromDTO(&body))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.CreateAgentGroupResultToDTO(res))
}

// GetAllGroups godoc
// @Id getAllAgentGroups
// @Tags agent-groups
// @Summary Get all agent groups
// @Produce json
// @Success 200 {array} dto.AgentGroupDTO
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups [get]
func (h *AgentGroupHandler) GetAllGroups(c *gin.Context) {
	groups, err := h.groupService.GetAllGroups(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, sharedUtil.Map(groups, mapper.AgentGroupToDTO))
}

// GetGroupByID godoc
// @Id getAgentGroupByID
// @Tags agent-groups
// @Summary Get agent group by ID
// @Produce json
// @Param id path string true "Agent group ID" format(uuid)
// @Success 200 {object} dto.AgentGroupDTO
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups/{id} [get]
func (h *AgentGroupHandler) GetGroupByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := getPathID(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	group, err := h.groupService.GetGroupByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.AgentGroupToDTO(group))
}

// UpdateGroup godoc
// @Id updateAgentGroup
// @Tags agent-groups
// @Summary Update agent group
// @Accept json
// @Param id path string true "Agent group ID" format(uuid)
// @Param request body dto.UpdateAgentGroupBody true "Update agent group body"
// @Success 200
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups/{id} [patch]
func (h *AgentGroupHandler) UpdateGroup(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := getPathID(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var body dto.UpdateAgentGroupBody
	err = c.ShouldBindJSON(&body)
	if err != nil {
		handleInvalidRequestBody(c, err)
		return
	}
	if id != body.ID {
		handleError(c, ErrorDiffParamAndBodyID)
		return
	}

	err = h.groupService.UpdateGroup(ctx, mapper.UpdateAgentGroupInputFromDTO(&body))
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteGroup godoc
// @Id deleteAgentGroup
// @Tags agent-groups
// @Summary Delete agent group by ID
// @Param id path string true "Agent group ID" format(uuid)
// @Success 200
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups/{id} [delete]
func (h *AgentGroupHandler) DeleteGroup(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := getPathID(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	err = h.groupService.DeleteGroupByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// MoveAgentsToGroup godoc
// @Id moveAgentsToGroup
// @Tags agent-groups
// @Summary Move agents to group by agent IDs
// @Param id path string true "Group ID" format(uuid)
// @Param request body dto.MoveAgentsToGroupRequestBody true "Move agents to group request body"
// @Success 200
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /groups/{id}/agents [POST]
func (h *AgentGroupHandler) MoveAgentsToGroup(c *gin.Context) {
	ctx := c.Request.Context()
	groupID, err := getPathID(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var body dto.MoveAgentsToGroupRequestBody
	err = c.ShouldBindJSON(&body)
	if err != nil {
		handleInvalidRequestBody(c, err)
		return
	}

	_, err = h.groupService.MoveAgentsToGroup(ctx, body.AgentIDs, groupID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// MoveAgentToGroup godoc
// @Id moveAgentToGroup
// @Tags agent-groups
// @Summary Move agent to group
// @Param id path string true "Agent ID" format(uuid)
// @Param groupId query string true "Group ID to move agent"
// @Success 200
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router /agents/{id}/group [PUT]
func (h *AgentGroupHandler) MoveAgentToGroup(c *gin.Context) {
	ctx := c.Request.Context()
	agentID, err := getPathID(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var params dto.MoveAgentToGroupQueryParams
	err = c.ShouldBindQuery(&params)
	if err != nil {
		handleError(c, fmt.Errorf("%w: %w", err, model.ErrBadRequest))
		return
	}
	groupID, err := uuid.Parse(params.GroupID)
	if err != nil {
		handleError(c, fmt.Errorf("%w: %w", err, model.ErrBadRequest))
		return
	}

	err = h.groupService.MoveAgentToGroup(ctx, agentID, groupID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
