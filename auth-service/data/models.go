package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Set timeout for DB Operation
const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the 'Models' struct,
// which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the struct containing many structs, and any model(struct type)
// that is included as a member in this will be available to us throughout
// the application, anywhere that the 'app' variable is used, provided that
// the model is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, 
	updated_at from users order by last_name`

	// Query DB
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	// Read 'user' info from each row, and add to 'users' slice.
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning: ", err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, 
	updated_at from users where email = $1`

	// Query DB
	row := db.QueryRowContext(ctx, query, email)

	var user User
	// Read 'user' info.
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error scanning: ", err)
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, 
	updated_at from users where id = $1`

	// Query DB
	row := db.QueryRowContext(ctx, query, id)

	var user User
	// Read 'user' info.
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error scanning: ", err)
		return nil, err
	}

	return &user, nil
}

// Update updates user info in the database, using the information stored
// in the receiver 'u'
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `update user set
			  email = $1,
			  first_name = $2, 
			  last_name = $3, 
			  user_active = $4,
			  updated_at = $5
			  where id = $6`

	// Update the user info in DB
	_, err := db.ExecContext(ctx, query,
		u.Email, u.FirstName, u.LastName, u.Active, time.Now(), u.ID)

	if err != nil {
		log.Println("Error updating: ", err)
		return err
	}

	return nil
}

// Delete deletes the given user from database, by 'u.ID'
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`

	// Delete the user info in DB
	_, err := db.ExecContext(ctx, query, u.ID)

	if err != nil {
		log.Println("Error deleting: ", err)
		return err
	}

	return nil
}

// DeleteByID deletes the given user from database, by 'id'
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`

	// Delete the user info in DB
	_, err := db.ExecContext(ctx, query, id)

	if err != nil {
		log.Println("Error deleting: ", err)
		return err
	}

	return nil
}

// Insert inserts a new user into database, and returns the ID of it
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("Error generating hashed password: ", err)
		return 0, err
	}

	var newID int
	query := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7) returning id`

	// Insert the user info in DB
	err = db.QueryRowContext(ctx, query,
		user.Email, user.FirstName, user.LastName, hashedPassword, user.Active, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		log.Println("Error inserting: ", err)
		return 0, err
	}

	return newID, nil
}

// ResetPassword is used to update the user's password in DB
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println("Error generating hashed password: ", err)
		return err
	}

	query := `update users set password = $1 where id = $2`

	// Update user's info in DB
	_, err = db.ExecContext(ctx, query, hashedPassword, u.ID)

	if err != nil {
		log.Println("Error updating password: ", err)
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare the hash of a user
// entered password with the pasword hash we have stored for that user in
// the DB. If both the hash matches, then return true, else return false.
func (u *User) PasswordMatches(enteredPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(enteredPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Wrong password
			log.Println("Wrong password entered: ", err)
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
