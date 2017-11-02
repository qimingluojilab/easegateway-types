package plugins

import (
	"net/http"

	"github.com/hexdecteam/easegateway-types/pipelines"
	"github.com/hexdecteam/easegateway-types/task"
)

// Plugin needs to cover follow rules:
//
// 1. Run(task.Task) method returns error only if
//    a) the plugin needs reconstruction, e.g. backend failure causes local client object invalidation;
//    b) the task has been cancelled by pipeline after running plugin is updated dynamically, task will
//    re-run on updated plugin;
//    The error caused by user input should be updated to task instead.
// 2. Should be implemented as stateless and be re-entry-able (idempotency) on the same task, a plugin
//    instance could be used in different pipeline or parallel running instances of same pipeline.
// 3. Prepare(pipelines.PipelineContext) guarantees it will be called on the same pipeline context against
//    the same plugin instance only once before executing Run(task.Task) on the pipeline.
type Plugin interface {
	Prepare(ctx pipelines.PipelineContext)
	Run(ctx pipelines.PipelineContext, t task.Task) (task.Task, error)
	Name() string
	CleanUp(ctx pipelines.PipelineContext)
	Close()
}

type Constructor func(conf Config) (Plugin, error)

type Config interface {
	PluginName() string
	Prepare(pipelineNames []string) error
}

type ConfigConstructor func() Config

////

type HTTPHandler func(w http.ResponseWriter, r *http.Request, path_params map[string]string)

type HTTPMuxEntry struct {
	Headers map[string][]string
	Handler HTTPHandler
}

type HTTPMux interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	AddFunc(pipeline, path, method string, headers map[string][]string, handler HTTPHandler) error
	AddFuncs(pipeline string, pipeline_rtable map[string]map[string]*HTTPMuxEntry) error
	DeleteFunc(pipeline, path, method string)
	DeleteFuncs(pipeline string) map[string]map[string]*HTTPMuxEntry
}

const (
	HTTP_SERVER_MUX_BUCKET_KEY                  = "HTTP_SERVER_MUX_BUCKET_KEY"
	HTTP_SERVER_PIPELINE_ROUTE_TABLE_BUCKET_KEY = "HTTP_SERVER_PIPELINE_ROUTE_TABLE_BUCKET_KEY"
	HTTP_SERVER_GONE_NOTIFIER_BUCKET_KEY        = "HTTP_SERVER_GONE_NOTIFIER_BUCKET_KEY"
)
