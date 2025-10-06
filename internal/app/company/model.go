package company

type Company struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	City    string `json:"city"`
}

var CreateCompanyReq struct {
	ID      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required,max=25"`
	Type    string `json:"type" binding:"required,max=25"`
	Address string `json:"address" binding:"required,max=255"`
	City    string `json:"city" binding:"required,max=100"`
}

var UpdateCompanyReq struct {
	Name    string `json:"name" binding:"required,max=25"`
	Type    string `json:"type" binding:"required,max=25"`
	Address string `json:"address" binding:"required,max=255"`
	City    string `json:"city" binding:"required,max=100"`
}
