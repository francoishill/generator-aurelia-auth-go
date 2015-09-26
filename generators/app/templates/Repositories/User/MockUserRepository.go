package User

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"

	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type mockRepo struct {
	Logger
}

func (m *mockRepo) GetById(id int64) User {
	m.Logger.Warn("Mock return dummy user by Id")
	return NewUser(id, "Mock Full Name", "mock@email.com")
}

func (m *mockRepo) GetByUUID(uuid interface{}) User {
	m.Logger.Warn("Mock return dummy user by UUID")
	return NewUser(-1, "Mock Full Name 2", "mock2@email.com")
}

func (m *mockRepo) GetByEmail(email string) User {
	m.Logger.Warn("Mock return dummy user by Email")
	return NewUser(-1, "Mock Full Name", email)
}

func (m *mockRepo) VerifyAndGetUserFromCredentials(email, password string) User {
	return m.GetByEmail(email)
}

func (m *mockRepo) Insert(fullName, email, rawPassword string) User {
	m.Logger.Warn("Mock just creating new dummy user")
	return NewUser(-1, fullName, email)
}

func (m *mockRepo) Update(user User) {
	m.Logger.Warn("Mock not doing anything in Update")
}

func NewMockUserRepository(logger Logger) UserRepository {
	return &mockRepo{
		logger,
	}
}
