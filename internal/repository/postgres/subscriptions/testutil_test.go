package subscriptions_test

import (
	"context"
	"os"
	"testing"
	"time"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository/postgres/subscriptions"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB *pgxpool.Pool
	repo   *subscriptions.Repository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",

		postgres.WithDatabase("subscriptions_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	runMigrations(connStr)

	testDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}

	repo = subscriptions.NewRepository(testDB)

	code := m.Run()

	testDB.Close()

	_ = container.Terminate(ctx)

	os.Exit(code)
}

func runMigrations(connStr string) {
	m, err := migrate.New("file://../../../../migrations", connStr)
	if err != nil {
		panic(err)
	}

	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil &&
		err != migrate.ErrNoChange {

		panic(err)
	}
}

func cleanupTables(t *testing.T) {
	t.Helper()

	queries := []string{"TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE"}

	for _, q := range queries {
		_, err := testDB.Exec(context.Background(), q)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupTest(t *testing.T) {
	t.Helper()

	cleanupTables(t)
}

func createTestSubscription(t *testing.T, repo *subscriptions.Repository) model.Subscription {
	t.Helper()

	endDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	sub := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     &endDate,
	}

	id, err := repo.Create(context.Background(), sub)
	if err != nil {
		t.Fatal(err)
	}

	sub.ID = id

	return sub
}
