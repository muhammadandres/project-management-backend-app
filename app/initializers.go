package app

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"manajemen_tugas_master/controller"
	"manajemen_tugas_master/repository"
	"manajemen_tugas_master/service"
)

// user
func InitializeRepositoryUser(db *gorm.DB) (repository.UserRepository, error) {
	return repository.NewUserRepository(db), nil
}

func InitializeServiceUser(userRepository repository.UserRepository) (service.UserService, error) {
	return service.NewUserService(userRepository, validator.New()), nil
}

func InitializeControllerUser(userService service.UserService) (controller.UserController, error) {
	return *controller.NewUserController(userService), nil
}

// task
func InitializeRepositoryTask(db *gorm.DB) (repository.TaskAndOwnerRepository, error) {
	return repository.NewTaskAndOwnerRepository(db), nil
}

func InitializeServiceTask(taskAndOwnerRepository repository.TaskAndOwnerRepository) (service.TaskAndOwnerService, error) {
	return service.NewTaskAndOwnerService(taskAndOwnerRepository, validator.New()), nil
}
func InitializeControllerTask(taskAndOwnerService service.TaskAndOwnerService) (controller.TaskAndOwnerController, error) {
	return *controller.NewTaskController(taskAndOwnerService), nil
}
