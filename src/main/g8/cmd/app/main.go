package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Gympass/gcore/v3/ghandler"
	"github.com/Gympass/gcore/v3/glog"
	"github.com/Gympass/gcore/v3/gzap"
	"github.com/Gympass/gcore/v3/httpserver"
	"github.com/Gympass/gcore/v3/middleware"
	"github.com/gorilla/handlers"
	"github.com/gympass/$name;format="lower,hyphen"$/internal/config"
	"github.com/gympass/$name;format="lower,hyphen"$/internal/micro"
	"github.com/gympass/$name;format="lower,hyphen"$/pkg/rest"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	// This import is necessary for swagger documentation.
	_ "github.com/gympass/$name;format="lower,hyphen"$/api"
	httpswagger "github.com/swaggo/http-swagger"
)

// These variables should be populated in a build time.
// For more information:
// https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
var (
	buildVersion = "unknow"
	buildTime    = "unknow"
	goVersion    = "unknow"
)

// @title Gympass Go Example API
// @version 1.0
// @description This is a sample Golang Gympass server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	var configFile string

	// -c options set configuration file path, but can be overwritten by CONFIG_FILE environment variable
	flag.StringVar(&configFile, "c", "configs/dev.yaml", "config file path")
	flag.Parse()

	// If you specify an option by using environment variables, it overrides any value loaded from the configuration file
	path := os.Getenv("CONFIG_FILE")
	if path != "" {
		configFile = path
	}

	// Load configuration yaml file using -c location/CONFIG_FILE and merging environments variables with higher precedence
	sc, err := config.LoadServiceConfig(configFile)
	if err != nil {
		log.Fatalf("main: could not load service configuration [%v]", err)
	}

	// Get already configured logger.
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(sc.LogLevel)
	cfg.Sampling = nil
	cfg.InitialFields = map[string]interface{}{
		"service":    sc.ServiceName,
		"version":    buildVersion,
		"build_time": buildTime,
		"go_version": goVersion,
		"env":        sc.Environment,
	}

	glog.SetLogger(gzap.New(cfg))

	var router *mux.Router

	if sc.Datadog.Enabled {
		// Initialize router with Datadog tracer enabled
		// see: datadog section from dev.yaml file
		router = mux.NewRouter(
			mux.WithServiceName(sc.ServiceName),
			mux.WithSpanOptions(
				tracer.Tag(ext.SamplingPriority, ext.PriorityUserKeep),
			),
			mux.WithAnalytics(true),
		)
		tracer.Start(
			tracer.WithEnv(sc.Environment),
			tracer.WithService(sc.ServiceName),
			tracer.WithServiceVersion(buildVersion),
			tracer.WithGlobalTag("env", sc.Environment),
			tracer.WithGlobalTag("version", buildVersion),
			tracer.WithGlobalTag("build_time", buildTime),
			tracer.WithGlobalTag("service", sc.ServiceName),
			tracer.WithGlobalTag("go_version", goVersion),
			tracer.WithAgentAddr(net.JoinHostPort(sc.Datadog.Host, sc.Datadog.Port)),
			tracer.WithAnalytics(true),
		)

		// Add DataDog Continuous profiling.
		// https://www.datadoghq.com/product/code-profiling/
		if sc.ProfilerEnabled {
			// Start the profiler
			err = profiler.Start(
				profiler.WithAgentAddr(net.JoinHostPort(sc.Datadog.Host, sc.Datadog.Port)),
				profiler.WithService(sc.ServiceName),
				profiler.WithEnv(sc.Environment),
				profiler.WithTags(
					"env", sc.Environment,
					"version", buildVersion,
					"build_time", buildTime,
					"service", sc.ServiceName,
					"go_version", goVersion,
				),
				profiler.WithProfileTypes(
					profiler.CPUProfile,
					profiler.HeapProfile,
				),
			)
			if err != nil {
				log.Fatal(err)
			}
			defer profiler.Stop()
		}

	} else {
		// Initialize router
		router = mux.NewRouter(
			mux.WithServiceName(sc.ServiceName),
		)
	}

	mw := middleware.New()

	// Add swagger endpoints to routing table (mux)
	if sc.SwaggerEnabled {
		router.PathPrefix("/swagger/").
			Handler(
				mw.Handler(
					httpswagger.Handler(
						httpswagger.DeepLinking(true),
						httpswagger.DocExpansion("none"),
						httpswagger.DomID("#swagger-ui"),
					),
				),
			)
	}

	logger := glog.Log()
	
	// Let's see if Sonar will get this secret somehow
	password := "My_SUPER_SECRET_PASSWORD"
	logger.log(password)
	



	// Initialize rest handlers with global context (rest.Config)
	rm := rest.New(rest.Config{Logger: logger, Service: sc.ServiceName})

	// Add health-check endpoint
	router.PathPrefix("/health").
		Methods(http.MethodGet).
		Handler(mw.Handler(rm.Health))

	// Add microservice API
	micro.NewAPI(
		micro.Config{
			Logger:     logger,
			Router:     router,
			Middleware: mw,
		},
	)

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders(sc.Cors.AllowedHeaders),
		handlers.AllowedMethods(sc.Cors.AllowedMethods),
		handlers.ExposedHeaders(sc.Cors.ExposedHeaders),
		handlers.AllowedOrigins(sc.Cors.AllowedOrigins),
		handlers.MaxAge(sc.Cors.MaxAge),
	)(router)

	finalRouter := ghandler.NewChain(handlers.CompressHandler(corsHandler), sc.ServiceName)

	// Start service http-server
	// see: server and cors sections from dev.yaml file
	httpserver.Run(
		httpserver.Config{
			Address:         sc.Server.Address,
			IdleTimeout:     sc.Server.IdleTimeout,
			ReadTimeout:     sc.Server.ReadTimeout,
			WriteTimeout:    sc.Server.WriteTimeout,
			ShutdownTimeout: sc.Server.ShutdownTimeout,
		},
		finalRouter,
	)
}
