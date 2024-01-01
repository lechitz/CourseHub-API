package service

import (
	"encoding/json"
	"fmt"
	"github.com/lechitz/CourseHub-API/application/domain"
	"github.com/lechitz/CourseHub-API/application/port/output"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type StudentService struct {
	LoggerSugar                     *zap.SugaredLogger
	StudentDomainDataBaseRepository output.IStudentDomainDataBaseRepository
	StudentDomainCacheRepository    output.IStudentDomainCacheRepository
}

var StudentCacheTTL = 10 * time.Minute

const (
	StudentCacheKeyTypeID = "ID"
)

const (
	StudentErrorToSaveInCache                     = "error to save student in cache"
	StudentErrorToGetByIDInCache                  = "error to get student in cache"
	StudentDomainCacheRepositoryNotInitialized    = "StudentDomainCacheRepository not initialized"
	StudentDomainDataBaseRepositoryNotInitialized = "StudentDomainDataBaseRepository not initialized"
)

func (service *StudentService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *StudentService) Create(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error) {

	student.RegistrationDate = time.Now()
	save, err := service.StudentDomainDataBaseRepository.Save(contextControl, student)
	if err != nil {
		return domain.StudentDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.StudentDomainCacheRepository.Set(contextControl,
		service.getCacheKey(StudentCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), StudentCacheTTL); err != nil {
		service.LoggerSugar.Infow(StudentErrorToSaveInCache, "student_id", save.ID)
	}

	return save, nil
}

func (service *StudentService) GetByID(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error) {
	student, exists, err := service.StudentDomainDataBaseRepository.GetByID(contextControl, ID)
	if err != nil {
		return domain.StudentDomain{}, exists, err
	}

	if !exists {
		return domain.StudentDomain{}, exists, nil
	}
	hash, _ := json.Marshal(student)
	if err = service.StudentDomainCacheRepository.Set(contextControl,
		service.getCacheKey(StudentCacheKeyTypeID, strconv.FormatInt(student.ID, 10)),
		string(hash), StudentCacheTTL); err != nil {
		service.LoggerSugar.Infow(StudentErrorToGetByIDInCache, "student_id", student.ID)
	}
	return student, exists, nil
}

func (service *StudentService) GetStudents(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error) {

	if service.StudentDomainDataBaseRepository == nil {
		return []domain.StudentDomain{}, fmt.Errorf(StudentDomainDataBaseRepositoryNotInitialized)
	}

	studentsFromDB, err := service.StudentDomainDataBaseRepository.GetStudents(contextControl, students)
	if err != nil {
		return []domain.StudentDomain{}, err
	}

	if service.StudentDomainCacheRepository == nil {
		return []domain.StudentDomain{}, fmt.Errorf(StudentDomainCacheRepositoryNotInitialized)
	}

	for _, student := range studentsFromDB {
		hash, err := json.Marshal(student)
		if err != nil {
			return []domain.StudentDomain{}, err
		}

		cacheKey := service.getCacheKey(StudentCacheKeyTypeID, strconv.FormatInt(student.ID, 10))

		if err := service.StudentDomainCacheRepository.Set(contextControl, cacheKey, string(hash), StudentCacheTTL); err != nil {
			service.LoggerSugar.Warnw("StudentErrorToSetInCache", "student_id", student.ID)
		}
	}

	return studentsFromDB, nil
}
