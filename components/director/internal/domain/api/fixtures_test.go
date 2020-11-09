package api_test

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/domain/version"

	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/repo"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

const (
	apiDefID         = "ddddddddd-dddd-dddd-dddd-dddddddddddd"
	tenantID         = "ttttttttt-tttt-tttt-tttt-tttttttttttt"
	externalTenantID = "eeeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	packageID        = "ppppppppp-pppp-pppp-pppp-pppppppppppp"
)

func fixAPIDefinitionModel(id string, pkgID string, name, targetURL string) *model.APIDefinition {
	return &model.APIDefinition{
		ID:        id,
		PackageID: pkgID,
		Name:      name,
		TargetURL: targetURL,
	}
}

func fixFullAPIDefinitionModel(placeholder string) model.APIDefinition {
	spec := &model.APISpec{
		Data:   str.Ptr("spec_data_" + placeholder),
		Format: model.SpecFormatYaml,
		Type:   model.APISpecTypeOpenAPI,
	}

	deprecated := false
	forRemoval := false

	v := &model.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	return model.APIDefinition{
		ID:          apiDefID,
		Tenant:      tenantID,
		PackageID:   packageID,
		Name:        placeholder,
		Description: str.Ptr("desc_" + placeholder),
		Spec:        spec,
		TargetURL:   fmt.Sprintf("https://%s.com", placeholder),
		Group:       str.Ptr("group_" + placeholder),
		Version:     v,
	}
}

func fixGQLAPIDefinition(id string, pkgId string, name, targetURL string) *externalschema.APIDefinition {
	return &externalschema.APIDefinition{
		ID:        id,
		PackageID: pkgId,
		Name:      name,
		TargetURL: targetURL,
	}
}

func fixFullGQLAPIDefinition(placeholder string) *externalschema.APIDefinition {
	data := externalschema.CLOB("spec_data_" + placeholder)
	format := externalschema.SpecFormatYaml

	spec := &externalschema.APISpec{
		Data:         &data,
		Format:       format,
		Type:         externalschema.APISpecTypeOpenAPI,
		DefinitionID: apiDefID,
	}

	deprecated := false
	forRemoval := false

	v := &externalschema.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	return &externalschema.APIDefinition{
		ID:          apiDefID,
		PackageID:   packageID,
		Name:        placeholder,
		Description: str.Ptr("desc_" + placeholder),
		Spec:        spec,
		TargetURL:   fmt.Sprintf("https://%s.com", placeholder),
		Group:       str.Ptr("group_" + placeholder),
		Version:     v,
	}
}

func fixModelAPIDefinitionInput(name, description string, group string) *model.APIDefinitionInput {
	data := "data"

	spec := &model.APISpecInput{
		Data:         &data,
		Type:         model.APISpecTypeOpenAPI,
		Format:       model.SpecFormatYaml,
		FetchRequest: &model.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &model.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &model.APIDefinitionInput{
		Name:        name,
		Description: &description,
		TargetURL:   "https://test-url.com",
		Group:       &group,
		Spec:        spec,
		Version:     v,
	}
}

func fixGQLAPIDefinitionInput(name, description string, group string) *externalschema.APIDefinitionInput {
	data := externalschema.CLOB("data")

	spec := &externalschema.APISpecInput{
		Data:         &data,
		Type:         externalschema.APISpecTypeOpenAPI,
		Format:       externalschema.SpecFormatYaml,
		FetchRequest: &externalschema.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &externalschema.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &externalschema.APIDefinitionInput{
		Name:        name,
		Description: &description,
		TargetURL:   "https://test-url.com",
		Group:       &group,
		Spec:        spec,
		Version:     v,
	}
}

func fixModelAuthInput(headers map[string][]string) *model.AuthInput {
	return &model.AuthInput{
		AdditionalHeaders: headers,
	}
}

func fixGQLAuthInput(headers map[string][]string) *externalschema.AuthInput {
	httpHeaders := externalschema.HttpHeaders(headers)

	return &externalschema.AuthInput{
		AdditionalHeaders: &httpHeaders,
	}
}

func fixModelAuth() *model.Auth {
	return &model.Auth{
		Credential: model.CredentialData{
			Basic: &model.BasicCredentialData{
				Username: "foo",
				Password: "bar",
			},
		},
		AdditionalHeaders:     map[string][]string{"test": {"foo", "bar"}},
		AdditionalQueryParams: map[string][]string{"test": {"foo", "bar"}},
		RequestAuth: &model.CredentialRequestAuth{
			Csrf: &model.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: "foo.url",
				Credential: model.CredentialData{
					Basic: &model.BasicCredentialData{
						Username: "boo",
						Password: "far",
					},
				},
				AdditionalHeaders:     map[string][]string{"test": {"foo", "bar"}},
				AdditionalQueryParams: map[string][]string{"test": {"foo", "bar"}},
			},
		},
	}
}

func fixGQLAuth() *externalschema.Auth {
	return &externalschema.Auth{
		Credential: &externalschema.BasicCredentialData{
			Username: "foo",
			Password: "bar",
		},
		AdditionalHeaders:     &externalschema.HttpHeaders{"test": {"foo", "bar"}},
		AdditionalQueryParams: &externalschema.QueryParams{"test": {"foo", "bar"}},
		RequestAuth: &externalschema.CredentialRequestAuth{
			Csrf: &externalschema.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: "foo.url",
				Credential: &externalschema.BasicCredentialData{
					Username: "boo",
					Password: "far",
				},
				AdditionalHeaders:     &externalschema.HttpHeaders{"test": {"foo", "bar"}},
				AdditionalQueryParams: &externalschema.QueryParams{"test": {"foo", "bar"}},
			},
		},
	}
}

func fixModelAPIRtmAuth(id string, auth *model.Auth) *model.APIRuntimeAuth {
	return &model.APIRuntimeAuth{
		ID:        str.Ptr("foo"),
		TenantID:  "tnt",
		RuntimeID: id,
		APIDefID:  "api_id",
		Value:     auth,
	}
}

func fixEntityAPIDefinition(id string, pkgID string, name, targetUrl string) api.Entity {
	return api.Entity{
		ID:        id,
		PkgID:     pkgID,
		Name:      name,
		TargetURL: targetUrl,
	}
}

func fixFullEntityAPIDefinition(apiDefID, placeholder string) api.Entity {
	boolPlaceholder := false

	return api.Entity{
		ID:          apiDefID,
		TenantID:    tenantID,
		PkgID:       packageID,
		Name:        placeholder,
		Description: repo.NewValidNullableString("desc_" + placeholder),
		Group:       repo.NewValidNullableString("group_" + placeholder),
		TargetURL:   fmt.Sprintf("https://%s.com", placeholder),
		EntitySpec: api.EntitySpec{
			SpecData:   repo.NewValidNullableString("spec_data_" + placeholder),
			SpecFormat: repo.NewValidNullableString(string(model.SpecFormatYaml)),
			SpecType:   repo.NewValidNullableString(string(model.APISpecTypeOpenAPI)),
		},
		Version: version.Version{
			VersionValue:           repo.NewNullableString(str.Ptr("v1.1")),
			VersionDepracated:      repo.NewNullableBool(&boolPlaceholder),
			VersionDepracatedSince: repo.NewNullableString(str.Ptr("v1.0")),
			VersionForRemoval:      repo.NewNullableBool(&boolPlaceholder),
		},
	}
}

func fixAPIDefinitionColumns() []string {
	return []string{"id", "tenant_id", "package_id", "name", "description", "group_name", "target_url", "spec_data",
		"spec_format", "spec_type", "version_value", "version_deprecated",
		"version_deprecated_since", "version_for_removal"}
}

func fixAPIDefinitionRow(id, placeholder string) []driver.Value {
	return []driver.Value{id, tenantID, packageID, placeholder, "desc_" + placeholder, "group_" + placeholder,
		fmt.Sprintf("https://%s.com", placeholder), "spec_data_" + placeholder, "YAML", "OPEN_API",
		"v1.1", false, "v1.0", false}
}

func fixAPICreateArgs(id string, api *model.APIDefinition) []driver.Value {
	return []driver.Value{id, tenantID, packageID, api.Name, api.Description, api.Group,
		api.TargetURL, api.Spec.Data, string(api.Spec.Format), string(api.Spec.Type),
		api.Version.Value, api.Version.Deprecated, api.Version.DeprecatedSince,
		api.Version.ForRemoval}
}

func fixDefaultAuth() string {
	return `{"Credential":{"Basic":null,"Oauth":null},"AdditionalHeaders":{"testHeader":["hval1","hval2"]},"AdditionalQueryParams":null,"RequestAuth":null}`
}

func fixModelFetchRequest(id, url string, timestamp time.Time) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     id,
		Tenant: tenantID,
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: model.FetchRequestStatusConditionInitial,
			Timestamp: timestamp,
		},
		ObjectType: model.APIFetchRequestReference,
		ObjectID:   "foo",
	}
}

func fixModelFetchRequestWithCondition(id, url string, timestamp time.Time, condition model.FetchRequestStatusCondition) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     id,
		Tenant: tenantID,
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: condition,
			Timestamp: timestamp,
		},
		ObjectType: model.APIFetchRequestReference,
		ObjectID:   "foo",
	}
}

func fixGQLFetchRequest(url string, timestamp time.Time) *externalschema.FetchRequest {
	return &externalschema.FetchRequest{
		Filter: nil,
		Mode:   externalschema.FetchModeSingle,
		Auth:   nil,
		URL:    url,
		Status: &externalschema.FetchRequestStatus{
			Timestamp: externalschema.Timestamp(timestamp),
			Condition: externalschema.FetchRequestStatusConditionInitial,
		},
	}
}
