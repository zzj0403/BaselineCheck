package compliance

import (
	"BaselineCheck/client/baselinelinux"
	"encoding/json"
	"errors"
	"strconv"

	"dario.cat/mergo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}
func (h *Handler) RegisterComplianceResult(ctx *gin.Context) {
	req := new(RegisterRequest)
	details := new(ComplianceDetails)

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 初始化 ComplianceResult
	result := &ComplianceResult{
		Hostname:      req.BaseInfo.HostName,
		BaselineCount: req.LengthComplianceInfo(),
		IP:            req.BaseInfo.LanIp,
	}
	// 判断是否已经存在
	existingResult, err := h.repo.GetComplianceResultByHostname(result.Hostname)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新记录
			if err := h.repo.CreateComplianceResult(result); err != nil {
				ctx.JSON(500, gin.H{"error": "创建合规结果失败"})
				return
			}
		} else {
			ctx.JSON(500, gin.H{"error": "查询合规结果失败"})
			return
		}
	} else {
		// 确保 existingResult 不为 nil
		if existingResult != nil {
			// 更新已存在的记录
			if err := mergo.MergeWithOverwrite(result, existingResult); err != nil {
				ctx.JSON(500, gin.H{"error": "合并合规结果失败"})
				return
			}
			result.ID = existingResult.ID // 确保使用已有记录的 ID
			if err := h.repo.UpdateComplianceResultByHostname(result); err != nil {
				ctx.JSON(500, gin.H{"error": "更新合规结果失败"})
				return
			}
		}
	}

	// 创建合规详情
	details.ResultId = result.ID
	details.Details, _ = json.Marshal(req)
	if err := h.repo.CreateComplianceDetails(details); err != nil {
		ctx.JSON(500, gin.H{"error": "创建合规详情失败"})
		return
	}

	// 返回成功响应
	ctx.JSON(201, result) // 返回201 Created和新创建的结果
}

// 展示合规结果
func (h *Handler) ShowComplianceHostList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))           // 默认第一页
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10")) // 默认每页10条

	// 将分页参数转换为合适的值
	Host, err := h.repo.GetComplianceHostList(page, pageSize)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "查询列表结果失败"})
		return
	}
	var res ComplianceHostnameListRespone
	res.List = Host
	res.Total = len(Host)
	res.Page.Page = page
	res.Page.PageSize = pageSize
	ctx.JSON(200, res)
	// 返回结果
}

// 展示合规结果
func (h *Handler) ShowComplianceResultList(ctx *gin.Context) {
	hostID := ctx.DefaultQuery("host_id", "") // 获取 host_id 参数
	// 查询数据库
	result, err := h.repo.GetComplianceDetailByhostId(hostID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "查询合规结果失败"})
		return
	}
	ctx.JSON(200, result)

}
func (h *Handler) ShowComplianceDetails(ctx *gin.Context) {
	// 获取url中的id
	id := ctx.Param("id")
	result_id := ctx.Param("result_id")

	// 查询数据库
	result, err := h.repo.GetComplianceDetailsByResultIdandDetailId(result_id, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "查询合规详情失败"})
		return
	}
	var res baselinelinux.Result
	if err := json.Unmarshal(result.Details, &res); err != nil {
		// JSON解码失败，返回400错误
		ctx.JSON(400, gin.H{"error": "合规详情解析失败"})
		return
	}

	ctx.HTML(200, "detail.html", gin.H{"Report": res})
}
