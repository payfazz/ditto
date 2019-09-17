package http

import (
	nHttp "net/http"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/payfazz/canfazz-cms/internal/form/endpoint"
	"github.com/payfazz/fazzkit/server"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/payfazz/canfazz-cms/config"
	formSvc "github.com/payfazz/canfazz-cms/internal/form"
	"github.com/payfazz/canfazz-cms/internal/form/model"
)

var prometheusNameSpace = strings.ReplaceAll(strings.ToLower(config.ServiceName), "-", "_")
var prometheusSubsystem = "form_http"

// CreateForm ...
func CreateForm(svc *formSvc.Service, logger log.Logger, opts []kithttp.ServerOption) nHttp.Handler {
	end := endpoint.Create(svc)

	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: prometheusNameSpace,
			Subsystem: prometheusSubsystem,
			Action:    "create",
		}
	}

	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.CreateFormInput{},
		Logger:      serverLogger,
	}, opts...)
}

// PublishForm ...
func PublishForm(svc *formSvc.Service, logger log.Logger, opts []kithttp.ServerOption) nHttp.Handler {
	end := endpoint.Publish(svc)

	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: prometheusNameSpace,
			Subsystem: prometheusSubsystem,
			Action:    "publish",
		}
	}

	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.PublishInput{},
		Logger:      serverLogger,
	}, opts...)
}

// GetLatestForm ...
func GetLatestForm(svc *formSvc.Service, logger log.Logger, opts []kithttp.ServerOption) nHttp.Handler {
	end := endpoint.GetLatest(svc)

	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: prometheusNameSpace,
			Subsystem: prometheusSubsystem,
			Action:    "get_latest",
		}
	}

	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.GetLatestInput{},
		Logger:      serverLogger,
	}, opts...)
}

// MatchValueWithForm ...
func MatchValueWithForm(svc *formSvc.Service, logger log.Logger, opts []kithttp.ServerOption) nHttp.Handler {
	end := endpoint.Value(svc)

	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: prometheusNameSpace,
			Subsystem: prometheusSubsystem,
			Action:    "value",
		}
	}

	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.ValidateInput{},
		Logger:      serverLogger,
	}, opts...)
}

// Submit ...
func Submit(svc *formSvc.Service, logger log.Logger, opts []kithttp.ServerOption) nHttp.Handler {
	end := endpoint.Submit(svc)

	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: prometheusNameSpace,
			Subsystem: prometheusSubsystem,
			Action:    "submit",
		}
	}

	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.ValidateInput{},
		Logger:      serverLogger,
	}, opts...)
}
