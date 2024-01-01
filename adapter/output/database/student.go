package database

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/CourseHub-API/application/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const (
	StudentNotFound = "student not found"
)

type StudentPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

func NewStudentPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) StudentPostgresDB {
	return StudentPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

type StudentDB struct {
	ID               int64     `gorm:"primaryKey, column:id"`
	Name             string    `gorm:"column:name"`
	RegistrationDate time.Time `gorm:"column:registration_date"`
}

func (*StudentDB) TableName() string {
	return "coursehub_api.student"
}

func (c *StudentDB) CopyToStudentDomain() domain.StudentDomain {
	return domain.StudentDomain{
		ID:               c.ID,
		Name:             c.Name,
		RegistrationDate: c.RegistrationDate,
	}
}

func (cp *StudentPostgresDB) Save(contextControl domain.ContextControl, studentDomain domain.StudentDomain) (domain.StudentDomain, error) {

	var studentDB StudentDB
	copier.Copy(&studentDB, &studentDomain)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&studentDB).Error; err != nil {
		cp.LoggerSugar.Errorw("error to save into postgres",
			"error", err.Error())
		return domain.StudentDomain{}, err
	}

	return studentDB.CopyToStudentDomain(), nil
}

func (cp *StudentPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error) {
	var studentDB StudentDB

	result := cp.DB.WithContext(contextControl.Context).First(&studentDB, ID)
	if result.RowsAffected == 0 {
		cp.LoggerSugar.Errorw(StudentNotFound)
		return domain.StudentDomain{}, false, nil
	}
	return studentDB.CopyToStudentDomain(), true, nil
}

func (cp *StudentPostgresDB) GetStudents(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error) {
	var studentsDB []StudentDB
	err := cp.DB.WithContext(contextControl.Context).Find(&studentsDB).Error
	if err != nil {
		cp.LoggerSugar.Errorw("error to get list of students", "error", err.Error())
		return nil, err
	}

	var domainStudents []domain.StudentDomain
	for _, studentDB := range studentsDB {
		domainStudents = append(domainStudents, studentDB.CopyToStudentDomain())
	}

	return domainStudents, nil
}
