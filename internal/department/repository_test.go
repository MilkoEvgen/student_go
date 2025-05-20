package department

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

func TestDepartmentExistsById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) > 0 FROM "departments" WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(true))

	repo := NewDepartmentRepository()
	exists, err := repo.ExistsById(1)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDepartmentSave(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments" ("name","head_of_department_id") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs("Science", nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	repo := NewDepartmentRepository()
	dept := &entity.Department{Name: "Science"}
	result, err := repo.Save(dept)

	assert.NoError(t, err)
	assert.Equal(t, "Science", result.Name)
}

func TestDepartmentFindById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM "departments" WHERE "departments"\."id" = \$1 ORDER BY "departments"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "Science", 10))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" = $1`)).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(10, "Dr. Smith"))

	repo := NewDepartmentRepository()
	dept, err := repo.FindById(1)

	require.NoError(t, err)
	require.NotNil(t, dept)
	assert.Equal(t, "Science", dept.Name)
	assert.Equal(t, "Dr. Smith", dept.HeadOfDepartment.Name)
}

func TestDepartmentUpdate(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments" SET "name"=$1 WHERE id = $2`)).
		WithArgs("Updated", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT \* FROM "departments" WHERE "departments"\."id" = \$1 ORDER BY "departments"\."id" LIMIT .*`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "Updated", 11))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" = $1`)).
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(11, "Prof. Jane"))

	repo := NewDepartmentRepository()
	updated, err := repo.Update(&entity.Department{ID: 1, Name: "Updated"})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, "Prof. Jane", updated.HeadOfDepartment.Name)
}

func TestDepartmentFindAll(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	page := 1
	limit := 2

	mock.ExpectQuery(`SELECT \* FROM "departments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "head_of_department_id"}).
			AddRow(1, "Science", 10).
			AddRow(2, "Arts", 11))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1,$2)`)).
		WithArgs(10, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(10, "Dr. Smith").
			AddRow(11, "Prof. Jane"))

	repo := NewDepartmentRepository()
	depts, err := repo.FindAll(page, limit)

	require.NoError(t, err)
	require.Len(t, depts, 2)
	assert.Equal(t, "Science", depts[0].Name)
	assert.Equal(t, "Arts", depts[1].Name)
	assert.Equal(t, "Dr. Smith", depts[0].HeadOfDepartment.Name)
	assert.Equal(t, "Prof. Jane", depts[1].HeadOfDepartment.Name)
}

func TestDepartmentDeleteById(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "departments" WHERE "departments"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewDepartmentRepository()
	err := repo.DeleteById(1)

	assert.NoError(t, err)
}

func TestDepartmentCount(t *testing.T) {
	db, mock, _ := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "departments"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	repo := NewDepartmentRepository()
	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
