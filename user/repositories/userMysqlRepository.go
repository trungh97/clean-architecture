package repositories

import (
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/user/entities"
)

type UserMysqlRepository struct {
	db database.Database
}

func NewUserMysqlRepository(db database.Database) UserRepository {
	return &UserMysqlRepository{
		db: db,
	}
}

func (r *UserMysqlRepository) GetUser(email string) (entities.User, error) {
	var user entities.User

	err := r.db.GetDb().Where("email = ?", email).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.User{}, errors.New("invalid username or password")
		}
		log.Errorf("failed to get user from the database: %v", err)
		return entities.User{}, errors.New("internal server error")
	}

	return user, err

}
