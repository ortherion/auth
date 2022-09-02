package test

import (
	"auth/internal/domain/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	strId    = "62c48a42e31ecb2af1d5d1c4"
	id, _    = primitive.ObjectIDFromHex(strId)
	userName = "test123"
	user     = models.User{
		ID:           id,
		Login:        userName,
		Password:     "$2a$04$3Fwej2KBe58nKVdo0n9mqugGQrEdwzvJqF1JBUgDI3TLLzntYOW96",
		Email:        "test123@ya.ru",
		FirstName:    "test",
		LastName:     "123",
		CreationDate: 0,
	}
)

func (tc *TestContainersSuite) TestUserRepoCreateSuccess() {
	newUser := &models.User{
		Login:        "user123",
		Password:     user.Password,
		Email:        "newUser123.ya.ru",
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CreationDate: uint64(time.Now().Unix()),
	}

	err := tc.userRepo.Insert(context.Background(), newUser)
	tc.NoError(err)

	users, err := tc.userRepo.GetAll(context.Background())

	tc.NoError(err)
	tc.NotNil(users, "users must be not nil")
	tc.Condition(func() bool {
		return len(users) > 0
	})
}

func (tc *TestContainersSuite) TestUserRepoCreateDuplicateName() {
	err := tc.userRepo.Insert(context.Background(), &user)

	tc.NotNil(err)
}

func (tc *TestContainersSuite) TestUserRepoGetAllSuccess() {
	users, err := tc.userRepo.GetAll(context.Background())

	tc.NoError(err)
	tc.NotNil(users, "users must be not nil")
	tc.Condition(func() bool {
		return len(users) > 0
	})
}

func (tc *TestContainersSuite) TestUserRepoGetByNameSuccess() {
	dbUser, err := tc.userRepo.GetByName(context.Background(), userName)

	tc.NoError(err)
	tc.NotNil(user, "user must be not nil")
	tc.Equal(user, *dbUser)
}

func (tc *TestContainersSuite) TestUserRepoGetSuccess() {
	tc.Suite.T().SkipNow()
	users, err := tc.userRepo.GetAll(context.Background())

	tc.NoError(err)
	tc.NotNil(users, "users must be not nil")
	tc.Condition(func() bool {
		return len(users) > 0
	})

	i, _ := primitive.ObjectIDFromHex(strId)

	tc.Require().Equal(i, users[0].ID)

	dbUser, err := tc.userRepo.Get(context.Background(), strId)

	tc.NoError(err)
	tc.NotNil(user, "user must be not nil")
	tc.Equal(user, *dbUser)
}

func (tc *TestContainersSuite) TestUserRepoGetInvalidId() {
	u, err := tc.userRepo.Get(context.Background(), "lkfjsdlajfds")

	tc.Nil(u)
	tc.NotNil(err)
}

func (tc *TestContainersSuite) TestUserRepoUpdateSuccess() {
	newUser := &models.User{
		ID:        user.ID,
		Login:     user.Login,
		Email:     user.Email + "x",
		FirstName: user.FirstName + "x",
		LastName:  user.LastName + "x",
	}

	err := tc.userRepo.Update(context.Background(), newUser)
	tc.Nil(err)

	dbUser, err := tc.userRepo.GetByName(context.Background(), user.Login)

	tc.Equal(dbUser.FirstName, newUser.FirstName)
	tc.Equal(dbUser.LastName, newUser.LastName)
	tc.Equal(dbUser.Email, newUser.Email)

	tc.NotEqualValues(dbUser.FirstName, user.FirstName)
	tc.NotEqualValues(dbUser.LastName, user.LastName)
	tc.NotEqualValues(dbUser.Email, user.Email)
}
