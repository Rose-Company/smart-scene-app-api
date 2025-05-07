package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"smart-scene-app-api/common"
	"smart-scene-app-api/pkg"
	"smart-scene-app-api/pkg/rest_service"
	logger2 "smart-scene-app-api/services/logger"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
)

type service interface {
	Run() error
	Stop() <-chan bool
	GetPrefix() string
	Get() interface{}
}

type server struct {
	prefix          string
	port            uint
	services        map[string]service
	restHandler     func() *gin.Engine
	logger          logger2.Loggers
	ctx             context.Context
	user            interface{}
	loggers         map[string]logger2.Loggers
	authorization   AuthorizationConfig
	jobs            map[string]*pkg.Job
	telegramService rest_service.RestInterface
	sesClient       *pkg.AWSSesClient
}

type JobHandler func() error

func NewServer(prefix string, port uint) *server {
	svs := make(map[string]service)
	logs := make(map[string]logger2.Loggers)
	jobs := map[string]*pkg.Job{}
	return &server{prefix: prefix, port: port, services: svs, loggers: logs, jobs: jobs}
}

func (s *server) Run() error {
	if err := s.configure(); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan error, 1)

	waitGroup := sync.WaitGroup{}
	for _, svc := range s.services {
		waitGroup.Add(1)
		go func(sv service, wg *sync.WaitGroup) {
			s.logger.Info().Println(fmt.Sprintf("%v is open", sv.GetPrefix()))
			defer wg.Done()
			if err := sv.Run(); err != nil {
				s.logger.Error().Println("Err: ", err)
				stop <- err
			} else {
				s.logger.Info().Println(fmt.Sprintf("%v is running", sv.GetPrefix()))
			}

		}(svc, &waitGroup)
	}
	waitGroup.Wait()

	if s.restHandler != nil {
		go func() {
			s.logger.Info().Println(fmt.Sprintf("%v is running in port %v", s.prefix, s.port))

			if err := s.restHandler().Run(fmt.Sprintf(":%v", s.port)); err != nil {
				stop <- err
			}
		}()

	}
	go func() {
		for _, job := range s.jobs {
			go func(job *pkg.Job) {
				s.logger.Info().Println(fmt.Sprintf("CronJob %v is running"))
				err := job.Run()
				if err != nil {
					return
				}
			}(job)
		}
	}()

	for {
		select {
		case err := <-stop:
			if err != nil {
				return err
			}

		case sig := <-sigs:
			if sig != nil {
				for _, svc := range s.services {
					<-svc.Stop()
					s.logger.Info().Println(fmt.Sprintf("%v is stopped", svc.GetPrefix()))
				}

				return errors.New(sig.String())
			}
		}
	}
}

func (s *server) configure() error {
	if err := s.initFlags(); err != nil {
		return err
	}

	if s.port == 0 {
		return errors.New(common.DataIsNullErr("Port"))
	}

	return nil
}

func (s *server) initFlags() error {
	return nil
}

func (s *server) InitService(svc service) {
	if has, ok := s.services[svc.GetPrefix()]; ok {
		s.logger.Error().Fatal(fmt.Sprintf("Service %v is duplicated", has.GetPrefix()))
	}

	s.services[svc.GetPrefix()] = svc
}

func (s *server) InitContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *server) GetContext() context.Context {
	return s.ctx
}

func (s *server) AddHandler(hdl func() *gin.Engine) {
	s.restHandler = hdl
}

func (s *server) AddLogger(loggers logger2.Loggers) {
	s.logger = loggers
}

func (s *server) GetService(prefix string) interface{} {
	if svc, ok := s.services[prefix]; ok {
		return svc.Get()
	}

	return nil
}

func (s *server) GetLogger() logger2.Loggers {
	return s.logger
}

func (s *server) GetUser() interface{} {
	return s.user
}

func (s *server) SetUser(value interface{}) {
	s.user = value
}

func (s *server) InitLogger(prefix string, logger logger2.Loggers) {
	if _, ok := s.loggers[prefix]; ok {
		s.logger.Error().Fatal(fmt.Sprintf("Service %v is duplicated", prefix))
	}

	s.loggers[prefix] = logger
}

func (s *server) GetLoggerWithPrefix(prefix string) logger2.Loggers {
	return s.loggers[prefix]
}

type RedisService interface {
	GetRedsync() redsync.Redsync
}

func (s *server) GetRedisRedsync(prefix string) redsync.Redsync {
	rd := s.services[prefix].(RedisService)
	return rd.GetRedsync()
}

func (s *server) AddJob(jobName string, job *pkg.Job) {
	s.jobs[jobName] = job
}

func (s *server) GetAuthConfig() *AuthorizationConfig {
	return &s.authorization
}

func (s *server) SetTelegramService(service rest_service.RestInterface) {
	s.telegramService = service
}

func (s *server) GetTelegramService() rest_service.RestInterface {
	return s.telegramService
}

func (s *server) SetAwsSes(service *pkg.AWSSesClient) {
	s.sesClient = service
}

func (s *server) GetAwsSes() *pkg.AWSSesClient {
	return s.sesClient
}
