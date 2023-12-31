package database

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/CourseHub-API/application/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CoursePostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

func NewCoursePostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) CoursePostgresDB {
	return CoursePostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (cp CoursePostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, error) {
	//TODO implement me
	panic("implement me")
}

type CourseDB struct {
	ID               int64     `gorm:"primaryKey, column:id"`
	Description      string    `gorm:"column:description"`
	Outline          string    `gorm:"column:outline"`
	RegistrationDate time.Time `gorm:"column:registration_date"`
}

func (CourseDB) TableName() string {
	return "coursehub_api.course"
}

func (c CourseDB) CopyToCourseDomain() domain.CourseDomain {
	return domain.CourseDomain{
		ID:               c.ID,
		Description:      c.Description,
		Outline:          c.Outline,
		RegistrationDate: c.RegistrationDate,
	}
}

func (cp CoursePostgresDB) Save(contextControl domain.ContextControl, courseDomain domain.CourseDomain) (domain.CourseDomain, error) {

	var courseDB CourseDB
	copier.Copy(&courseDB, &courseDomain)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&courseDB).Error; err != nil {
		cp.LoggerSugar.Errorw("error to save into postgres",
			"error", err.Error())
		return domain.CourseDomain{}, err
	}

	return courseDB.CopyToCourseDomain(), nil
}
