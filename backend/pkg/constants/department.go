package constants

// Department 部门枚举
type Department string

const (
	DepartmentSales          Department = "销售部"
	DepartmentRD             Department = "研发部"
	DepartmentAdministration Department = "行政与综合管理部"
	DepartmentQuality        Department = "质量管理部"
	DepartmentManagement     Department = "管理部"
)

// GetAllDepartments 获取所有部门列表
func GetAllDepartments() []Department {
	return []Department{
		DepartmentSales,
		DepartmentRD,
		DepartmentAdministration,
		DepartmentQuality,
		DepartmentManagement,
	}
}

// GetDepartmentList 获取部门字符串列表
func GetDepartmentList() []string {
	departments := GetAllDepartments()
	result := make([]string, len(departments))
	for i, dept := range departments {
		result[i] = string(dept)
	}
	return result
}

// IsValidDepartment 验证部门是否有效
func IsValidDepartment(dept string) bool {
	for _, d := range GetAllDepartments() {
		if string(d) == dept {
			return true
		}
	}
	return false
}

