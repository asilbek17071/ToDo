package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asilbek17071/ToDo/genproto"
	l "github.com/asilbek17071/ToDo/pkg/logger"
	"github.com/asilbek17071/ToDo/storage"
)

type TaskService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewTaskService ...
func NewTaskService(storage storage.IStorage, log l.Logger) *TaskService {
	return &TaskService{
		storage: storage,
		logger:  log,
	}
}

func (s *TaskService) Create(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	task, err := s.storage.Task().Create(*req)
	if err != nil {
		s.logger.Error("failed to create task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	return &task, nil
}

func (s *TaskService) Get(ctx context.Context, req *pb.ByIdReq) (*pb.Task, error) {
	task, err := s.storage.Task().Get(req.GetId())
	if err != nil {
		s.logger.Error("failed to get task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	return &task, nil
}

func (s *TaskService) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	tasks, count, err := s.storage.Task().List(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list tasks")
	}

	return &pb.ListResp{
		Tasks: tasks,
		Count: count,
	}, nil
}

func (s *TaskService) Update(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	task, err := s.storage.Task().Update(*req)
	if err != nil {
		s.logger.Error("failed to update task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	return &task, nil
}

func (s *TaskService) Delete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyRes, error) {
	err := s.storage.Task().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	return &pb.EmptyRes{}, nil
}

func (s *TaskService) Overdue(ctx context.Context, req *pb.OverdueReq) (*pb.OverdueResp, error) {
	tasks, count, err := s.storage.Task().Overdue(req.Timed, req.Limit, req.Page)
	if err != nil {
		s.logger.Error("failed to list tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list tasks")
	}

	return &pb.OverdueResp{
		Overres: tasks,
		Count:   count,
	}, nil
}
