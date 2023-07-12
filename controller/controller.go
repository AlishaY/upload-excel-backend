package controller

import (
	"upload-excel-backend/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIController struct {
	db *gorm.DB
}

func NewKPIController(db *gorm.DB) *KPIController {
	return &KPIController{
		db: db,
	}
}

func (c *KPIController) GetKPIs(ctx *gin.Context) {
	var kpis []model.KPI
	c.db.Find(&kpis)

	ctx.JSON(200, kpis)
}

func (c *KPIController) CreateKPI(ctx *gin.Context) {
	var kpi model.KPI
	if err := ctx.ShouldBindJSON(&kpi); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	kpi.Date = time.Now()

	c.db.Create(&kpi)

	ctx.JSON(201, kpi)
}