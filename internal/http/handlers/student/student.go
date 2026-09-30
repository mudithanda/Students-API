package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	//"github.com/mudithanda/students-api/internal/http/handlers/student"
	"github.com/mudithanda/students-api/internal/storage"
	"github.com/mudithanda/students-api/internal/types"
	"github.com/mudithanda/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"))) // custom err
			return 
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {

			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return 
		}

		// creating student
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully",slog.String("userId",fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}


		//w.Write([]byte("Welcome to students api"))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("updating student")

		idStr := r.PathValue("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		} 
		
		var student types.Student

		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		// Update student in db
		rowsAffected, err := storage.UpdateStudent(
			id,
			student.Name,
			student.Email,
			student.Age,

		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
		
		response.WriteJson(w, http.StatusOK, rowsAffected)
	}
}

func DeleteList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Deleting student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		rowsAffected, err := storage.DeleteStudent(intId)
		if err != nil {
			slog.Error("error deleting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, rowsAffected)
	}
}