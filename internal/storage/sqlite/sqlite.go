package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // _ because being used indirectly
	"github.com/mudithanda/students-api/internal/config"
	"github.com/mudithanda/students-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

// creating instance of sqlite struct
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CRETE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	name TEXT
	email TEXT
	age INTEGER
	)`)
	
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil // err
}

func(s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)") // protect fron sql injection

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

	// return 0, nil
}

func(s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}
	return  student, nil
}

func(s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}
	return students, nil
}

func (s* Sqlite) UpdateStudent(id int64, name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (s *Sqlite) DeleteStudent(id int64) (int64, error){
	stmt, err := s.Db.Prepare(`DELETE FROM students WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}