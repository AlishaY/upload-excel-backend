package model

import "time"

type KPI struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Value    float64   `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
}