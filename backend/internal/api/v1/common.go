package v1

import (
	"erb/pkg/constants"
	"erb/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetDepartments 获取部门列表
// @Summary 获取部门列表
// @Description 获取所有可用的部门列表
// @Tags 通用接口
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/common/departments [get]
func GetDepartments(c *gin.Context) {
	departments := constants.GetDepartmentList()
	response.Success(c, departments)
}

