package teacher

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestExistsById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) > 0 FROM "teachers" WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(true))

	repo := NewTeacherRepository()
	exists, err := repo.ExistsById(1)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestSave(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers" ("name") VALUES ($1) RETURNING "id"`)).
		WithArgs("John").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	repo := NewTeacherRepository()
	tch := &entity.Teacher{Name: "John"}
	result, err := repo.Save(tch)

	assert.NoError(t, err)
	assert.Equal(t, "John", result.Name)
}

func TestFindById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."id" = \$1 ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Alice"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."teacher_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "teacher_id"}).
			AddRow(1, "Math", 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "departments" WHERE "departments"."head_of_department_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "Science", 1))

	repo := NewTeacherRepository()
	tch, err := repo.FindById(1)

	require.NoError(t, err)
	require.NotNil(t, tch)
	assert.Equal(t, "Alice", tch.Name)
}

func TestUpdate(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WithArgs("UpdatedName", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."id" = \$1 ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "UpdatedName"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."teacher_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "teacher_id"}).
			AddRow(1, "Math", 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "departments" WHERE "departments"."head_of_department_id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "Physics", 1))

	repo := NewTeacherRepository()
	tch := &entity.Teacher{ID: 1, Name: "UpdatedName"}
	updated, err := repo.Update(tch)

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "UpdatedName", updated.Name)
}

func TestFindAll(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	page := 2
	limit := 2
	offset := (page - 1) * limit

	mock.ExpectQuery(`SELECT \* FROM "teachers" OFFSET \$\d+`).
		WithArgs(offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Alice").
			AddRow(2, "Bob"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."teacher_id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "teacher_id"}).
			AddRow(1, "Math", 1).
			AddRow(2, "CS", 2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "departments" WHERE "departments"."head_of_department_id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "MathDept", 1).
			AddRow(2, "CSEDept", 2))

	repo := NewTeacherRepository()
	teachers, err := repo.FindAll(page, limit)

	require.NoError(t, err)
	require.Len(t, teachers, 2)
	assert.Equal(t, "Alice", teachers[0].Name)
	assert.Equal(t, "Bob", teachers[1].Name)
}

func TestDeleteById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "teachers" WHERE "teachers"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewTeacherRepository()
	err := repo.DeleteById(1)

	assert.NoError(t, err)
}

func TestCount(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "teachers"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	repo := NewTeacherRepository()
	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
