package course

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

func TestCourseExistsById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) > 0 FROM "courses" WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(true))

	repo := NewCourseRepository()
	exists, err := repo.ExistsById(1)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCourseSave(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "courses" ("title","teacher_id") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs("Math", nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	repo := NewCourseRepository()
	course := &entity.Course{Title: "Math"}
	result, err := repo.Save(course)

	assert.NoError(t, err)
	assert.Equal(t, "Math", result.Title)
}

func TestCourseFindById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE "courses"\."id" = \$1 ORDER BY "courses"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(1, "Math", 101))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."course_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id", "student_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" WHERE "students"."id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@example.com").
			AddRow(2, "Bob", "bob@example.com"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" = $1`)).
		WithArgs(101).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(101, "Dr. Smith"))

	repo := NewCourseRepository()
	course, err := repo.FindById(1)

	require.NoError(t, err)
	require.NotNil(t, course)
	assert.Equal(t, "Math", course.Title)
	require.Len(t, course.Students, 2)
	assert.Equal(t, "Alice", course.Students[0].Name)
	assert.Equal(t, "Bob", course.Students[1].Name)
	assert.Equal(t, "Dr. Smith", course.Teacher.Name)
}

func TestCourseUpdate(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "courses" SET "title"=$1 WHERE id = $2`)).
		WithArgs("Updated Title", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE "courses"\."id" = \$1 ORDER BY "courses"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(1, "Updated Title", 101))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."course_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id", "student_id"}).
			AddRow(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" WHERE "students"."id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@example.com"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" = $1`)).
		WithArgs(101).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(101, "Dr. Smith"))

	repo := NewCourseRepository()
	c := &entity.Course{ID: 1, Title: "Updated Title"}
	updated, err := repo.Update(c)

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Dr. Smith", updated.Teacher.Name)
	assert.Len(t, updated.Students, 1)
	assert.Equal(t, "Alice", updated.Students[0].Name)
}

func TestCourseFindAll(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	page := 1
	limit := 2

	mock.ExpectQuery(`SELECT \* FROM "courses"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "teacher_id"}).
			AddRow(1, "Math", 101).
			AddRow(2, "Physics", 102))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course_student" WHERE "course_student"."course_id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"course_id", "student_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" WHERE "students"."id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@example.com").
			AddRow(2, "Bob", "bob@example.com"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1,$2)`)).
		WithArgs(101, 102).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(101, "Dr. Smith").
			AddRow(102, "Prof. Jane"))

	repo := NewCourseRepository()
	courses, err := repo.FindAll(page, limit)

	require.NoError(t, err)
	require.Len(t, courses, 2)
	assert.Equal(t, "Math", courses[0].Title)
	assert.Equal(t, "Physics", courses[1].Title)
	assert.Equal(t, "Dr. Smith", courses[0].Teacher.Name)
	assert.Equal(t, "Prof. Jane", courses[1].Teacher.Name)
}

func TestCourseDeleteById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "courses" WHERE "courses"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCourseRepository()
	err := repo.DeleteById(1)

	assert.NoError(t, err)
}

func TestCourseCount(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "courses"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	repo := NewCourseRepository()
	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
