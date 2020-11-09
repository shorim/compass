package document

import (
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

type converter struct {
	frConverter FetchRequestConverter
}

func NewConverter(frConverter FetchRequestConverter) *converter {
	return &converter{frConverter: frConverter}
}

func (c *converter) ToGraphQL(in *model.Document) *externalschema.Document {
	if in == nil {
		return nil
	}

	var clob *externalschema.CLOB
	if in.Data != nil {
		tmp := externalschema.CLOB([]byte(*in.Data))
		clob = &tmp
	}

	return &externalschema.Document{
		ID:          in.ID,
		PackageID:   in.PackageID,
		Title:       in.Title,
		DisplayName: in.DisplayName,
		Description: in.Description,
		Format:      externalschema.DocumentFormat(in.Format),
		Kind:        in.Kind,
		Data:        clob,
	}
}

func (c *converter) MultipleToGraphQL(in []*model.Document) []*externalschema.Document {
	var documents []*externalschema.Document
	for _, r := range in {
		if r == nil {
			continue
		}

		documents = append(documents, c.ToGraphQL(r))
	}

	return documents
}

func (c *converter) InputFromGraphQL(in *externalschema.DocumentInput) (*model.DocumentInput, error) {
	if in == nil {
		return nil, nil
	}

	var data *string
	if in.Data != nil {
		tmp := string(*in.Data)
		data = &tmp
	}

	fetchReq, err := c.frConverter.InputFromGraphQL(in.FetchRequest)
	if err != nil {
		return nil, errors.Wrap(err, "while converting FetchRequestInput input")
	}

	return &model.DocumentInput{
		Title:        in.Title,
		DisplayName:  in.DisplayName,
		Description:  in.Description,
		Format:       model.DocumentFormat(in.Format),
		Kind:         in.Kind,
		Data:         data,
		FetchRequest: fetchReq,
	}, nil
}

func (c *converter) MultipleInputFromGraphQL(in []*externalschema.DocumentInput) ([]*model.DocumentInput, error) {
	var inputs []*model.DocumentInput
	for _, r := range in {
		if r == nil {
			continue
		}

		docInput, err := c.InputFromGraphQL(r)
		if err != nil {
			return nil, err
		}

		inputs = append(inputs, docInput)
	}

	return inputs, nil
}

func (c *converter) ToEntity(in model.Document) (Entity, error) {
	kind := repo.NewNullableString(in.Kind)
	data := repo.NewNullableString(in.Data)

	out := Entity{
		ID:          in.ID,
		PkgID:       in.PackageID,
		TenantID:    in.Tenant,
		Title:       in.Title,
		DisplayName: in.DisplayName,
		Description: in.Description,
		Format:      string(in.Format),
		Kind:        kind,
		Data:        data,
	}

	return out, nil
}

func (c *converter) FromEntity(in Entity) (model.Document, error) {
	kind := repo.StringPtrFromNullableString(in.Kind)
	data := repo.StringPtrFromNullableString(in.Data)

	out := model.Document{
		ID:          in.ID,
		PackageID:   in.PkgID,
		Tenant:      in.TenantID,
		Title:       in.Title,
		DisplayName: in.DisplayName,
		Description: in.Description,
		Format:      model.DocumentFormat(in.Format),
		Kind:        kind,
		Data:        data,
	}
	return out, nil
}
