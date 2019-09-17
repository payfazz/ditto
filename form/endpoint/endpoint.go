package endpoint

import (
	"context"

	"github.com/payfazz/kitx/pkg/event"

	"github.com/payfazz/canfazz-cms/internal/global"
	"github.com/payfazz/kitx/pkg/db"

	"github.com/payfazz/canfazz-cms/config"
	"github.com/payfazz/canfazz-cms/internal/form/model"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/canfazz-cms/internal/form"
)

func GetLatest(s *form.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData := request.(*model.GetLatestInput)
		return s.GetLatestVersion(ctx, reqData.TypeForm, reqData.ID, true)
	}
}

func Publish(s *form.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData := request.(*model.PublishInput)
		err = db.RunInTransaction(ctx, global.GetQuery(ctx).DB(), func(ctx context.Context) error {
			err = s.Publish(ctx, reqData.ID)
			return err
		})

		return response, err
	}
}

func Value(s *form.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData := request.(*model.ValidateInput)
		return s.ValidateInput(ctx, *reqData)
	}
}

func Submit(s *form.Service) endpoint.Endpoint {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData := request.(*model.ValidateInput)
		err = db.RunInTransaction(ctx, global.GetQuery(ctx).DB(), func(ctx context.Context) error {
			response, err = s.ValidateInput(ctx, *reqData)
			return err
		})
		if err != nil {
			return nil, err
		}
		dataStore := event.StoreDetail{
			Domain:      config.ServiceName,
			Subject:     cfg.PubFormSubmitted,
			EventSource: "",
			Data:        response,
		}
		event.Store(dataStore)
		return response, err
	}
}

func Create(s *form.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData := request.(*model.CreateFormInput)
		err = db.RunInTransaction(ctx, global.GetQuery(ctx).DB(), func(ctx context.Context) error {
			response, err = s.Create(ctx, *reqData)
			return err
		})

		return response, err
	}
}
