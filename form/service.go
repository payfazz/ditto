package form

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/payfazz/canfazz-cms/internal/form/repository"

	"github.com/payfazz/canfazz-cms/internal/lib/validate"

	"github.com/google/uuid"
	"github.com/payfazz/canfazz-authz/pkg/chi/middleware"
	formModel "github.com/payfazz/canfazz-cms/internal/database/models/form"
	"github.com/payfazz/canfazz-cms/internal/database/models/presentation"
	"github.com/payfazz/canfazz-cms/internal/form/model"
	rError "github.com/payfazz/kitx/pkg/error"
)

// Service ...
type Service struct {
	Actor    string
	validate validate.IValidate
	repo     repository.Interface
}

// New Service
func New(actor string, validate validate.IValidate, p repository.Interface) *Service {
	return &Service{
		Actor:    actor,
		validate: validate,
		repo:     p,
	}
}

//Create for create new form
func (formSvc *Service) Create(ctx context.Context, inputModel model.CreateFormInput) (*formModel.Model, error) {
	var err error
	var structure map[string]interface{}
	actor := formSvc.Actor
	userID := uuid.Nil
	authData, ok := middleware.GetAuthData(ctx).(map[string]interface{})
	if ok {
		userID = uuid.MustParse(authData["Sub"].(string))
	}
	if userID != uuid.Nil {
		actor = userID.String()
	}
	var version = 1
	err = json.Unmarshal([]byte(*inputModel.Structure), &structure)
	if err != nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_jsonb`)
	}
	err = formSvc.validateSection(ctx, structure, make(map[string]bool))
	if err != nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure`)
	}

	if *inputModel.Type == "kyc" {
		latestForm, err := formSvc.GetLatestVersion(ctx, *inputModel.Type, uuid.Nil, false)
		if err != nil {
			return nil, err
		}
		if latestForm != nil {
			version = latestForm.Version + 1
		}
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
	}

	err = formSvc.repo.SaveForm(ctx, newUUID, *inputModel.Name, *inputModel.Description, *inputModel.Type, *inputModel.Structure, version, false, actor)

	if err != nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `something_went_wrong`)
	}

	insertedData, err := formSvc.repo.GetFormByID(ctx, newUUID)

	if err != nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `something_went_wrong`)
	}

	return insertedData, nil
}

//GetLatestVersion for get latest version from type from and id
func (formSvc *Service) GetLatestVersion(ctx context.Context, typeForm string, id uuid.UUID, resolve bool) (*formModel.Model, error) {
	var err error

	if typeForm == "task" && id == uuid.Nil {
		return nil, nil
	}

	model, err := formSvc.repo.GetLatestVersion(ctx, typeForm, id, resolve)
	if err != nil {
		return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
	}

	if model == nil {
		return nil, rError.New(errors.New(`form_not_found`), rError.Enum.UNPROCESSABLEENTITY, `form_not_found`)
	}

	if resolve {
		var structure map[string]interface{}
		err = json.Unmarshal([]byte(model.Structure), &structure)
		if err != nil {
			return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_jsonb`)
		}
		err = formSvc.validateSection(ctx, structure, make(map[string]bool))
		if err != nil {
			return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure`)
		}

		structure, err = formSvc.resolveSection(ctx, structure)
		if err != nil {
			return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure`)
		}
		jsonStr, _ := json.Marshal(structure)
		model.Structure = string(jsonStr)
	}

	return model, nil
}

//Publish for publishing form
func (formSvc *Service) Publish(ctx context.Context, id uuid.UUID) error {
	var err error
	actor := formSvc.Actor
	userID := uuid.Nil
	authData, ok := middleware.GetAuthData(ctx).(map[string]interface{})
	if ok {
		userID = uuid.MustParse(authData["Sub"].(string))
	}
	if userID != uuid.Nil {
		actor = userID.String()
	}
	formData, err := formSvc.repo.GetFormByID(ctx, id)
	if err != nil {
		return rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
	}
	if formData == nil {
		return rError.New(errors.New(`invalid_form_id`), rError.Enum.UNPROCESSABLEENTITY, `invalid_form_id`)
	}

	err = formSvc.repo.Publish(ctx, id, actor)

	if err != nil {
		return rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `something_went_wrong`)
	}
	return nil
}

//ValidateInput for validate input form Endpoint
func (formSvc *Service) ValidateInput(ctx context.Context, inputedForm model.ValidateInput) (*presentation.ValidationInputModel, error) {
	userID := uuid.Nil
	userName := ""
	userCode := ""
	authData, ok := middleware.GetAuthData(ctx).(map[string]interface{})
	if ok {
		userID = uuid.MustParse(authData["Sub"].(string))
		userName = authData["SubName"].(string)
		userCode = authData["SubCode"].(string)
	}

	var err error
	reference := uuid.Nil
	if inputedForm.ID != nil {
		reference = *inputedForm.ID
	}
	latestVersion, err := formSvc.GetLatestVersion(ctx, inputedForm.TypeForm, reference, false)
	if err != nil {
		return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
	}

	if latestVersion == nil {
		return nil, rError.New(errors.New(`invalid_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_form`)
	}

	var structureForm map[string]interface{}
	err = json.Unmarshal([]byte(latestVersion.Structure), &structureForm)
	if err != nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_jsonb`)
	}
	var errFields []map[string]string
	structure, err := formSvc.validateAndGenerateFormInput(ctx, inputedForm.InputedForm, structureForm, &errFields)
	if err != nil && structure == nil {
		return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_input`)
	}

	var validationResult *presentation.ValidationInputModel

	var tmpAdditionalData map[string]interface{}
	if inputedForm.AdditionalData == nil {
		tmpAdditionalData = make(map[string]interface{})
		tmpAdditionalData["user_id"] = userID
		tmpAdditionalData["user_name"] = userName
		tmpAdditionalData["user_code"] = userCode
	} else {
		tmpAdditionalData = *inputedForm.AdditionalData
		tmpAdditionalData["user_id"] = userID
		tmpAdditionalData["user_name"] = userName
		tmpAdditionalData["user_code"] = userCode
	}
	inputedForm.AdditionalData = &tmpAdditionalData

	byteStructure, err := json.Marshal(structure)
	if err != nil {
		return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
	}
	validationResult = &presentation.ValidationInputModel{
		Version:        latestVersion.Version,
		Structure:      string(byteStructure),
		Value:          inputedForm.InputedForm,
		Type:           inputedForm.TypeForm,
		ID:             inputedForm.ID,
		AdditionalData: inputedForm.AdditionalData,
		FieldErrors:    errFields,
	}
	if structure["status"] != nil {
		validationResult.Valid = false
	} else {
		validationResult.Valid = true
		validationResult.ResolvedValue = formSvc.resolveValue(inputedForm.InputedForm, *inputedForm.AdditionalData)
	}

	return validationResult, nil
}

//extractArrayMap for extract from array interface to array map (object)
func (formSvc *Service) extractArrayMap(child []interface{}) []map[string]interface{} {
	var arrMap []map[string]interface{}
	var tmpMap map[string]interface{}
	for _, v := range child {
		tmpMap = v.(map[string]interface{})
		arrMap = append(arrMap, tmpMap)
	}
	return arrMap
}

//getAge for get age from birth day
func (formSvc *Service) getAge(birthday time.Time) int {
	now := time.Now()
	years := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		years--
	}
	return years
}
