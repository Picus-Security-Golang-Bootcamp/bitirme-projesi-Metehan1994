package user

import (
	"errors"
	"fmt"

	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Migration() {
	u.db.AutoMigrate(&models.User{})
}

func (u *UserRepository) InsertSampleData(user *models.User) {
	result := u.db.Unscoped().Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	if result.Error != nil {
		zap.L().Fatal("Cannot insert data into DB") //Check Error
	}
}

//GetUserList provides a list of user
func (u *UserRepository) GetUserList() []models.User {
	zap.L().Debug("user.repo.GetuserList")
	var users []models.User
	u.db.Find(&users)

	return users
}

//GetUser finds the user matched with given email and password
func (u *UserRepository) GetUser(email, password string) (models.User, bool) {
	zap.L().Debug("user.repo.Getuser")
	var user models.User
	u.db.Where("email = ?", email).Find(&user)
	if user.Password == password && user.Email == email {
		return user, true
	}
	return user, false
}

//GetUserByEmail finds the user by email
func (u *UserRepository) GetUserByEmail(email string) *models.User {
	zap.L().Debug("user.repo.GetUserByEmail", zap.Reflect("email", email))
	var user models.User
	results := u.db.Where("email = ?", email).Find(&user)
	if results.Error != nil {
		zap.L().Error(results.Error.Error())
	}
	return &user
}

//GetUserByEmail finds the user by id
func (u *UserRepository) GetUserByID(Id string) (*models.User, error) {
	zap.L().Debug("user.repo.GetUserByID", zap.Reflect("Id", Id))
	//uuid, _ := uuid.FromBytes([]byte(Id))
	fmt.Println(Id)
	var user models.User
	userUUID := uuid.MustParse(Id)
	results := u.db.Where("id = ?", userUUID).Find(&user)
	if results.Error != nil {
		return nil, results.Error
	}
	if user.ID == uuid.Nil {
		return nil, errors.New("user not found. You need to login to the system")
	}
	return &user, nil
}

//GetUserByEmail finds the user by username
func (u *UserRepository) GetUserByUsername(username string) *models.User {
	zap.L().Debug("user.repo.GetUserByUsername", zap.Reflect("username", username))
	var user models.User
	results := u.db.Where("username = ?", username).Find(&user)
	if results.Error != nil {
		zap.L().Error(results.Error.Error())
	}
	return &user
}

//CreateNewUser creates a new user based on the sign up process.
func (u *UserRepository) CreateNewUser(user models.User) (*models.User, error) {
	zap.L().Debug("user.repo.CreateNewUser", zap.Reflect("user", user))
	result := u.db.Where("email = ?", user.Email).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
