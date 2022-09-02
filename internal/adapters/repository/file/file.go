package file

import (
	"auth/internal/domain/models"
	"context"
	"encoding/json"
	"os"
)

type FileRepo struct {
	File *os.File
}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func (r *FileRepo) GetUser(ctx context.Context, login string) (models.User, error) {
	user := models.User{}
	file, err := os.ReadFile(os.Getenv("CONFIG_PATH") + "user_data.json")
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(file, &user)
	if err != nil {
		return user, err
	}

	if user.Login != login {
		return user, models.ErrUserNotFound
	}

	//err = hash.ComparePasswordAndHash(user.Password, password)
	//if err != nil {
	//	return nil, ErrUserNotFound
	//}

	return user, nil
}

func (r *FileRepo) CreateUser(ctx context.Context, user models.User) error {
	var err error
	r.File, err = os.OpenFile(os.Getenv("CONFIG_PATH")+"user_data.json", os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer r.File.Close()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = r.File.Write(data)
	if err != nil {
		return err
	}
	return nil
}

//func (r *UserRepo) GetAll() (Users, error){
//	users := make(Users)
//	file, err := os.ReadFile(os.Getenv("CONFIG_PATH") + "user_data.json")
//	if err != nil {
//		return nil, err
//	}
//	err = json.Unmarshal(file, &users)
//	if err != nil {
//		return nil, err
//	}
//	return users, nil
//}
