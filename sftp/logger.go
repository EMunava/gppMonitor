package sftp

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) GetFilesInPath(path string) (file []File, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetFilesInPath",
			"path", path,
			"files", file,
			"eror", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.GetFilesInPath(path)
}

func (s *loggingService) RetrieveFile(path, file string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RetrieveFile",
			"path", path,
			"files", file,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.RetrieveFile(path, file)
}
