package storage

import "github.com/mudithanda/students-api/internal/types"

type Storage interface {
	// creating method
	CreateStudent(name string, email string, age int) (int64, error)

	// getting student
	GetStudentById(id int64) (types.Student, error)

	GetStudents() ([]types.Student, error)

	UpdateStudent(id int64, name string, email string, age int) (int64, error)
	DeleteStudent(id int64) (int64, error)
}