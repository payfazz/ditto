package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/payfazz/canfazz-cms/internal/global"

	"github.com/google/uuid"
	"github.com/payfazz/canfazz-cms/internal/database/models/dynamicoptions"
	"github.com/payfazz/canfazz-cms/internal/database/models/form"
	"github.com/payfazz/canfazz-cms/internal/database/models/listoptions"
	"github.com/payfazz/canfazz-cms/internal/database/models/objectlist"
	"github.com/payfazz/canfazz-cms/internal/database/models/supportedtype"
)

//Postgres ...
type Postgres struct{}

//NewPostgres ...
func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) GetFormByID(ctx context.Context, id uuid.UUID) (*form.Model, error) {
	return form.FetchByID(ctx, *global.GetQuery(ctx), id)
}

func (p *Postgres) SaveForm(ctx context.Context, id uuid.UUID, name, description, typeForm, structure string, version int, published bool, actor string) error {
	var query strings.Builder
	var err error

	query.WriteString(fmt.Sprintf(`INSERT INTO form (%s,%s,%s,%s,%s,%s,%s,%s`,
		form.Columns.ID,
		form.Columns.Name,
		form.Columns.Description,
		form.Columns.Type,
		form.Columns.Version,
		form.Columns.Structure,
		form.Columns.Published,
		form.Columns.CreatedBy))
	query.WriteString(`) VALUES (
		:id,
		:name,
		:description,
		:type,
		:version,
		:structure,
		false,
		:createdby)
		`)
	params := make(map[string]interface{})
	params["id"] = id
	params["name"] = name
	params["description"] = description
	params["type"] = typeForm
	params["version"] = version
	params["structure"] = structure
	params["createdby"] = actor
	if ctx == nil {
		_, err = global.GetQuery(ctx).NamedExec(query.String(), params)
	} else {
		_, err = global.GetQuery(ctx).NamedExecContext(ctx, query.String(), params)
	}

	if err != nil {
		return err
	}
	return form.Log(ctx, *global.GetQuery(ctx), id)
}

func (p *Postgres) GetLatestVersion(ctx context.Context, typeForm string, id uuid.UUID, published bool) (*form.Model, error) {
	var query strings.Builder
	var params []interface{}
	var row *sqlx.Row
	var err error
	query.WriteString(fmt.Sprintf(`SELECT %s FROM form WHERE %s=$1`, form.FieldsString(""), form.Columns.Type))
	params = append(params, typeForm)
	if id != uuid.Nil {
		fmt.Println("id", id)
		query.WriteString(fmt.Sprintf(` AND %s=$2`, form.Columns.ID))
		params = append(params, id)
	}
	query.WriteString(fmt.Sprintf(` AND %s IS NULL`, form.Columns.DeletedDateUTC))

	if published {
		query.WriteString(fmt.Sprintf(` AND %s=true`, form.Columns.Published))
	}

	query.WriteString(fmt.Sprintf(` ORDER BY %s DESC LIMIT 1`, form.Columns.Version))

	model := &form.Model{}

	if ctx == nil {
		row = global.GetQuery(ctx).QueryRowx(query.String(), params...)
	} else {
		row = global.GetQuery(ctx).QueryRowxContext(ctx, query.String(), params...)
	}
	err = model.Scan(row)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return model, nil
}

func (p *Postgres) Publish(ctx context.Context, id uuid.UUID, actor string) error {
	params := make(map[string]interface{})
	var err error
	query := fmt.Sprintf(`UPDATE %s SET %s=true, %s=current_timestamp, %s=current_timestamp, %s=:updatedby WHERE %s=:id`, form.TableName, form.Columns.Published, form.Columns.PublishedDateUTC, form.Columns.UpdatedDateUTC, form.Columns.UpdatedBy, form.Columns.ID)
	params["id"] = id
	params["updatedby"] = actor
	if ctx == nil {
		_, err = global.GetQuery(ctx).NamedExec(query, params)
	} else {
		_, err = global.GetQuery(ctx).NamedExecContext(ctx, query, params)
	}

	if err != nil {
		return err
	}
	return form.Log(ctx, *global.GetQuery(ctx), id)
}

func (p *Postgres) GetSupportedTypeByTypeAndValue(ctx context.Context, typeStructure, value string) (*supportedtype.Model, error) {
	var err error
	var query strings.Builder
	var row *sqlx.Row
	var params []interface{}
	query.WriteString(`SELECT `)
	query.WriteString(supportedtype.FieldsString(""))
	query.WriteString(` FROM `)
	query.WriteString(supportedtype.TableName)
	query.WriteString(fmt.Sprintf(` WHERE %s=$1 AND %s=$2 AND %s IS NULL`, supportedtype.Columns.Type, supportedtype.Columns.Value, supportedtype.Columns.DeletedDateUTC))
	params = append(params, typeStructure)
	params = append(params, value)

	model := &supportedtype.Model{}

	if ctx == nil {
		row = global.GetQuery(ctx).QueryRowx(query.String(), params...)
	} else {
		row = global.GetQuery(ctx).QueryRowxContext(ctx, query.String(), params...)
	}
	err = model.Scan(row)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return model, nil
}

func (p *Postgres) GetListOptionByType(ctx context.Context, code string) (*listoptions.Model, error) {
	return listoptions.FetchByCode(ctx, *global.GetQuery(ctx), code)
}

func (p *Postgres) GetDynamicOptionByCode(ctx context.Context, code string) (*dynamicoptions.Model, error) {
	return dynamicoptions.FetchByCode(ctx, *global.GetQuery(ctx), code)
}

func (p *Postgres) GetObjectListByID(ctx context.Context, id uuid.UUID) (*objectlist.Model, error) {
	return objectlist.FetchByID(ctx, *global.GetQuery(ctx), id)
}

func (p *Postgres) GetObjectListByCode(ctx context.Context, code string) (*objectlist.Model, error) {
	return objectlist.FetchByCode(ctx, *global.GetQuery(ctx), code)
}

func (p *Postgres) GetSupportedTypeByValue(ctx context.Context, value string) (*supportedtype.Model, error) {
	return supportedtype.FetchByValue(ctx, *global.GetQuery(ctx), value)
}

func (p *Postgres) GetAllObjectListByCode(ctx context.Context, code string) ([]objectlist.Model, error) {
	var err error
	var rows *sqlx.Rows
	var query strings.Builder
	var params []interface{}
	query.WriteString(`SELECT ` + objectlist.FieldsString("") + ` FROM ` + objectlist.TableName + ` WHERE ` + objectlist.Columns.Code + `=$1 AND ` + objectlist.Columns.DeletedDateUTC + ` IS NULL`)
	params = append(params, code)
	if ctx != nil {
		rows, err = global.GetQuery(ctx).QueryxContext(ctx, query.String(), params...)
	} else {
		rows, err = global.GetQuery(ctx).Queryx(query.String(), params...)
	}

	if err != nil {
		return nil, err
	}

	var result []objectlist.Model
	for rows.Next() {
		tmpModel := objectlist.Model{}
		err = tmpModel.Scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, tmpModel)
	}

	return result, nil
}
