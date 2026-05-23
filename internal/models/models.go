// api-structure/internal/models
package models

import "time"

type Department struct {
	ID        int          `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:200;not null;index:idx_department_parent_name,unique" json:"name"`
	ParentID  *int         `gorm:"index:idx_department_parent_name,unique" json:"parent_id,omitempty"`
	Parent    *Department  `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Children  []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Employees []Employee   `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type Employee struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	DepartmentID int        `gorm:"not null;index" json:"department_id"`
	Department   Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	FullName     string     `gorm:"size:200;not null" json:"full_name"`
	Position     string     `gorm:"size:200;not null" json:"position"`
	HiredAt      *time.Time `gorm:"type:date" json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
