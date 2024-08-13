package app

import (
	"manajemen_tugas_master/controller"
	"manajemen_tugas_master/repository"
	"manajemen_tugas_master/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

// user
func InitializeRepositoryUser(db *gorm.DB) (repository.UserRepository, error) {
	return repository.NewUserRepository(db), nil
}

func InitializeServiceUser(userRepository repository.UserRepository) (service.UserService, error) {
	return service.NewUserService(userRepository, validator.New()), nil
}

func InitializeControllerUser(userService service.UserService, store *session.Store) (controller.UserController, error) {
	return *controller.NewUserController(userService, store), nil
}

// board
func InitializeRepositoryBoard(db *gorm.DB) (repository.BoardRepository, error) {
	return repository.NewBoardRepository(db), nil
}

func InitializeServiceBoard(boardRepository repository.BoardRepository) (service.BoardService, error) {
	return service.NewBoardService(boardRepository), nil
}
func InitializeControllerBoard(boardService service.BoardService) (controller.BoardController, error) {
	return *controller.NewBoardController(boardService), nil
}

// task
func InitializeRepositoryTask(db *gorm.DB) (repository.TaskAndOwnerRepository, error) {
	return repository.NewTaskAndOwnerRepository(db), nil
}

func InitializeServiceTask(taskAndOwnerRepository repository.TaskAndOwnerRepository, boardRepository repository.BoardRepository) (service.TaskAndOwnerService, error) {
	return service.NewTaskAndOwnerService(taskAndOwnerRepository, boardRepository, validator.New()), nil
}

func InitializeControllerTask(taskAndOwnerService service.TaskAndOwnerService) (controller.TaskAndOwnerController, error) {
	return *controller.NewTaskController(taskAndOwnerService), nil
}
