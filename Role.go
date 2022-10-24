package main

import "gorm.io/gorm"

type RoleModel struct {
	Id    uint   `gorm:"primarykey;index"`
	Value uint64 `gorm:"type:numeric;not null"`
	Name  string
}

func (role *RoleModel) TableName() string {
	return "roles"
}

func initRoles(db *gorm.DB) {
	role := new(RoleModel)

	if err := db.Where("name = user ").First(role).Error; err != nil {
		role.Value = 1 << 0
		role.Name = "user"
		tx := db.Create(role)
		tx.Commit()
	}

	if err := db.Where("name = dev ").First(role).Error; err != nil {
		role.Value = 1 << 1
		role.Name = "dev"
		tx := db.Create(role)
		tx.Commit()
	}

	if err := db.Where("name = owner ").First(role).Error; err != nil {
		role.Value = 1 << 2
		role.Name = "owner"
		tx := db.Create(role)
		tx.Commit()
	}
}
