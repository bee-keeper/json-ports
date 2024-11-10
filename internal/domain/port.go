package domain

import (
	"gorm.io/datatypes"
)

// Port entity
type Port struct {
	ID          uint           `gorm:"primaryKey"`
	Unloc       string         `gorm:"unique;not null"`
	Name        string         `json:"name"`
	City        string         `json:"city"`
	Country     string         `json:"country"`
	Alias       datatypes.JSON `json:"alias,omitempty"`
	Regions     datatypes.JSON `json:"regions,omitempty"`
	Coordinates datatypes.JSON `json:"coordinates"`
	Province    string         `json:"province"`
	Timezone    string         `json:"timezone"`
	Unlocs      datatypes.JSON `json:"unlocs"`
	Code        string         `json:"code"`
}
