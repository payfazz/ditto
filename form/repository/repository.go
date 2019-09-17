package repository

import (
	"context"

	"github.com/payfazz/canfazz-cms/internal/database/models/objectlist"

	"github.com/google/uuid"
	"github.com/payfazz/canfazz-cms/internal/database/models/dynamicoptions"
	"github.com/payfazz/canfazz-cms/internal/database/models/form"
	"github.com/payfazz/canfazz-cms/internal/database/models/listoptions"
	"github.com/payfazz/canfazz-cms/internal/database/models/supportedtype"
)

type Interface interface {
	GetFormByID(ctx context.Context, id uuid.UUID) (*form.Model, error)
	SaveForm(ctx context.Context, id uuid.UUID, name, description, typeForm, structure string, version int, published bool, actor string) error
	GetLatestVersion(ctx context.Context, typeForm string, id uuid.UUID, published bool) (*form.Model, error)
	Publish(ctx context.Context, id uuid.UUID, actor string) error
	GetSupportedTypeByTypeAndValue(ctx context.Context, typeStructure, value string) (*supportedtype.Model, error)
	GetSupportedTypeByValue(ctx context.Context, value string) (*supportedtype.Model, error)
	GetListOptionByType(ctx context.Context, code string) (*listoptions.Model, error)
	GetDynamicOptionByCode(ctx context.Context, code string) (*dynamicoptions.Model, error)
	GetObjectListByID(ctx context.Context, id uuid.UUID) (*objectlist.Model, error)
	GetObjectListByCode(ctx context.Context, code string) (*objectlist.Model, error)
	GetAllObjectListByCode(ctx context.Context, code string) ([]objectlist.Model, error)
}
