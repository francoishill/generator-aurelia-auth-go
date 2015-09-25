package User

import (
	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type mockRepo struct{}

func (m *mockRepo) GetById(id int64) User {
	return NewUser(id, "Mock Full Name", "mock@email.com")
}

func (m *mockRepo) GetByUUID(uuid interface{}) User {
	return NewUser(-1, "Mock Full Name 2", "mock2@email.com")
}

func (m *mockRepo) GetByEmail(email string) User {
	return NewUser(-1, "Mock Full Name", email)
}

func (m *mockRepo) VerifyAndGetUserFromCredentials(email, password string) User {
	return m.GetByEmail(email)
}

func (r *mockRepo) Insert(fullName, email, rawPassword string) User {
	return NewUser(-1, fullName, email)
}

func (r *mockRepo) Update(user User) {
	//Do nothing
}

func NewMockUserRepository() UserRepository {
	return &mockRepo{}
}
