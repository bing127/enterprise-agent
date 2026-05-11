package handler

import (
	"context"

	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/bing127/enterprise-agent/module-xxx/domain/model"
	"github.com/bing127/enterprise-agent/module-xxx/service"
	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// XxxHandler HTTP 处理器。
type XxxHandler struct {
	svc service.XxxService
	log *zap.Logger
}

// NewXxxHandler 创建 XxxHandler 实例。
func NewXxxHandler(svc service.XxxService) *XxxHandler {
	return &XxxHandler{
		svc: svc,
		log: logger.L().Named("xxx"),
	}
}

// ── 请求体 DTO ────────────────────────────────────────────────────────────────

type createXxxRequest struct {
	Name string `json:"name" vd:"len($)>0"`
}

// Create godoc
// @Summary      创建 Xxx
// @Tags         xxx
// @Accept       json
// @Produce      json
// @Param        body  body      createXxxRequest  true  "请求体"
// @Success      201   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Security     BearerAuth
// @Router       /api/v1/xxx [post]
func (h *XxxHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req createXxxRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("Create: bind failed", zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	xxx := &model.Xxx{Name: req.Name}
	if err := h.svc.Create(ctx, xxx); err != nil {
		h.log.Error("Create: service error", zap.Error(err))
		response.InternalError(ctx, c, "创建失败")
		return
	}
	response.Created(ctx, c, xxx)
}

// GetByID godoc
// @Summary      获取 Xxx 详情
// @Tags         xxx
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Security     BearerAuth
// @Router       /api/v1/xxx/{id} [get]
func (h *XxxHandler) GetByID(ctx context.Context, c *app.RequestContext) {
	// TODO: 解析 id，调用 h.svc.GetByID
	response.OK(ctx, c, nil)
}

// List godoc
// @Summary      获取 Xxx 列表
// @Tags         xxx
// @Produce      json
// @Success      200  {object}  response.Response
// @Security     BearerAuth
// @Router       /api/v1/xxx [get]
func (h *XxxHandler) List(ctx context.Context, c *app.RequestContext) {
	// TODO: 解析分页参数，调用 h.svc.List
	response.OK(ctx, c, nil)
}

// Update godoc
// @Summary      更新 Xxx
// @Tags         xxx
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "ID"
// @Param        body  body      createXxxRequest  true  "请求体"
// @Success      200   {object}  response.Response
// @Security     BearerAuth
// @Router       /api/v1/xxx/{id} [put]
func (h *XxxHandler) Update(ctx context.Context, c *app.RequestContext) {
	// TODO: 解析 id + body，调用 h.svc.Update
	response.OK(ctx, c, nil)
}

// Delete godoc
// @Summary      删除 Xxx
// @Tags         xxx
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  response.Response
// @Security     BearerAuth
// @Router       /api/v1/xxx/{id} [delete]
func (h *XxxHandler) Delete(ctx context.Context, c *app.RequestContext) {
	// TODO: 解析 id，调用 h.svc.Delete
	response.OK(ctx, c, nil)
}
