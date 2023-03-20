package initializers

import "golang-auth/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
