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

var CacheTTL = 10 * time.Minute

const (
	CourseCacheKeyTypeID = "ID"
)

const (
	CourseErrorToSaveInCache = "error to save course in cache."
)

func (service CourseService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service CourseService) Create(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error) {

	course.RegistrationDate = time.Now()
	save, err := service.CourseDomainDataBaseRepository.Save(contextControl, course)
	if err != nil {
		return domain.CourseDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.CourseDomainCacheRepository.Set(contextControl,
		service.getCacheKey(CourseCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), CacheTTL); err != nil {
		service.LoggerSugar.Infow(CourseErrorToSaveInCache, "course_id", save.ID)
	}

	return save, nil
}
