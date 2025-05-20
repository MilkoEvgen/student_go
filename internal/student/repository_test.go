package student

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	dbcontext.DB = gormDB
	return db, mock, gormDB
}

func TestStudentExistsById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) > 0 FROM "students" WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(true))

	repo := NewStudentRepository()
	exists, err := repo.ExistsById(1)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestStudentSave(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "students" ("name","email") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs("John", "john@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	repo := NewStudentRepository()
	student := &entity.Student{Name: "John", Email: "john@example.com"}
	result, err := repo.Save(student)

	assert.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
}

func TestStudentFindById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM "students" WHERE "students"\."id" = \$1 ORDER BY "students"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@example.com"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."student_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"student_id", "course_id"}).
			AddRow(1, 101).
			AddRow(1, 102))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."id" IN ($1,$2)`)).
		WithArgs(101, 102).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(101, "Math", 201).
			AddRow(102, "Physics", 202))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1,$2)`)).
		WithArgs(201, 202).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(201, "Dr. Smith").
			AddRow(202, "Prof. Jane"))

	repo := NewStudentRepository()
	student, err := repo.FindById(1)

	require.NoError(t, err)
	require.NotNil(t, student)
	assert.Equal(t, "Alice", student.Name)
	assert.Equal(t, "alice@example.com", student.Email)
	require.Len(t, student.Courses, 2)
	assert.Equal(t, "Math", student.Courses[0].Title)
	assert.Equal(t, "Physics", student.Courses[1].Title)
}

func TestStudentUpdate(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students"`)).
		WithArgs("UpdatedName", "", 1). // Name, Email, ID
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT \* FROM "students" WHERE "students"\."id" = \$1 ORDER BY "students"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "UpdatedName", ""))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."student_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"student_id", "course_id"}).
			AddRow(1, 101).
			AddRow(1, 102))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."id" IN ($1,$2)`)).
		WithArgs(101, 102).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(101, "Math", 201).
			AddRow(102, "Physics", 202))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1,$2)`)).
		WithArgs(201, 202).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(201, "Dr. Smith").
			AddRow(202, "Prof. Jane"))

	repo := NewStudentRepository()
	st := &entity.Student{ID: 1, Name: "UpdatedName"}
	updated, err := repo.Update(st)

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "UpdatedName", updated.Name)
	require.Len(t, updated.Courses, 2)
	assert.Equal(t, "Math", updated.Courses[0].Title)
	assert.Equal(t, "Physics", updated.Courses[1].Title)
	assert.Equal(t, "Dr. Smith", updated.Courses[0].Teacher.Name)
	assert.Equal(t, "Prof. Jane", updated.Courses[1].Teacher.Name)
}

func TestStudentFindAll(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	page := 1
	limit := 2

	mock.ExpectQuery(`SELECT \* FROM "students" LIMIT \$\d+`).
		WithArgs(limit).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@example.com").
			AddRow(2, "Bob", "bob@example.com"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."student_id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"student_id", "course_id"}).
			AddRow(1, 101).
			AddRow(2, 102))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."id" IN ($1,$2)`)).
		WithArgs(101, 102).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(101, "Math", 201).
			AddRow(102, "Physics", 202))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1,$2)`)).
		WithArgs(201, 202).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(201, "Dr. Smith").
			AddRow(202, "Prof. Jane"))

	repo := NewStudentRepository()
	students, err := repo.FindAll(page, limit)

	require.NoError(t, err)
	require.Len(t, students, 2)

	assert.Equal(t, "Alice", students[0].Name)
	assert.Equal(t, "Bob", students[1].Name)

	require.Len(t, students[0].Courses, 1)
	require.Len(t, students[1].Courses, 1)

	assert.Equal(t, "Math", students[0].Courses[0].Title)
	assert.Equal(t, "Dr. Smith", students[0].Courses[0].Teacher.Name)

	assert.Equal(t, "Physics", students[1].Courses[0].Title)
	assert.Equal(t, "Prof. Jane", students[1].Courses[0].Teacher.Name)
}

func TestStudentDeleteById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "students" WHERE "students"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewStudentRepository()
	err := repo.DeleteById(1)

	assert.NoError(t, err)
}

func TestStudentCount(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "students"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	repo := NewStudentRepository()
	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
