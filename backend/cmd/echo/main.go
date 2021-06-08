package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bytebase/bytebase/server"
	"github.com/bytebase/bytebase/store"
	"go.uber.org/zap"
)

// const DSN = ":memory:"
const DSN = "./data/bytebase_dev.db"

type Main struct {
	l *zap.Logger

	server *server.Server

	db *store.DB
}

func NewMain() *Main {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to create logger.")
	}
	defer logger.Sync()
	return &Main{
		l: logger,
	}
}

func (m *Main) Run() error {
	db := store.NewDB(m.l, DSN)
	if err := db.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	m.db = db

	server := server.NewServer(m.l)
	server.PrincipalService = store.NewPrincipalService(m.l, db)
	server.MemberService = store.NewMemberService(m.l, db)
	server.ProjectService = store.NewProjectService(m.l, db)
	server.ProjectMemberService = store.NewProjectMemberService(m.l, db)
	server.EnvironmentService = store.NewEnvironmentService(m.l, db)
	server.InstanceService = store.NewInstanceService(m.l, db)
	server.DatabaseService = store.NewDatabaseService(m.l, db)
	server.TableService = store.NewTableService(m.l, db)
	server.DataSourceService = store.NewDataSourceService(m.l, db)
	server.IssueService = store.NewIssueService(m.l, db)
	server.PipelineService = store.NewPipelineService(m.l, db)
	server.StageService = store.NewStageService(m.l, db)
	server.TaskService = store.NewTaskService(m.l, db, store.NewTaskRunService(m.l, db))
	server.ActivityService = store.NewActivityService(m.l, db)
	server.BookmarkService = store.NewBookmarkService(m.l, db)
	server.VCSService = store.NewVCSService(m.l, db)

	m.server = server
	if err := server.Run(); err != nil {
		return err
	}

	return nil
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	m.l.Info("Trying to stop bytebase...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if m.server != nil {
		m.l.Info("Trying to gracefully shutdown server...")
		m.server.Shutdown(ctx)
	}

	if m.db != nil {
		m.l.Info("Trying to close database connections...")
		if err := m.db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	m := NewMain()

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := m.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		cancel()
	}()

	// Execute program.
	if err := m.Run(); err != nil {
		if err != http.ErrServerClosed {
			m.Close()
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	m.l.Info("Bytebase stopped properly.")
}
