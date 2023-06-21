package adminapiservice

import (
	taskmanager "admin-panel/gen/task_manager"
	"context"
	"log"
)

// taskManager service example implementation.
// The example methods log the requests and return zero values.
type taskManagersrvc struct {
	logger *log.Logger
}

// NewTaskManager returns the taskManager service implementation.
func NewTaskManager(logger *log.Logger) taskmanager.Service {
	return &taskManagersrvc{logger}
}

// TaskList implements taskList.
func (s *taskManagersrvc) TaskList(ctx context.Context) (res *taskmanager.TaskListResult, err error) {
	res = &taskmanager.TaskListResult{}
	s.logger.Print("taskManager.taskList")
	return
}

// TaskDeploy implements taskDeploy.
func (s *taskManagersrvc) TaskDeploy(ctx context.Context, p *taskmanager.TaskDeploy2) (res *taskmanager.TaskDeployResult, err error) {
	res = &taskmanager.TaskDeployResult{}
	s.logger.Print("taskManager.taskDeploy")
	return
}

// UnDeploy implements unDeploy.
func (s *taskManagersrvc) UnDeploy(ctx context.Context, p *taskmanager.TaskDeploy2) (res *taskmanager.UnDeployResult, err error) {
	res = &taskmanager.UnDeployResult{}
	s.logger.Print("taskManager.unDeploy")
	return
}

// TaskCreate implements taskCreate.
func (s *taskManagersrvc) TaskCreate(ctx context.Context, p *taskmanager.TaskItem) (res *taskmanager.TaskCreateResult, err error) {
	res = &taskmanager.TaskCreateResult{}
	s.logger.Print("taskManager.taskCreate")
	return
}
