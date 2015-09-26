package User

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"net/http"
	"strconv"

	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Encryption"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Errors"
	. "github.com/francoishill/golang-common-ddd/Interface/Storage/DbStorage"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"

	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type repo struct {
	logger        Logger
	errorsService ErrorsService
	encryption    EncryptionService
	db            DbStorage
}

type userWithFullnameAndEmail struct {
	Id       int64  `db:"id"`
	FullName string `db:"full_name"`
	Email    string `db:"email"`
}

type userWithFullnameEmailAndPassword struct {
	*userWithFullnameAndEmail
	HexHashedPassword string `db:"hashed_password"`
	HexSalt           string `db:"salt"`
}

func (t *userWithFullnameAndEmail) toUser() User {
	return NewUser(t.Id, t.FullName, t.Email)
}

func (r *repo) GetById(id int64) User {
	usr := &userWithFullnameAndEmail{}
	row := r.db.QueryRow(`
        SELECT
            u.id, u.full_name, u.email
            FROM user AS u
            WHERE u.activated_on IS NOT NULL AND banned_on IS NULL AND u.id = ?
            LIMIT ?`,
		id,
		1,
	)
	err := row.StructScan(usr)
	if err != nil {
		panic(r.errorsService.CreateClientError(http.StatusUnauthorized, "[1443259048] User missing or not activated"))
	}
	CheckError(err)

	return usr.toUser()
}

func (r *repo) GetByUUID(uuid interface{}) User {
	uuidStr := uuid.(string)
	uuidInt64, err := strconv.ParseInt(uuidStr, 10, 64)
	CheckError(err)
	return r.GetById(uuidInt64)
}

func (r *repo) GetByEmail(email string) User {
	usr := &userWithFullnameAndEmail{}
	row := r.db.QueryRow(`
        SELECT
            u.id, u.full_name, u.email
            FROM user AS u
            WHERE u.activated_on IS NOT NULL AND banned_on IS NULL AND email = ?
            LIMIT ?`,
		email,
		1,
	)
	err := row.StructScan(usr)
	if err != nil {
		panic(r.errorsService.CreateClientError(http.StatusUnauthorized, "[1443259049] User missing or not activated"))
	}
	CheckError(err)

	return usr.toUser()
}

func (r *repo) VerifyAndGetUserFromCredentials(email, guessedPassword string) User {
	usrWithPwd := &userWithFullnameEmailAndPassword{}
	row := r.db.QueryRow(`
        SELECT
            u.id, u.full_name, u.email, u.hashed_password, u.salt
            FROM user AS u
            WHERE u.activated_on IS NOT NULL AND banned_on IS NULL AND email = ?
            LIMIT ?`,
		email,
		1)
	err := row.StructScan(usrWithPwd)
	if err != nil {
		panic(r.errorsService.CreateClientError(http.StatusUnauthorized, "[1443259050] User missing or not activated"))
	}
	CheckError(err)

	matched := r.encryption.PasswordMatchHex(guessedPassword, usrWithPwd.HexSalt, usrWithPwd.HexHashedPassword)
	if !matched {
		panic(r.errorsService.CreateClientError(http.StatusUnauthorized, "User not found"))
	}

	return usrWithPwd.toUser()
}

func (r *repo) emailAlreadyExists(email string) bool {
	var cnt int
	err := r.db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).
		Scan(&cnt)
	CheckError(err)
	return cnt > 0
}

func (r *repo) Insert(fullName, email, rawPassword string) User {
	password := r.encryption.CreatePassword(rawPassword)
	result, err := r.db.Exec("INSERT INTO user (full_name, email, hashed_password, salt) VALUES (?, ?, ?, ?)",
		fullName, email, password.GetHashedPasswordHex(), password.GetSaltHex())

	if err != nil {
		if r.emailAlreadyExists(email) {
			panic(r.errorsService.CreateClientError(http.StatusBadRequest, "This email already exists"))
		}

		r.logger.Error("Unexpected insert user error: %s", err.Error())
		panic(r.errorsService.CreateClientError(http.StatusInternalServerError, ""))
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("Unexpected insert user (LastInsertId) error: %s", err.Error())
		panic(r.errorsService.CreateClientError(http.StatusInternalServerError, ""))
	}

	return NewUser(lastId, fullName, email)
}

func (r *repo) Update(user User) {
	r.db.MustExec("UPDATE user SET full_name = ? WHERE id = ?", user.FullName(), user.Id())
}

func NewDbUserRepository(logger Logger, errorsService ErrorsService, encryptionService EncryptionService, dbStorage DbStorage) UserRepository {
	return &repo{
		logger,
		errorsService,
		encryptionService,
		dbStorage,
	}
}
