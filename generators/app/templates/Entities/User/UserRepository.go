package User

type UserRepository interface {
	GetById(id int64) User
	GetByUUID(uuid interface{}) User
	GetByEmail(email string) User
	VerifyAndGetUserFromCredentials(email, guessedPassword string) User

	Insert(fullName, email, rawPassword string) User
	Update(User)
}
