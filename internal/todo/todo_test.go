package todo_test

import (
	"context"
	"my-first-api/internal/db"
	"my-first-api/internal/todo"
	"reflect"
	"testing"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil
}

func TestService_Search(t *testing.T) {
	type fields struct {
		todos []todo.Item
	}
	type args struct {
		query string
	}
	tests := []struct {
		name           string
		toDosToAdd     []string
		query          string
		expectedResult any
	}{
		{
			name:           "given a todo of shop and a search of sh, i should get shop back",
			toDosToAdd:     []string{"shop"},
			query:          "sh",
			expectedResult: []string{"shop"},
		},
		{
			name:           "still returns shop, even if the case doesn't match",
			toDosToAdd:     []string{"Shopping"},
			query:          "sh",
			expectedResult: []string{"Shopping"},
		},
		{
			name:           "spaces",
			toDosToAdd:     []string{"go Shopping"},
			query:          "go",
			expectedResult: []string{"go Shopping"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := svc.Search(tt.query)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func TestService_Add(t *testing.T) {
	type fields struct {
		todos []todo.Item
	}
	tests := []struct {
		name       string
		toDosToAdd []string
		wantErr    bool
	}{
		{
			name:       "Add throws error is duplicate todos",
			toDosToAdd: []string{"Shopping", "Shopping"},
			wantErr:    true,
		},
		{
			name:       "Doesnt throw error if no duplicates",
			toDosToAdd: []string{"Shopping", "Cleaning", "Laundry"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					if !tt.wantErr {
						t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
					}
				}

			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type fields struct {
		todos []todo.Item
	}
	tests := []struct {
		name       string
		toDosToAdd []string
		want       []todo.Item
	}{
		{name: "GetAll returns all added todos",
			toDosToAdd: []string{"Shopping", "Cleaning", "Laundry"},
			want: []todo.Item{
				{
					Task:   "Shopping",
					Status: "TO_BE_STARTED",
				},
				{
					Task:   "Cleaning",
					Status: "TO_BE_STARTED",
				},
				{
					Task:   "Laundry",
					Status: "TO_BE_STARTED",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := svc.GetAll()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
