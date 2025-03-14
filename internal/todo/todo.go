package todo

import (
	"context"
	"errors"
	"fmt"
	"my-first-api/internal/db"
	"strings"
)

type Item struct {
	Task   string
	Status string
}

type Manager interface {
	InsertItem(ctx context.Context, item db.Item) error
	GetAllItems(ctx context.Context) ([]db.Item, error)
}

type Service struct {
	db Manager
	//var todos = make(map[int]string) map was used by me for delete
}

func NewService(db Manager) *Service {
	return &Service{
		db: db,
	}
}

func (svc *Service) Add(todo string) error {
	//present := slices.Contains(svc.todos, todo)
	//if !present {
	//	svc.todos = append(svc.todos, todo)
	//}

	items, err := svc.GetAll()
	if err != nil {
		return fmt.Errorf("could not get items: %w", err)
	}

	for _, t := range items {
		if t.Task == todo {
			return errors.New("todo is not unique")
		}
	}
	err = svc.db.InsertItem(context.Background(), db.Item{
		Task:   todo,
		Status: "TO_BE_STARTED",
	})
	if err != nil {
		return fmt.Errorf("could not insert item: %w", err)
	}

	return nil
}

func (svc *Service) Search(query string) ([]string, error) {
	items, err := svc.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get items: %w", err)
	}

	var result []string
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), strings.ToLower(query)) {
			result = append(result, todo.Task)
		}
	}
	return result, nil
}

func (svc *Service) GetAll() ([]Item, error) {
	var results []Item
	items, err := svc.db.GetAllItems(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to read from db: %w", err)
	}
	for _, item := range items {
		results = append(results, Item{
			Task:   item.Task,
			Status: item.Status,
		})
	}
	return results, nil
}
