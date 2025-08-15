package bootstrap

import (
	"atm-test/internal/config"
	"atm-test/internal/pkg/cors"
	"atm-test/internal/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const gracefulDeadline = 5 * time.Second

type App struct {
	db       *sql.DB
	http     *http.Server
	cfg      config.Config
	log      logger.Logger
	ctx      context.Context
	teardown []func()
}

func New(cfg config.Config, log logger.Logger, ctx context.Context) *App {
	teardown := make([]func(), 0)

	app := App{
		cfg:      cfg,
		log:      log,
		teardown: teardown,
		ctx:      ctx,
	}

	app.initConnections()

	repo := initRepos(app.db, app.log)
	//useCase := initUseCases(
	//	repo,
	//	client,
	//	app.cfg,
	//	app.log)

	router := gin.Default()
	router.Use(cors.CORSMiddleware())

	//server := rest.New(
	//	router,
	//	app.log,
	//	useCase.authUserUc,
	//	useCase.authGoogleUc,
	//	useCase.authAnonymousUc,
	//	useCase.sessionsUc,
	//	useCase.cycleUserUc,
	//	useCase.cycleAnonymousUserUc,
	//	useCase.chatQuestionsUc,
	//	//useCase.chatAdminUc,
	//	useCase.chatAnswersUc,
	//	useCase.onboardingBlocksUC,
	//	useCase.onboardingQuestionsUC,
	//	useCase.onboardingOptionsUC,
	//	useCase.onboardingAnswersUC,
	//	// Partners
	//	useCase.partnersUC,
	//)

	app.http = &http.Server{
		Addr:        cfg.HTTPPort,
		Handler:     server,
		ReadTimeout: 10 * time.Second,
	}

	return &app
}

func (app *App) initConnections() {
	// Connecting to Database
	db, err := sql.Open("postgres", app.cfg.Postgres.PostgresURL())
	if err != nil {
		panic(fmt.Sprintf("sql.Open: %s", err))
	}
	app.log.Info("Database connection established")

	app.teardown = append(app.teardown, func() {
		app.log.Info("Database connection closing...")
		if err := db.Close(); err != nil {
			app.log.Error(err.Error())
		}
		app.log.Info("Database connection closed")
	})

	app.teardown = append(app.teardown, func() {
		app.log.Info("Firebase Messaging client cleanup")
	})

	app.db = db
	app.teardown = append(app.teardown, func() {
		app.log.Info("HTTP is shutting down")
		ctxShutDown, cancel := context.WithTimeout(app.ctx, gracefulDeadline)
		defer cancel()
		if err = app.http.Shutdown(ctxShutDown); err != nil {
			app.log.Error(fmt.Sprintf("server Shutdown Failed:%s", err))
			if err == http.ErrServerClosed {
				err = nil
			}
			return
		}
		app.log.Info("HTTP is shut down")
	})
}

func (app *App) Run(ctx context.Context) {
	go func() {
		app.log.Info("REST Server started at port " + app.cfg.HTTPPort)
		if err := app.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.log.Fatal(fmt.Sprintf("Failed To Run REST Server: %s\n", err.Error()))
		}
	}()
	<-ctx.Done()
	for i := range app.teardown {
		app.teardown[i]()
	}
}
