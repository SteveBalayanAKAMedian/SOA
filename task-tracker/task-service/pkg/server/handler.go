package server

import (
    "context"
    "database/sql"
    "errors"
    "time"
	"log"

    "google.golang.org/protobuf/types/known/timestamppb"
    "task-service/proto"
	"github.com/google/uuid"
)

type Server struct {
	proto.UnimplementedTaskServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.GetTaskResponse, error) {
	task := &proto.Task{}
	var createdAt time.Time
	var updatedAt time.Time

	err := s.db.QueryRow(
		"SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1",
		req.Id,
	).Scan(&task.Id, &task.UserId, &task.Title, &task.Description, &task.Status, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Task not found with id: %s", req.Id) // Debug information
			return nil, errors.New("task not found")
		}
		log.Printf("Error querying task with id: %s, error: %v", req.Id, err) // Debug information
		return nil, err
	}

	task.CreatedAt = timestamppb.New(createdAt)
	task.UpdatedAt = timestamppb.New(updatedAt)

	return &proto.GetTaskResponse{Task: task}, nil
}



func (s *Server) CreateTask(ctx context.Context, req *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
    taskID := uuid.New().String()

    _, err := s.db.ExecContext(ctx, "INSERT INTO tasks (id, user_id, title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        taskID, req.UserId, req.Title, req.Description, "open", time.Now(), time.Now())
    if err != nil {
        return nil, errors.New("failed to create task")
    }

    return &proto.CreateTaskResponse{
        Task: &proto.Task{
            Id:          taskID,
            UserId:      req.UserId,
            Title:       req.Title,
            Description: req.Description,
            Status:      "open",
            CreatedAt:   timestamppb.Now(),
            UpdatedAt:   timestamppb.Now(),
        },
    }, nil
}

func (s *Server) UpdateTask(ctx context.Context, req *proto.UpdateTaskRequest) (*proto.UpdateTaskResponse, error) {
    _, err := s.db.ExecContext(ctx, "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5 AND user_id = $6",
        req.Title, req.Description, req.Status, time.Now(), req.Id, req.UserId)
    if err != nil {
        return nil, errors.New("failed to update task")
    }

    return &proto.UpdateTaskResponse{
        Task: &proto.Task{
            Id:          req.Id,
            UserId:      req.UserId,
            Title:       req.Title,
            Description: req.Description,
            Status:      req.Status,
            UpdatedAt:   timestamppb.Now(),
        },
    }, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
    _, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1 AND user_id = $2", req.Id, req.UserId)
    if err != nil {
        return nil, errors.New("failed to delete task")
    }

    return &proto.DeleteTaskResponse{Success: true}, nil
}

func (s *Server) ListTasks(ctx context.Context, req *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	rows, err := s.db.Query(
		"SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks LIMIT $1 OFFSET $2",
		req.PageSize, (req.Page-1)*req.PageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*proto.Task{}
	for rows.Next() {
		task := &proto.Task{}
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&task.Id, &task.UserId, &task.Title, &task.Description, &task.Status, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		task.CreatedAt = timestamppb.New(createdAt)
		task.UpdatedAt = timestamppb.New(updatedAt)

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &proto.ListTasksResponse{Tasks: tasks}, nil
}

