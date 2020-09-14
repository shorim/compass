// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
)

type CredentialData interface {
	IsCredentialData()
}

type OneTimeToken interface {
	IsOneTimeToken()
}

// Every query that implements pagination returns object that implements Pageable interface.
// To specify page details, query specify two parameters: `first` and `after`.
// `first` specify page size, `after` is a cursor for the next page. When requesting first page, set `after` to empty value.
// For requesting next page, set `after` to `pageInfo.endCursor` returned from previous query.
type Pageable interface {
	IsPageable()
}

type APIDefinitionInput struct {
	// TODO: Validation if it is guid
	ID              *string `json:"ID"`
	OpenDiscoveryID *string `json:"openDiscoveryID"`
	// **Validation:** ASCII printable characters, max=100
	Title            string `json:"title"`
	ShortDescription string `json:"shortDescription"`
	// **Validation:** max=2000
	Description *string `json:"description"`
	// **Validation:** valid URL, max=512
	EntryPoint string `json:"entryPoint"`
	// **Validation:** max=36
	Group            *string         `json:"group"`
	Specs            []*APISpecInput `json:"specs"`
	Version          *VersionInput   `json:"version"`
	APIDefinitions   JSON            `json:"apiDefinitions"`
	Tags             *JSON           `json:"tags"`
	Documentation    *string         `json:"documentation"`
	ChangelogEntries *JSON           `json:"changelogEntries"`
	Logo             *string         `json:"logo"`
	Image            *string         `json:"image"`
	URL              *string         `json:"url"`
	// should be ENUM
	ReleaseStatus string `json:"releaseStatus"`
	// should be ENUM
	APIProtocol string    `json:"apiProtocol"`
	Actions     JSON      `json:"actions"`
	LastUpdated Timestamp `json:"lastUpdated"`
	Extensions  *JSON     `json:"extensions"`
}

type APIDefinitionPage struct {
	Data       []*APIDefinition `json:"data"`
	PageInfo   *PageInfo        `json:"pageInfo"`
	TotalCount int              `json:"totalCount"`
}

func (APIDefinitionPage) IsPageable() {}

// **Validation:**
// - for ODATA type, accepted formats are XML and JSON, for OPEN_API accepted formats are YAML and JSON
// - data or fetchRequest required
type APISpecInput struct {
	Data         *CLOB              `json:"data"`
	Type         APISpecType        `json:"type"`
	CustomType   *string            `json:"customType"`
	Format       SpecFormat         `json:"format"`
	FetchRequest *FetchRequestInput `json:"fetchRequest"`
}

type ApplicationEventingConfiguration struct {
	DefaultURL string `json:"defaultURL"`
}

// **Validation:** provided placeholders' names are unique
type ApplicationFromTemplateInput struct {
	// **Validation:** ASCII printable characters, max=100
	TemplateName string                `json:"templateName"`
	Values       []*TemplateValueInput `json:"values"`
}

type ApplicationPage struct {
	Data       []*Application `json:"data"`
	PageInfo   *PageInfo      `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (ApplicationPage) IsPageable() {}

type ApplicationRegisterInput struct {
	// **Validation:**  Up to 36 characters long. Cannot start with a digit. The characters allowed in names are: digits (0-9), lower case letters (a-z),-, and .
	Name string `json:"name"`
	// **Validation:** max=256
	ProviderName *string `json:"providerName"`
	// **Validation:** max=2000
	Description *string `json:"description"`
	// **Validation:** label key is alphanumeric with underscore
	Labels   *Labels         `json:"labels"`
	Webhooks []*WebhookInput `json:"webhooks"`
	// **Validation:** valid URL, max=256
	HealthCheckURL      *string                     `json:"healthCheckURL"`
	Bundles             []*BundleInput              `json:"bundles"`
	IntegrationSystemID *string                     `json:"integrationSystemID"`
	StatusCondition     *ApplicationStatusCondition `json:"statusCondition"`
}

type ApplicationStatus struct {
	Condition ApplicationStatusCondition `json:"condition"`
	Timestamp Timestamp                  `json:"timestamp"`
}

type ApplicationTemplate struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Description      *string                        `json:"description"`
	ApplicationInput string                         `json:"applicationInput"`
	Placeholders     []*PlaceholderDefinition       `json:"placeholders"`
	AccessLevel      ApplicationTemplateAccessLevel `json:"accessLevel"`
}

// **Validation:** provided placeholders' names are unique and used in applicationInput
type ApplicationTemplateInput struct {
	// **Validation:** ASCII printable characters, max=100
	Name string `json:"name"`
	// **Validation:** max=2000
	Description      *string                        `json:"description"`
	ApplicationInput *ApplicationRegisterInput      `json:"applicationInput"`
	Placeholders     []*PlaceholderDefinitionInput  `json:"placeholders"`
	AccessLevel      ApplicationTemplateAccessLevel `json:"accessLevel"`
}

type ApplicationTemplatePage struct {
	Data       []*ApplicationTemplate `json:"data"`
	PageInfo   *PageInfo              `json:"pageInfo"`
	TotalCount int                    `json:"totalCount"`
}

func (ApplicationTemplatePage) IsPageable() {}

type ApplicationUpdateInput struct {
	// **Validation:** max=256
	ProviderName *string `json:"providerName"`
	// **Validation:** max=2000
	Description *string `json:"description"`
	// **Validation:** valid URL, max=256
	HealthCheckURL      *string                     `json:"healthCheckURL"`
	IntegrationSystemID *string                     `json:"integrationSystemID"`
	StatusCondition     *ApplicationStatusCondition `json:"statusCondition"`
}

type Auth struct {
	Credential                      CredentialData         `json:"credential"`
	AdditionalHeaders               *HttpHeaders           `json:"additionalHeaders"`
	AdditionalHeadersSerialized     *HttpHeadersSerialized `json:"additionalHeadersSerialized"`
	AdditionalQueryParams           *QueryParams           `json:"additionalQueryParams"`
	AdditionalQueryParamsSerialized *QueryParamsSerialized `json:"additionalQueryParamsSerialized"`
	RequestAuth                     *CredentialRequestAuth `json:"requestAuth"`
}

type AuthInput struct {
	Credential *CredentialDataInput `json:"credential"`
	// **Validation:** if provided, headers name and value required
	AdditionalHeaders           *HttpHeaders           `json:"additionalHeaders"`
	AdditionalHeadersSerialized *HttpHeadersSerialized `json:"additionalHeadersSerialized"`
	// **Validation:** if provided, query parameters name and value required
	AdditionalQueryParams           *QueryParams                `json:"additionalQueryParams"`
	AdditionalQueryParamsSerialized *QueryParamsSerialized      `json:"additionalQueryParamsSerialized"`
	RequestAuth                     *CredentialRequestAuthInput `json:"requestAuth"`
}

type AutomaticScenarioAssignment struct {
	ScenarioName string `json:"scenarioName"`
	Selector     *Label `json:"selector"`
}

type AutomaticScenarioAssignmentPage struct {
	Data       []*AutomaticScenarioAssignment `json:"data"`
	PageInfo   *PageInfo                      `json:"pageInfo"`
	TotalCount int                            `json:"totalCount"`
}

func (AutomaticScenarioAssignmentPage) IsPageable() {}

type AutomaticScenarioAssignmentSetInput struct {
	ScenarioName string `json:"scenarioName"`
	// Runtimes and Applications which contain labels with equal key and value are matched
	Selector *LabelSelectorInput `json:"selector"`
}

type BasicCredentialData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (BasicCredentialData) IsCredentialData() {}

type BasicCredentialDataInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BundleInput struct {
	// TODO: Validation if it is guid
	ID              *string `json:"ID"`
	OpenDiscoveryID *string `json:"openDiscoveryID"`
	// **Validation:** ASCII printable characters, max=100
	Title string `json:"title"`
	// **Validation:** max=2000
	ShortDescription               string                  `json:"shortDescription"`
	Description                    *string                 `json:"description"`
	Tags                           *JSON                   `json:"tags"`
	LastUpdated                    Timestamp               `json:"lastUpdated"`
	Extensions                     *JSON                   `json:"extensions"`
	InstanceAuthRequestInputSchema *JSONSchema             `json:"instanceAuthRequestInputSchema"`
	DefaultInstanceAuth            *AuthInput              `json:"defaultInstanceAuth"`
	APIDefinitions                 []*APIDefinitionInput   `json:"apiDefinitions"`
	EventDefinitions               []*EventDefinitionInput `json:"eventDefinitions"`
	Documents                      []*DocumentInput        `json:"documents"`
}

type BundleInstanceAuth struct {
	ID string `json:"id"`
	// Context of BundleInstanceAuth - such as Runtime ID, namespace
	Context *JSON `json:"context"`
	// User input while requesting Bundle Instance Auth
	InputParams *JSON `json:"inputParams"`
	// It may be empty if status is PENDING.
	// Populated with `bundle.defaultAuth` value if `bundle.defaultAuth` is defined. If not, Compass notifies Application/Integration System about the Auth request.
	Auth   *Auth                     `json:"auth"`
	Status *BundleInstanceAuthStatus `json:"status"`
}

type BundleInstanceAuthRequestInput struct {
	// Context of BundleInstanceAuth - such as Runtime ID, namespace, etc.
	Context *JSON `json:"context"`
	// **Validation:** JSON validated against bundle.instanceAuthRequestInputSchema
	InputParams *JSON `json:"inputParams"`
}

type BundleInstanceAuthSetInput struct {
	// **Validation:** If not provided, the status has to be set. If provided, the status condition  must be "SUCCEEDED".
	Auth *AuthInput `json:"auth"`
	// **Validation:** Optional if the auth is provided.
	// If the status condition is "FAILED", auth must be empty.
	Status *BundleInstanceAuthStatusInput `json:"status"`
}

type BundleInstanceAuthStatus struct {
	Condition BundleInstanceAuthStatusCondition `json:"condition"`
	Timestamp Timestamp                         `json:"timestamp"`
	Message   string                            `json:"message"`
	// Possible reasons:
	// - PendingNotification
	// - NotificationSent
	// - CredentialsProvided
	// - CredentialsNotProvided
	// - PendingDeletion
	Reason string `json:"reason"`
}

type BundleInstanceAuthStatusInput struct {
	Condition BundleInstanceAuthSetStatusConditionInput `json:"condition"`
	// **Validation:** required, if condition is FAILED
	Message string `json:"message"`
	// Example reasons:
	// - PendingNotification
	// - NotificationSent
	// - CredentialsProvided
	// - CredentialsNotProvided
	// - PendingDeletion
	//
	//    **Validation**: required, if condition is FAILED
	Reason string `json:"reason"`
}

type BundlePackageRelationInput struct {
	PackageID string `json:"packageID"`
	BundleID  string `json:"bundleID"`
}

type BundlePage struct {
	Data       []*Bundle `json:"data"`
	PageInfo   *PageInfo `json:"pageInfo"`
	TotalCount int       `json:"totalCount"`
}

func (BundlePage) IsPageable() {}

type CSRFTokenCredentialRequestAuth struct {
	TokenEndpointURL                string                 `json:"tokenEndpointURL"`
	Credential                      CredentialData         `json:"credential"`
	AdditionalHeaders               *HttpHeaders           `json:"additionalHeaders"`
	AdditionalHeadersSerialized     *HttpHeadersSerialized `json:"additionalHeadersSerialized"`
	AdditionalQueryParams           *QueryParams           `json:"additionalQueryParams"`
	AdditionalQueryParamsSerialized *QueryParamsSerialized `json:"additionalQueryParamsSerialized"`
}

type CSRFTokenCredentialRequestAuthInput struct {
	// **Validation:** valid URL
	TokenEndpointURL string               `json:"tokenEndpointURL"`
	Credential       *CredentialDataInput `json:"credential"`
	// **Validation:** if provided, headers name and value required
	AdditionalHeaders           *HttpHeaders           `json:"additionalHeaders"`
	AdditionalHeadersSerialized *HttpHeadersSerialized `json:"additionalHeadersSerialized"`
	// **Validation:** if provided, query parameters name and value required
	AdditionalQueryParams           *QueryParams           `json:"additionalQueryParams"`
	AdditionalQueryParamsSerialized *QueryParamsSerialized `json:"additionalQueryParamsSerialized"`
}

// **Validation:** basic or oauth field required
type CredentialDataInput struct {
	Basic *BasicCredentialDataInput `json:"basic"`
	Oauth *OAuthCredentialDataInput `json:"oauth"`
}

type CredentialRequestAuth struct {
	Csrf *CSRFTokenCredentialRequestAuth `json:"csrf"`
}

type CredentialRequestAuthInput struct {
	// **Validation:** required
	Csrf *CSRFTokenCredentialRequestAuthInput `json:"csrf"`
}

type DocumentInput struct {
	// TODO: Validation if it is guid
	ID *string `json:"ID"`
	// **Validation:** max=128
	Title string `json:"title"`
	// **Validation:** max=128
	DisplayName string `json:"displayName"`
	// **Validation:** max=2000
	Description string         `json:"description"`
	Format      DocumentFormat `json:"format"`
	// **Validation:** max=256
	Kind         *string            `json:"kind"`
	Data         *CLOB              `json:"data"`
	FetchRequest *FetchRequestInput `json:"fetchRequest"`
}

type DocumentPage struct {
	Data       []*Document `json:"data"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

func (DocumentPage) IsPageable() {}

type EventDefinitionInput struct {
	// TODO: Validation if it is guid
	ID              *string `json:"ID"`
	OpenDiscoveryID *string `json:"openDiscoveryID"`
	// **Validation:** ASCII printable characters, max=100
	Title            string `json:"title"`
	ShortDescription string `json:"shortDescription"`
	// **Validation:** max=2000
	Description *string           `json:"description"`
	Specs       []*EventSpecInput `json:"specs"`
	// **Validation:** max=36
	Group            *string       `json:"group"`
	Version          *VersionInput `json:"version"`
	EventDefinitions JSON          `json:"eventDefinitions"`
	Tags             *JSON         `json:"tags"`
	Documentation    *string       `json:"documentation"`
	ChangelogEntries *JSON         `json:"changelogEntries"`
	Logo             *string       `json:"logo"`
	Image            *string       `json:"image"`
	URL              *string       `json:"url"`
	// should be ENUM
	ReleaseStatus string    `json:"releaseStatus"`
	LastUpdated   Timestamp `json:"lastUpdated"`
	Extensions    *JSON     `json:"extensions"`
}

type EventDefinitionPage struct {
	Data       []*EventDefinition `json:"data"`
	PageInfo   *PageInfo          `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

func (EventDefinitionPage) IsPageable() {}

// **Validation:**
// - data or fetchRequest required
// - for ASYNC_API type, accepted formats are YAML and JSON
type EventSpecInput struct {
	Data         *CLOB              `json:"data"`
	Type         EventSpecType      `json:"type"`
	CustomType   *string            `json:"customType"`
	Format       SpecFormat         `json:"format"`
	FetchRequest *FetchRequestInput `json:"fetchRequest"`
}

// Compass performs fetch to validate if request is correct and stores a copy
type FetchRequest struct {
	URL    string              `json:"url"`
	Auth   *Auth               `json:"auth"`
	Mode   FetchMode           `json:"mode"`
	Filter *string             `json:"filter"`
	Status *FetchRequestStatus `json:"status"`
}

type FetchRequestInput struct {
	// **Validation:** valid URL, max=256
	URL string `json:"url"`
	// Currently unsupported, providing it will result in a failure
	Auth *AuthInput `json:"auth"`
	// Currently unsupported, providing it will result in a failure
	Mode *FetchMode `json:"mode"`
	// **Validation:** max=256
	// Currently unsupported, providing it will result in a failure
	Filter *string `json:"filter"`
}

type FetchRequestStatus struct {
	Condition FetchRequestStatusCondition `json:"condition"`
	Message   *string                     `json:"message"`
	Timestamp Timestamp                   `json:"timestamp"`
}

type HealthCheck struct {
	Type      HealthCheckType            `json:"type"`
	Condition HealthCheckStatusCondition `json:"condition"`
	Origin    *string                    `json:"origin"`
	Message   *string                    `json:"message"`
	Timestamp Timestamp                  `json:"timestamp"`
}

type HealthCheckPage struct {
	Data       []*HealthCheck `json:"data"`
	PageInfo   *PageInfo      `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (HealthCheckPage) IsPageable() {}

type IntegrationSystemInput struct {
	// **Validation:**  Up to 36 characters long. Cannot start with a digit. The characters allowed in names are: digits (0-9), lower case letters (a-z),-, and .
	Name string `json:"name"`
	// **Validation:** max=2000
	Description *string `json:"description"`
}

type IntegrationSystemPage struct {
	Data       []*IntegrationSystem `json:"data"`
	PageInfo   *PageInfo            `json:"pageInfo"`
	TotalCount int                  `json:"totalCount"`
}

func (IntegrationSystemPage) IsPageable() {}

type Label struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type LabelDefinition struct {
	Key    string      `json:"key"`
	Schema *JSONSchema `json:"schema"`
}

type LabelDefinitionInput struct {
	// **Validation:** max=256, alphanumeric chartacters and underscore
	Key    string      `json:"key"`
	Schema *JSONSchema `json:"schema"`
}

type LabelFilter struct {
	// Label key. If query for the filter is not provided, returns every object with given label key regardless of its value.
	Key string `json:"key"`
	// Optional SQL/JSON Path expression. If query is not provided, returns every object with given label key regardless of its value.
	// Currently only a limited subset of expressions is supported.
	Query *string `json:"query"`
}

type LabelInput struct {
	// **Validation:** max=256, alphanumeric chartacters and underscore
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type LabelSelectorInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OAuthCredentialData struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	// URL for getting access token
	URL string `json:"url"`
}

func (OAuthCredentialData) IsCredentialData() {}

type OAuthCredentialDataInput struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	// **Validation:** valid URL
	URL string `json:"url"`
}

type Package struct {
	ID               string      `json:"id"`
	OpenDiscoveryID  string      `json:"openDiscoveryID"`
	ApplicationID    string      `json:"applicationID"`
	Title            string      `json:"title"`
	ShortDescription string      `json:"shortDescription"`
	Description      string      `json:"description"`
	Version          string      `json:"version"`
	Licence          *string     `json:"licence"`
	LicenceType      *string     `json:"licenceType"`
	TermsOfService   *string     `json:"termsOfService"`
	Logo             *string     `json:"logo"`
	Image            *string     `json:"image"`
	Provider         *JSON       `json:"provider"`
	Actions          *JSON       `json:"actions"`
	Tags             *JSON       `json:"tags"`
	LastUpdated      Timestamp   `json:"lastUpdated"`
	Extensions       *JSON       `json:"extensions"`
	Bundles          *BundlePage `json:"bundles"`
	Bundle           *Bundle     `json:"bundle"`
}

type PackageInput struct {
	// TODO: Validation if it is guid
	ID              *string `json:"ID"`
	OpenDiscoveryID string  `json:"openDiscoveryID"`
	// **Validation:** ASCII printable characters, max=100
	Title string `json:"title"`
	// **Validation:** max=2000
	ShortDescription string         `json:"shortDescription"`
	Description      string         `json:"description"`
	Version          string         `json:"version"`
	Licence          *string        `json:"licence"`
	LicenceType      *string        `json:"licenceType"`
	TermsOfService   *string        `json:"termsOfService"`
	Logo             *string        `json:"logo"`
	Image            *string        `json:"image"`
	Provider         *JSON          `json:"provider"`
	Actions          *JSON          `json:"actions"`
	Tags             *JSON          `json:"tags"`
	LastUpdated      Timestamp      `json:"lastUpdated"`
	Extensions       *JSON          `json:"extensions"`
	Bundles          []*BundleInput `json:"bundles"`
}

type PackagePage struct {
	Data       []*Package `json:"data"`
	PageInfo   *PageInfo  `json:"pageInfo"`
	TotalCount int        `json:"totalCount"`
}

func (PackagePage) IsPageable() {}

type PageInfo struct {
	StartCursor PageCursor `json:"startCursor"`
	EndCursor   PageCursor `json:"endCursor"`
	HasNextPage bool       `json:"hasNextPage"`
}

type PlaceholderDefinition struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type PlaceholderDefinitionInput struct {
	// **Validation:**  Up to 36 characters long. Cannot start with a digit. The characters allowed in names are: digits (0-9), lower case letters (a-z),-, and .
	Name string `json:"name"`
	// **Validation:**  max=2000
	Description *string `json:"description"`
}

type RuntimeEventingConfiguration struct {
	DefaultURL string `json:"defaultURL"`
}

type RuntimeInput struct {
	// **Validation:**  Up to 36 characters long. Cannot start with a digit. The characters allowed in names are: digits (0-9), lower case letters (a-z),-, and .
	Name string `json:"name"`
	// **Validation:**  max=2000
	Description *string `json:"description"`
	// **Validation:** key: required, alphanumeric with underscore
	Labels          *Labels                 `json:"labels"`
	StatusCondition *RuntimeStatusCondition `json:"statusCondition"`
}

type RuntimeMetadata struct {
	CreationTimestamp Timestamp `json:"creationTimestamp"`
}

type RuntimePage struct {
	Data       []*Runtime `json:"data"`
	PageInfo   *PageInfo  `json:"pageInfo"`
	TotalCount int        `json:"totalCount"`
}

func (RuntimePage) IsPageable() {}

type RuntimeStatus struct {
	Condition RuntimeStatusCondition `json:"condition"`
	Timestamp Timestamp              `json:"timestamp"`
}

type SystemAuth struct {
	ID   string `json:"id"`
	Auth *Auth  `json:"auth"`
}

type TemplateValueInput struct {
	// **Validation:**  Up to 36 characters long. Cannot start with a digit. The characters allowed in names are: digits (0-9), lower case letters (a-z),-, and .
	Placeholder string `json:"placeholder"`
	Value       string `json:"value"`
}

type Tenant struct {
	ID          string  `json:"id"`
	InternalID  string  `json:"internalID"`
	Name        *string `json:"name"`
	Initialized *bool   `json:"initialized"`
}

type Version struct {
	// for example 4.6
	Value      string `json:"value"`
	Deprecated *bool  `json:"deprecated"`
	// for example 4.5
	DeprecatedSince *string `json:"deprecatedSince"`
	// if true, will be removed in the next version
	ForRemoval *bool `json:"forRemoval"`
}

type VersionInput struct {
	// **Validation:** max=256
	Value      string `json:"value"`
	Deprecated *bool  `json:"deprecated"`
	// **Validation:** max=256
	DeprecatedSince *string `json:"deprecatedSince"`
	ForRemoval      *bool   `json:"forRemoval"`
}

type Viewer struct {
	ID   string     `json:"id"`
	Type ViewerType `json:"type"`
}

type Webhook struct {
	ID            string                 `json:"id"`
	ApplicationID string                 `json:"applicationID"`
	Type          ApplicationWebhookType `json:"type"`
	URL           string                 `json:"url"`
	Auth          *Auth                  `json:"auth"`
}

type WebhookInput struct {
	Type ApplicationWebhookType `json:"type"`
	// **Validation:** valid URL, max=256
	URL  string     `json:"url"`
	Auth *AuthInput `json:"auth"`
}

type APISpecType string

const (
	APISpecTypeOdata   APISpecType = "ODATA"
	APISpecTypeOpenAPI APISpecType = "OPEN_API"
	APISpecTypeCustom  APISpecType = "CUSTOM"
)

var AllAPISpecType = []APISpecType{
	APISpecTypeOdata,
	APISpecTypeOpenAPI,
	APISpecTypeCustom,
}

func (e APISpecType) IsValid() bool {
	switch e {
	case APISpecTypeOdata, APISpecTypeOpenAPI, APISpecTypeCustom:
		return true
	}
	return false
}

func (e APISpecType) String() string {
	return string(e)
}

func (e *APISpecType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = APISpecType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid APISpecType", str)
	}
	return nil
}

func (e APISpecType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ApplicationStatusCondition string

const (
	ApplicationStatusConditionInitial   ApplicationStatusCondition = "INITIAL"
	ApplicationStatusConditionConnected ApplicationStatusCondition = "CONNECTED"
	ApplicationStatusConditionFailed    ApplicationStatusCondition = "FAILED"
)

var AllApplicationStatusCondition = []ApplicationStatusCondition{
	ApplicationStatusConditionInitial,
	ApplicationStatusConditionConnected,
	ApplicationStatusConditionFailed,
}

func (e ApplicationStatusCondition) IsValid() bool {
	switch e {
	case ApplicationStatusConditionInitial, ApplicationStatusConditionConnected, ApplicationStatusConditionFailed:
		return true
	}
	return false
}

func (e ApplicationStatusCondition) String() string {
	return string(e)
}

func (e *ApplicationStatusCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ApplicationStatusCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ApplicationStatusCondition", str)
	}
	return nil
}

func (e ApplicationStatusCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ApplicationTemplateAccessLevel string

const (
	ApplicationTemplateAccessLevelGlobal ApplicationTemplateAccessLevel = "GLOBAL"
)

var AllApplicationTemplateAccessLevel = []ApplicationTemplateAccessLevel{
	ApplicationTemplateAccessLevelGlobal,
}

func (e ApplicationTemplateAccessLevel) IsValid() bool {
	switch e {
	case ApplicationTemplateAccessLevelGlobal:
		return true
	}
	return false
}

func (e ApplicationTemplateAccessLevel) String() string {
	return string(e)
}

func (e *ApplicationTemplateAccessLevel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ApplicationTemplateAccessLevel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ApplicationTemplateAccessLevel", str)
	}
	return nil
}

func (e ApplicationTemplateAccessLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ApplicationWebhookType string

const (
	ApplicationWebhookTypeConfigurationChanged ApplicationWebhookType = "CONFIGURATION_CHANGED"
	ApplicationWebhookTypeOpenDiscovery        ApplicationWebhookType = "OPEN_DISCOVERY"
)

var AllApplicationWebhookType = []ApplicationWebhookType{
	ApplicationWebhookTypeConfigurationChanged,
	ApplicationWebhookTypeOpenDiscovery,
}

func (e ApplicationWebhookType) IsValid() bool {
	switch e {
	case ApplicationWebhookTypeConfigurationChanged, ApplicationWebhookTypeOpenDiscovery:
		return true
	}
	return false
}

func (e ApplicationWebhookType) String() string {
	return string(e)
}

func (e *ApplicationWebhookType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ApplicationWebhookType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ApplicationWebhookType", str)
	}
	return nil
}

func (e ApplicationWebhookType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type BundleInstanceAuthSetStatusConditionInput string

const (
	BundleInstanceAuthSetStatusConditionInputSucceeded BundleInstanceAuthSetStatusConditionInput = "SUCCEEDED"
	BundleInstanceAuthSetStatusConditionInputFailed    BundleInstanceAuthSetStatusConditionInput = "FAILED"
)

var AllBundleInstanceAuthSetStatusConditionInput = []BundleInstanceAuthSetStatusConditionInput{
	BundleInstanceAuthSetStatusConditionInputSucceeded,
	BundleInstanceAuthSetStatusConditionInputFailed,
}

func (e BundleInstanceAuthSetStatusConditionInput) IsValid() bool {
	switch e {
	case BundleInstanceAuthSetStatusConditionInputSucceeded, BundleInstanceAuthSetStatusConditionInputFailed:
		return true
	}
	return false
}

func (e BundleInstanceAuthSetStatusConditionInput) String() string {
	return string(e)
}

func (e *BundleInstanceAuthSetStatusConditionInput) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BundleInstanceAuthSetStatusConditionInput(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BundleInstanceAuthSetStatusConditionInput", str)
	}
	return nil
}

func (e BundleInstanceAuthSetStatusConditionInput) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type BundleInstanceAuthStatusCondition string

const (
	// When creating, before Application sets the credentials
	BundleInstanceAuthStatusConditionPending   BundleInstanceAuthStatusCondition = "PENDING"
	BundleInstanceAuthStatusConditionSucceeded BundleInstanceAuthStatusCondition = "SUCCEEDED"
	BundleInstanceAuthStatusConditionFailed    BundleInstanceAuthStatusCondition = "FAILED"
	// When Runtime requests deletion and Application has to revoke the credentials
	BundleInstanceAuthStatusConditionUnused BundleInstanceAuthStatusCondition = "UNUSED"
)

var AllBundleInstanceAuthStatusCondition = []BundleInstanceAuthStatusCondition{
	BundleInstanceAuthStatusConditionPending,
	BundleInstanceAuthStatusConditionSucceeded,
	BundleInstanceAuthStatusConditionFailed,
	BundleInstanceAuthStatusConditionUnused,
}

func (e BundleInstanceAuthStatusCondition) IsValid() bool {
	switch e {
	case BundleInstanceAuthStatusConditionPending, BundleInstanceAuthStatusConditionSucceeded, BundleInstanceAuthStatusConditionFailed, BundleInstanceAuthStatusConditionUnused:
		return true
	}
	return false
}

func (e BundleInstanceAuthStatusCondition) String() string {
	return string(e)
}

func (e *BundleInstanceAuthStatusCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BundleInstanceAuthStatusCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BundleInstanceAuthStatusCondition", str)
	}
	return nil
}

func (e BundleInstanceAuthStatusCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DocumentFormat string

const (
	DocumentFormatMarkdown DocumentFormat = "MARKDOWN"
)

var AllDocumentFormat = []DocumentFormat{
	DocumentFormatMarkdown,
}

func (e DocumentFormat) IsValid() bool {
	switch e {
	case DocumentFormatMarkdown:
		return true
	}
	return false
}

func (e DocumentFormat) String() string {
	return string(e)
}

func (e *DocumentFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DocumentFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DocumentFormat", str)
	}
	return nil
}

func (e DocumentFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EventSpecType string

const (
	EventSpecTypeAsyncAPI EventSpecType = "ASYNC_API"
	EventSpecTypeCustom   EventSpecType = "CUSTOM"
)

var AllEventSpecType = []EventSpecType{
	EventSpecTypeAsyncAPI,
	EventSpecTypeCustom,
}

func (e EventSpecType) IsValid() bool {
	switch e {
	case EventSpecTypeAsyncAPI, EventSpecTypeCustom:
		return true
	}
	return false
}

func (e EventSpecType) String() string {
	return string(e)
}

func (e *EventSpecType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventSpecType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventSpecType", str)
	}
	return nil
}

func (e EventSpecType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FetchMode string

const (
	FetchModeSingle  FetchMode = "SINGLE"
	FetchModePackage FetchMode = "PACKAGE"
	FetchModeIndex   FetchMode = "INDEX"
)

var AllFetchMode = []FetchMode{
	FetchModeSingle,
	FetchModePackage,
	FetchModeIndex,
}

func (e FetchMode) IsValid() bool {
	switch e {
	case FetchModeSingle, FetchModePackage, FetchModeIndex:
		return true
	}
	return false
}

func (e FetchMode) String() string {
	return string(e)
}

func (e *FetchMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FetchMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FetchMode", str)
	}
	return nil
}

func (e FetchMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FetchRequestStatusCondition string

const (
	FetchRequestStatusConditionInitial   FetchRequestStatusCondition = "INITIAL"
	FetchRequestStatusConditionSucceeded FetchRequestStatusCondition = "SUCCEEDED"
	FetchRequestStatusConditionFailed    FetchRequestStatusCondition = "FAILED"
)

var AllFetchRequestStatusCondition = []FetchRequestStatusCondition{
	FetchRequestStatusConditionInitial,
	FetchRequestStatusConditionSucceeded,
	FetchRequestStatusConditionFailed,
}

func (e FetchRequestStatusCondition) IsValid() bool {
	switch e {
	case FetchRequestStatusConditionInitial, FetchRequestStatusConditionSucceeded, FetchRequestStatusConditionFailed:
		return true
	}
	return false
}

func (e FetchRequestStatusCondition) String() string {
	return string(e)
}

func (e *FetchRequestStatusCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FetchRequestStatusCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FetchRequestStatusCondition", str)
	}
	return nil
}

func (e FetchRequestStatusCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type HealthCheckStatusCondition string

const (
	HealthCheckStatusConditionSucceeded HealthCheckStatusCondition = "SUCCEEDED"
	HealthCheckStatusConditionFailed    HealthCheckStatusCondition = "FAILED"
)

var AllHealthCheckStatusCondition = []HealthCheckStatusCondition{
	HealthCheckStatusConditionSucceeded,
	HealthCheckStatusConditionFailed,
}

func (e HealthCheckStatusCondition) IsValid() bool {
	switch e {
	case HealthCheckStatusConditionSucceeded, HealthCheckStatusConditionFailed:
		return true
	}
	return false
}

func (e HealthCheckStatusCondition) String() string {
	return string(e)
}

func (e *HealthCheckStatusCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HealthCheckStatusCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HealthCheckStatusCondition", str)
	}
	return nil
}

func (e HealthCheckStatusCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type HealthCheckType string

const (
	HealthCheckTypeManagementPlaneApplicationHealthcheck HealthCheckType = "MANAGEMENT_PLANE_APPLICATION_HEALTHCHECK"
)

var AllHealthCheckType = []HealthCheckType{
	HealthCheckTypeManagementPlaneApplicationHealthcheck,
}

func (e HealthCheckType) IsValid() bool {
	switch e {
	case HealthCheckTypeManagementPlaneApplicationHealthcheck:
		return true
	}
	return false
}

func (e HealthCheckType) String() string {
	return string(e)
}

func (e *HealthCheckType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HealthCheckType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HealthCheckType", str)
	}
	return nil
}

func (e HealthCheckType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RuntimeStatusCondition string

const (
	RuntimeStatusConditionInitial      RuntimeStatusCondition = "INITIAL"
	RuntimeStatusConditionProvisioning RuntimeStatusCondition = "PROVISIONING"
	RuntimeStatusConditionConnected    RuntimeStatusCondition = "CONNECTED"
	RuntimeStatusConditionFailed       RuntimeStatusCondition = "FAILED"
)

var AllRuntimeStatusCondition = []RuntimeStatusCondition{
	RuntimeStatusConditionInitial,
	RuntimeStatusConditionProvisioning,
	RuntimeStatusConditionConnected,
	RuntimeStatusConditionFailed,
}

func (e RuntimeStatusCondition) IsValid() bool {
	switch e {
	case RuntimeStatusConditionInitial, RuntimeStatusConditionProvisioning, RuntimeStatusConditionConnected, RuntimeStatusConditionFailed:
		return true
	}
	return false
}

func (e RuntimeStatusCondition) String() string {
	return string(e)
}

func (e *RuntimeStatusCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RuntimeStatusCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RuntimeStatusCondition", str)
	}
	return nil
}

func (e RuntimeStatusCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SpecFormat string

const (
	SpecFormatYaml SpecFormat = "YAML"
	SpecFormatJSON SpecFormat = "JSON"
	SpecFormatXML  SpecFormat = "XML"
)

var AllSpecFormat = []SpecFormat{
	SpecFormatYaml,
	SpecFormatJSON,
	SpecFormatXML,
}

func (e SpecFormat) IsValid() bool {
	switch e {
	case SpecFormatYaml, SpecFormatJSON, SpecFormatXML:
		return true
	}
	return false
}

func (e SpecFormat) String() string {
	return string(e)
}

func (e *SpecFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SpecFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SpecFormat", str)
	}
	return nil
}

func (e SpecFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ViewerType string

const (
	ViewerTypeRuntime           ViewerType = "RUNTIME"
	ViewerTypeApplication       ViewerType = "APPLICATION"
	ViewerTypeIntegrationSystem ViewerType = "INTEGRATION_SYSTEM"
	ViewerTypeUser              ViewerType = "USER"
)

var AllViewerType = []ViewerType{
	ViewerTypeRuntime,
	ViewerTypeApplication,
	ViewerTypeIntegrationSystem,
	ViewerTypeUser,
}

func (e ViewerType) IsValid() bool {
	switch e {
	case ViewerTypeRuntime, ViewerTypeApplication, ViewerTypeIntegrationSystem, ViewerTypeUser:
		return true
	}
	return false
}

func (e ViewerType) String() string {
	return string(e)
}

func (e *ViewerType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ViewerType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ViewerType", str)
	}
	return nil
}

func (e ViewerType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
