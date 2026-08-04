package lifecycle

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestManagerStartAndStopOrder(t *testing.T) {
	manager := NewManager()
	var calls []string

	manager.Register(Hook{
		Name: "database",
		Start: func(context.Context) error {
			calls = append(calls, "start database")
			return nil
		},
		Stop: func(context.Context) error {
			calls = append(calls, "stop database")
			return nil
		},
	})

	manager.Register(Hook{
		Name: "http",
		Start: func(context.Context) error {
			calls = append(calls, "start http")
			return nil
		},
		Stop: func(context.Context) error {
			calls = append(calls, "stop http")
			return nil
		},
	})

	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	want := []string{
		"start database",
		"start http",
		"stop http",
		"stop database",
	}

	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestManagerStartReturnsHookNameOnError(t *testing.T) {
	manager := NewManager()
	startErr := errors.New("connect failed")

	manager.Register(Hook{
		Name: "database",
		Start: func(context.Context) error {
			return startErr
		},
	})

	err := manager.Start(context.Background())
	if !errors.Is(err, startErr) {
		t.Fatalf("err = %v, want wrapped %v", err, startErr)
	}
}
