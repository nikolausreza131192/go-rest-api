package repository

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	userConstants "github.com/nikolausreza131192/pos/constants"
	"github.com/nikolausreza131192/pos/entity"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo define the interface for user repository
type UserRepo interface {
	GetByUsername(username string) (entity.User, error)
	GetUserPassword(username string) (string, error)
	CreateUser(user entity.User) (string, error)
}

type userRepo struct {
	db *sqlx.DB
}

// UserRepoParam will be used as repository parameter
type UserRepoParam struct {
	DB *sqlx.DB
}

// NewUser initialize item repository
func NewUser(param UserRepoParam) UserRepo {
	r := &userRepo{
		db: param.DB,
	}

	return r
}

// GetByUsername return user by specified username
func (r *userRepo) GetByUsername(username string) (entity.User, error) {
	result := entity.User{}
	row := r.db.QueryRowx(`
		SELECT 
			id, 
			nama, 
			username, 
			COALESCE(email, '') as email, 
			role, 
			status, 
			created_by, 
			updated_by, 
			created_at, 
			updated_at 
		FROM m_user 
		WHERE username = ?
	`, username)
	if err := row.StructScan(&result); err != nil {
		log.Println("GetByUsername Error struct scan", username, err)
		return result, err
	}

	return result, nil
}

// GetUserPassword return password of specified user
func (r *userRepo) GetUserPassword(username string) (string, error) {
	password := ""
	row := r.db.QueryRow(`
		SELECT password 
		FROM m_user 
		WHERE username = ?
	`, username)
	if err := row.Scan(&password); err == sql.ErrNoRows {
		return password, errors.New("User is not found")
	} else if err != nil {
		log.Printf("func GetUserPassword Error query. Username: %s. Error: %s", username, err)
		return password, err
	}

	return password, nil
}

// CreateUser will add new row to user table
func (r *userRepo) CreateUser(user entity.User) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("func CreateUser Error begin transaction. User: %+v. Error: %s", user, err)
		return "", err
	}

	// Generate random password
	userPassword := generateRandomString()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

	query := `
		INSERT INTO m_user(nama, username, password, email, role, status, created_by, updated_by, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(
		query,
		user.Name,
		user.Username,
		string(hashedPassword),
		user.Email,
		user.Role,
		userConstants.UserDefaultStatus,
		user.CreatedBy,
		user.UpdatedBy,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Printf("func CreateUser Error execute query. User: %+v. Error: %s", user, err)
		return "", err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("func CreateUser Error commit transaction. User: %+v. Error: %s", user, err)
		return "", err
	}

	return userPassword, nil
}

// generateRandomString is used to generate 6 random alphabet letters
func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
