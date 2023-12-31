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

type CourseService struct {
	LoggerSugar                    *zap.SugaredLogger
	CourseDomainDataBaseRepository output.ICourseDomainDataBaseRepository
	CourseDomainCacheRepository    output.ICourseDomainCacheRepository
}

var CourseCacheTTL = 10 * time.Minute

const (
	CourseCacheKeyTypeID = "ID"
)

const (
	CourseErrorToSaveInCache                     = "error to save course in cache"
	CourseErrorToGetByIDInCache                  = "error to get course in cache"
	CourseDomainCacheRepositoryNotInitialized    = "CourseDomainCacheRepository not initialized"
	CourseDomainDataBaseRepositoryNotInitialized = "CourseDomainDataBaseRepository not initialized"
)

func (service *CourseService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *CourseService) Create(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error) {

	course.RegistrationDate = time.Now()
	save, err := service.CourseDomainDataBaseRepository.Save(contextControl, course)
	if err != nil {
		return domain.CourseDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.CourseDomainCacheRepository.Set(contextControl,
		service.getCacheKey(CourseCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), CourseCacheTTL); err != nil {
		service.LoggerSugar.Infow(CourseErrorToSaveInCache, "course_id", save.ID)
	}

	return save, nil
}

func (service *CourseService) GetByID(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, bool, error) {
	course, exists, err := service.CourseDomainDataBaseRepository.GetByID(contextControl, ID)
	if err != nil {
		return domain.CourseDomain{}, exists, err
	}

	if !exists {
		return domain.CourseDomain{}, exists, nil
	}
	hash, _ := json.Marshal(course)
	if err = service.CourseDomainCacheRepository.Set(contextControl,
		service.getCacheKey(CourseCacheKeyTypeID, strconv.FormatInt(course.ID, 10)),
		string(hash), CourseCacheTTL); err != nil {
		service.LoggerSugar.Infow(CourseErrorToGetByIDInCache, "id", course.ID) //TODO: to adjust the keyAndValues
	}
	return course, exists, nil
}

func (service *CourseService) GetCourses(contextControl domain.ContextControl, courses []domain.CourseDomain) ([]domain.CourseDomain, error) {

	if service.CourseDomainDataBaseRepository == nil {
		return []domain.CourseDomain{}, fmt.Errorf(CourseDomainDataBaseRepositoryNotInitialized)
	}

	coursesFromDB, err := service.CourseDomainDataBaseRepository.GetCourses(contextControl, courses)
	if err != nil {
		return []domain.CourseDomain{}, err
	}

	if service.CourseDomainCacheRepository == nil {
		return []domain.CourseDomain{}, fmt.Errorf(CourseDomainCacheRepositoryNotInitialized)
	}

	for _, course := range coursesFromDB {
		hash, err := json.Marshal(course)
		if err != nil {
			return []domain.CourseDomain{}, err
		}

		cacheKey := service.getCacheKey(CourseCacheKeyTypeID, strconv.FormatInt(course.ID, 10))

		if err := service.CourseDomainCacheRepository.Set(contextControl, cacheKey, string(hash), CourseCacheTTL); err != nil {
			service.LoggerSugar.Warnw("CourseErrorToSetInCache", "id", course.ID)
		}
	}

	return coursesFromDB, nil
}
