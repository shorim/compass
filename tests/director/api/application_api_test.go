package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/tests/director/pkg/jwtbuilder"

	"github.com/kyma-incubator/compass/tests/director/pkg/ptr"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	gcli "github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	eventingCategory            = "eventing"
	registerApplicationCategory = "register application"
	queryApplicationsCategory   = "query applications"
	queryApplicationCategory    = "query application"
	deleteWebhookCategory       = "delete webhook"
	addWebhookCategory          = "add webhook"
	updateWebhookCategory       = "update webhook"
	webhookURL                  = "https://kyma-project.io"
	defaultScenario             = "DEFAULT"
)

var integrationSystemID = "69230297-3c81-4711-aac2-3afa8cb42e2d"

func TestRegisterApplicationWithAllSimpleFieldsProvided(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	in := graphql.ApplicationRegisterInput{
		Name:           "wordpress",
		ProviderName:   ptr.String("provider name"),
		Description:    ptr.String("my first wordpress application"),
		HealthCheckURL: ptr.String("http://mywordpress.com/health"),
		Labels: &graphql.Labels{
			"group":     []interface{}{"production", "experimental"},
			"scenarios": []interface{}{"DEFAULT"},
		},
	}

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)

	actualApp := graphql.ApplicationExt{}

	// WHEN
	request := fixRegisterApplicationRequest(appInputGQL)
	err = tc.RunOperation(ctx, request, &actualApp)

	//THEN
	saveExampleInCustomDir(t, request.Query(), registerApplicationCategory, "register application")
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)
	assertApplication(t, in, actualApp)
	assert.Equal(t, graphql.ApplicationStatusConditionInitial, actualApp.Status.Condition)
}

func TestRegisterApplicationWithStatusCondition(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	statusCond := graphql.ApplicationStatusConditionConnected
	in := graphql.ApplicationRegisterInput{
		Name:           "wordpress",
		ProviderName:   ptr.String("provider name"),
		Description:    ptr.String("my first wordpress application"),
		HealthCheckURL: ptr.String("http://mywordpress.com/health"),
		Labels: &graphql.Labels{
			"group":     []interface{}{"production", "experimental"},
			"scenarios": []interface{}{"DEFAULT"},
		},
		StatusCondition: &statusCond,
	}

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)

	request := fixRegisterApplicationRequest(appInputGQL)

	// WHEN
	actualApp := graphql.ApplicationExt{}
	err = tc.RunOperation(ctx, request, &actualApp)

	//THEN
	saveExampleInCustomDir(t, request.Query(), registerApplicationCategory, "register application with status")
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)
	assertApplication(t, in, actualApp)
	assert.Equal(t, statusCond, actualApp.Status.Condition)
}

func TestRegisterApplicationWithWebhooks(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	in := graphql.ApplicationRegisterInput{
		Name:         "wordpress",
		ProviderName: ptr.String("compass"),
		Webhooks: []*graphql.WebhookInput{
			{
				Type: graphql.ApplicationWebhookTypeConfigurationChanged,
				Auth: fixBasicAuth(t),
				URL:  "http://mywordpress.com/webhooks1",
			},
		},
		Labels: &graphql.Labels{
			"scenarios": []interface{}{"DEFAULT"},
		},
	}

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	actualApp := graphql.ApplicationExt{}

	// WHEN
	request := fixRegisterApplicationRequest(appInputGQL)
	saveExampleInCustomDir(t, request.Query(), registerApplicationCategory, "register application with webhooks")
	err = tc.RunOperation(ctx, request, &actualApp)

	//THEN
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)
	assertApplication(t, in, actualApp)
}

func TestRegisterApplicationWithPackages(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	in := fixApplicationRegisterInputWithPackages(t)
	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	actualApp := graphql.ApplicationExt{}

	// WHEN
	request := fixRegisterApplicationRequest(appInputGQL)
	err = tc.RunOperation(ctx, request, &actualApp)

	//THEN
	saveExampleInCustomDir(t, request.Query(), registerApplicationCategory, "register application with packages")
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)
	assertApplication(t, in, actualApp)
}

func TestCreateApplicationWithNonExistentIntegrationSystem(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	in := fixSampleApplicationCreateInputWithIntegrationSystem("placeholder")
	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	actualApp := graphql.ApplicationExt{}

	request := fixRegisterApplicationRequest(appInputGQL)
	// WHEN
	err = tc.RunOperation(ctx, request, &actualApp)

	//THEN
	require.Error(t, err)
	require.NotNil(t, err.Error())
	require.Contains(t, err.Error(), "Object not found")
}

func TestUpdateApplication(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	actualApp := registerApplication(t, ctx, "before")
	defer unregisterApplication(t, actualApp.ID)

	updateStatusCond := graphql.ApplicationStatusConditionConnected

	expectedApp := actualApp
	expectedApp.Name = "before"
	expectedApp.ProviderName = ptr.String("after")
	expectedApp.Description = ptr.String("after")
	expectedApp.HealthCheckURL = ptr.String(webhookURL)
	expectedApp.Status.Condition = updateStatusCond
	expectedApp.Labels["name"] = "before"

	updateInput := fixSampleApplicationUpdateInput("after")
	updateInput.StatusCondition = &updateStatusCond
	updateInputGQL, err := tc.graphqlizer.ApplicationUpdateInputToGQL(updateInput)
	require.NoError(t, err)
	request := fixUpdateApplicationRequest(actualApp.ID, updateInputGQL)
	updatedApp := graphql.ApplicationExt{}

	//WHEN
	err = tc.RunOperation(ctx, request, &updatedApp)

	//THEN
	require.NoError(t, err)
	assert.Equal(t, expectedApp.ID, updatedApp.ID)
	assert.Equal(t, expectedApp.Name, updatedApp.Name)
	assert.Equal(t, expectedApp.ProviderName, updatedApp.ProviderName)
	assert.Equal(t, expectedApp.Description, updatedApp.Description)
	assert.Equal(t, expectedApp.HealthCheckURL, updatedApp.HealthCheckURL)
	assert.Equal(t, expectedApp.Status.Condition, updatedApp.Status.Condition)

	saveExample(t, request.Query(), "update application")
}

func TestUpdateApplicationWithNonExistentIntegrationSystem(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	actualApp := registerApplication(t, ctx, "before")
	defer unregisterApplication(t, actualApp.ID)

	updateInput := fixSampleApplicationUpdateInputWithIntegrationSystem("after")
	updateInputGQL, err := tc.graphqlizer.ApplicationUpdateInputToGQL(updateInput)
	require.NoError(t, err)
	request := fixUpdateApplicationRequest(actualApp.ID, updateInputGQL)
	updatedApp := graphql.ApplicationExt{}

	//WHEN
	err = tc.RunOperation(ctx, request, &updatedApp)

	//THEN
	require.Error(t, err)
	require.NotNil(t, err.Error())
	require.Contains(t, err.Error(), "Object not found")
}

func TestCreateApplicationWithDuplicatedNamesWithinTenant(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	appName := "samename"

	actualApp := registerApplication(t, ctx, appName)
	defer unregisterApplication(t, actualApp.ID)

	t.Run("Error when creating second Application with same name", func(t *testing.T) {
		in := fixSampleApplicationRegisterInputWithName("first", appName)
		appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
		require.NoError(t, err)
		request := fixRegisterApplicationRequest(appInputGQL)

		// WHEN
		err = tc.RunOperation(ctx, request, nil)

		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not unique")
	})
}

func TestDeleteApplication(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	in := fixSampleApplicationRegisterInput("app")

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	createReq := fixRegisterApplicationRequest(appInputGQL)
	actualApp := graphql.ApplicationExt{}
	err = tc.RunOperation(ctx, createReq, &actualApp)
	require.NoError(t, err)

	require.NotEmpty(t, actualApp.ID)

	// WHEN
	delReq := fixUnregisterApplicationRequest(actualApp.ID)
	saveExample(t, delReq.Query(), "unregister application")
	err = tc.RunOperation(ctx, delReq, &actualApp)

	//THEN
	require.NoError(t, err)
}

func TestUpdateApplicationParts(t *testing.T) {
	ctx := context.Background()
	placeholder := "app"
	in := fixSampleApplicationRegisterInput(placeholder)

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	createReq := fixRegisterApplicationRequest(appInputGQL)
	actualApp := graphql.ApplicationExt{}
	err = tc.RunOperation(ctx, createReq, &actualApp)
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)

	t.Run("labels manipulation", func(t *testing.T) {
		expectedLabel := graphql.Label{Key: "brand_new_label", Value: []interface{}{"aaa", "bbb"}}

		// add label
		createdLabel := &graphql.Label{}

		addReq := fixSetApplicationLabelRequest(actualApp.ID, expectedLabel.Key, []string{"aaa", "bbb"})
		saveExample(t, addReq.Query(), "set application label")
		err := tc.RunOperation(ctx, addReq, &createdLabel)
		require.NoError(t, err)
		assert.Equal(t, &expectedLabel, createdLabel)
		actualApp := getApplication(t, ctx, actualApp.ID)
		assert.Contains(t, actualApp.Labels[expectedLabel.Key], "aaa")
		assert.Contains(t, actualApp.Labels[expectedLabel.Key], "bbb")

		// delete label value
		deletedLabel := graphql.Label{}
		delReq := fixDeleteApplicationLabelRequest(actualApp.ID, expectedLabel.Key)
		saveExample(t, delReq.Query(), "delete application label")
		err = tc.RunOperation(ctx, delReq, &deletedLabel)
		require.NoError(t, err)
		assert.Equal(t, expectedLabel, deletedLabel)
		actualApp = getApplication(t, ctx, actualApp.ID)
		assert.Nil(t, actualApp.Labels[expectedLabel.Key])

	})

	t.Run("manage webhooks", func(t *testing.T) {
		// add
		webhookInStr, err := tc.graphqlizer.WebhookInputToGQL(&graphql.WebhookInput{
			URL:  "http://new-webhook.url",
			Type: graphql.ApplicationWebhookTypeConfigurationChanged,
		})

		require.NoError(t, err)
		addReq := fixAddWebhookRequest(actualApp.ID, webhookInStr)
		saveExampleInCustomDir(t, addReq.Query(), addWebhookCategory, "add application webhook")

		actualWebhook := graphql.Webhook{}
		err = tc.RunOperation(ctx, addReq, &actualWebhook)
		require.NoError(t, err)
		assert.Equal(t, "http://new-webhook.url", actualWebhook.URL)
		assert.Equal(t, graphql.ApplicationWebhookTypeConfigurationChanged, actualWebhook.Type)
		id := actualWebhook.ID
		require.NotNil(t, id)

		// get all webhooks
		updatedApp := getApplication(t, ctx, actualApp.ID)
		assert.Len(t, updatedApp.Webhooks, 2)

		// update
		webhookInStr, err = tc.graphqlizer.WebhookInputToGQL(&graphql.WebhookInput{
			URL: "http://updated-webhook.url", Type: graphql.ApplicationWebhookTypeConfigurationChanged,
		})

		require.NoError(t, err)
		updateReq := fixUpdateWebhookRequest(actualWebhook.ID, webhookInStr)
		saveExampleInCustomDir(t, updateReq.Query(), updateWebhookCategory, "update application webhook")
		err = tc.RunOperation(ctx, updateReq, &actualWebhook)
		require.NoError(t, err)
		assert.Equal(t, "http://updated-webhook.url", actualWebhook.URL)

		// delete

		//GIVEN
		deleteReq := fixDeleteWebhookRequest(actualWebhook.ID)
		saveExampleInCustomDir(t, deleteReq.Query(), deleteWebhookCategory, "delete application webhook")

		//WHEN
		err = tc.RunOperation(ctx, deleteReq, &actualWebhook)

		//THEN
		require.NoError(t, err)
		assert.Equal(t, "http://updated-webhook.url", actualWebhook.URL)

	})

	t.Run("refetch API", func(t *testing.T) {
		// TODO later
	})

	t.Run("refetch Event Spec", func(t *testing.T) {
		// TODO later
	})
}

func TestQueryApplications(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		in := graphql.ApplicationRegisterInput{
			Name: fmt.Sprintf("app-%d", i),
		}

		appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
		require.NoError(t, err)
		actualApp := graphql.Application{}
		request := fixRegisterApplicationRequest(appInputGQL)
		err = tc.RunOperation(ctx, request, &actualApp)
		require.NoError(t, err)
		defer unregisterApplication(t, actualApp.ID)
	}
	actualAppPage := graphql.ApplicationPage{}

	// WHEN
	queryReq := fixApplicationsRequest()
	err := tc.RunOperation(ctx, queryReq, &actualAppPage)
	saveExampleInCustomDir(t, queryReq.Query(), queryApplicationsCategory, "query applications")

	//THEN
	require.NoError(t, err)
	assert.Len(t, actualAppPage.Data, 3)
	assert.Equal(t, 3, actualAppPage.TotalCount)
}

func TestQueryApplicationsPageable(t *testing.T) {
	// GIVEN
	appAmount := 7
	after := 3
	cursor := ""
	ctx := context.Background()

	apps := make(map[string]*graphql.ApplicationExt)
	for i := 0; i < appAmount; i++ {
		app := registerApplication(t, ctx, fmt.Sprintf("app-%d", i))
		defer unregisterApplication(t, app.ID)
		apps[app.ID] = &app
	}
	appsPage := graphql.ApplicationPageExt{}

	// WHEN
	queriesForFullPage := appAmount / after
	for i := 0; i < queriesForFullPage; i++ {
		appReq := fixApplicationsPageableRequest(after, cursor)
		err := tc.RunOperation(ctx, appReq, &appsPage)
		require.NoError(t, err)

		//THEN
		assert.Equal(t, cursor, string(appsPage.PageInfo.StartCursor))
		assert.True(t, appsPage.PageInfo.HasNextPage)
		assert.Len(t, appsPage.Data, after)
		assert.Equal(t, appAmount, appsPage.TotalCount)
		for _, app := range appsPage.Data {
			assert.Equal(t, app, apps[app.ID])
			delete(apps, app.ID)
		}
		cursor = string(appsPage.PageInfo.EndCursor)
	}

	appReq := fixApplicationsPageableRequest(after, cursor)
	err := tc.RunOperation(ctx, appReq, &appsPage)
	require.NoError(t, err)

	assert.False(t, appsPage.PageInfo.HasNextPage)
	assert.Empty(t, appsPage.PageInfo.EndCursor)
	assert.Equal(t, appAmount, appsPage.TotalCount)
	require.Len(t, appsPage.Data, 1)
	delete(apps, appsPage.Data[0].ID)
	assert.Len(t, apps, 0)
}

func TestQuerySpecificApplication(t *testing.T) {
	// GIVEN
	in := graphql.ApplicationRegisterInput{
		Name:         "app",
		ProviderName: ptr.String("compass"),
	}

	appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)

	actualApp := graphql.Application{}
	request := fixRegisterApplicationRequest(appInputGQL)
	err = tc.RunOperation(context.Background(), request, &actualApp)
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	appID := actualApp.ID
	defer unregisterApplication(t, appID)

	t.Run("Query Application With Consumer User", func(t *testing.T) {
		actualApp := graphql.Application{}

		// WHEN
		queryAppReq := fixApplicationRequest(appID)
		err = tc.RunOperation(context.Background(), queryAppReq, &actualApp)
		saveExampleInCustomDir(t, queryAppReq.Query(), queryApplicationCategory, "query application")

		//THE
		require.NoError(t, err)
		assert.Equal(t, appID, actualApp.ID)
	})

	ctx := context.Background()

	runtime := registerRuntime(t, ctx, "runtime-test")
	defer unregisterRuntime(t, runtime.ID)

	scenarios := []string{defaultScenario, "test-scenario"}

	// update label definitions
	updateScenariosLabelDefinitionWithinTenant(t, ctx, testTenants.GetDefaultTenantID(), scenarios)
	defer updateScenariosLabelDefinitionWithinTenant(t, ctx, testTenants.GetDefaultTenantID(), scenarios[:1])

	runtimeConsumer := tc.NewOperation(ctx).WithConsumer(&jwtbuilder.Consumer{
		ID:   runtime.ID,
		Type: jwtbuilder.RuntimeConsumer,
	})

	t.Run("Query Application With Consumer Runtime in same scenario", func(t *testing.T) {
		// set application scenarios label
		setApplicationLabel(t, ctx, appID, scenariosLabel, scenarios[1:])
		defer setApplicationLabel(t, ctx, appID, scenariosLabel, scenarios[:1])

		// set runtime scenarios label
		setRuntimeLabel(t, ctx, runtime.ID, scenariosLabel, scenarios[1:])
		defer setRuntimeLabel(t, ctx, runtime.ID, scenariosLabel, scenarios[:1])

		actualApp := graphql.Application{}

		// WHEN
		queryAppReq := fixApplicationRequest(appID)
		err = runtimeConsumer.Run(queryAppReq, &actualApp)

		//THEN
		require.NoError(t, err)
		assert.Equal(t, appID, actualApp.ID)
	})

	t.Run("Query Application With Consumer Runtime not in same scenario", func(t *testing.T) {
		// set application scenarios label
		setApplicationLabel(t, ctx, appID, scenariosLabel, scenarios[:1])

		// set runtime scenarios label
		setRuntimeLabel(t, ctx, runtime.ID, scenariosLabel, scenarios[1:])
		defer setRuntimeLabel(t, ctx, runtime.ID, scenariosLabel, scenarios[:1])

		actualApp := graphql.Application{}

		// WHEN
		queryAppReq := fixApplicationRequest(appID)
		err = runtimeConsumer.Run(queryAppReq, &actualApp)

		//THEN
		require.Error(t, err)
		require.Contains(t, err.Error(), "The operation is not allowed")
	})
}

func TestTenantSeparation(t *testing.T) {
	// GIVEN
	appIn := fixSampleApplicationRegisterInput("tenantseparation")
	inStr, err := tc.graphqlizer.ApplicationRegisterInputToGQL(appIn)
	require.NoError(t, err)
	createReq := fixRegisterApplicationRequest(inStr)
	actualApp := graphql.ApplicationExt{}
	ctx := context.Background()
	err = tc.RunOperation(ctx, createReq, &actualApp)
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)
	defer unregisterApplication(t, actualApp.ID)

	// WHEN
	getAppReq := fixApplicationsRequest()
	customTenant := testTenants.GetIDByName(t, "Test1")
	anotherTenantsApps := graphql.ApplicationPage{}
	// THEN
	err = tc.RunOperationWithCustomTenant(ctx, customTenant, getAppReq, &anotherTenantsApps)
	require.NoError(t, err)
	assert.Empty(t, anotherTenantsApps.Data)
}

func TestApplicationsForRuntime(t *testing.T) {
	//GIVEN
	ctx := context.Background()
	tenantID := testTenants.GetIDByName(t, "Test1")
	otherTenant := testTenants.GetIDByName(t, "Test2")
	tenantApplications := []*graphql.Application{}
	defaultValue := "DEFAULT"
	scenarios := []string{defaultValue, "black-friday-campaign", "christmas-campaign", "summer-campaign"}

	jsonSchema := map[string]interface{}{
		"type":        "array",
		"minItems":    1,
		"uniqueItems": true,
		"items": map[string]interface{}{
			"type": "string",
			"enum": scenarios,
		},
	}
	var schema interface{} = jsonSchema

	createLabelDefinitionWithinTenant(t, ctx, scenariosLabel, schema, tenantID)
	createLabelDefinitionWithinTenant(t, ctx, scenariosLabel, schema, otherTenant)

	applications := []struct {
		ApplicationName string
		Tenant          string
		WithinTenant    bool
		Scenarios       []string
	}{
		{
			Tenant:          tenantID,
			ApplicationName: "first",
			WithinTenant:    true,
			Scenarios:       []string{defaultValue},
		},
		{
			Tenant:          tenantID,
			ApplicationName: "second",
			WithinTenant:    true,
			Scenarios:       []string{defaultValue, "black-friday-campaign"},
		},
		{
			Tenant:          tenantID,
			ApplicationName: "third",
			WithinTenant:    true,
			Scenarios:       []string{"black-friday-campaign", "christmas-campaign", "summer-campaign"},
		},
		{
			Tenant:          tenantID,
			ApplicationName: "allscenarios",
			WithinTenant:    true,
			Scenarios:       []string{defaultValue, "black-friday-campaign", "christmas-campaign", "summer-campaign"},
		},
		{
			Tenant:          otherTenant,
			ApplicationName: "test",
			WithinTenant:    false,
			Scenarios:       []string{defaultValue, "black-friday-campaign"},
		},
	}

	for _, testApp := range applications {
		applicationInput := fixSampleApplicationRegisterInput(testApp.ApplicationName)
		applicationInput.Labels = &graphql.Labels{scenariosLabel: testApp.Scenarios}
		appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(applicationInput)
		require.NoError(t, err)

		createApplicationReq := fixRegisterApplicationRequest(appInputGQL)
		application := graphql.Application{}

		err = tc.RunOperationWithCustomTenant(ctx, testApp.Tenant, createApplicationReq, &application)

		require.NoError(t, err)
		require.NotEmpty(t, application.ID)

		defer unregisterApplicationInTenant(t, application.ID, testApp.Tenant)
		if testApp.WithinTenant {
			tenantApplications = append(tenantApplications, &application)
		}
	}

	//create runtime
	runtimeInput := fixRuntimeInput("runtime")
	(*runtimeInput.Labels)[scenariosLabel] = scenarios
	runtimeInputGQL, err := tc.graphqlizer.RuntimeInputToGQL(runtimeInput)
	require.NoError(t, err)
	registerRuntimeRequest := fixRegisterRuntimeRequest(runtimeInputGQL)
	runtime := graphql.Runtime{}
	err = tc.RunOperationWithCustomTenant(ctx, tenantID, registerRuntimeRequest, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)
	defer unregisterRuntimeWithinTenant(t, runtime.ID, tenantID)

	t.Run("Applications For Runtime Query", func(t *testing.T) {
		request := fixApplicationForRuntimeRequest(runtime.ID)
		applicationPage := graphql.ApplicationPage{}

		err = tc.RunOperationWithCustomTenant(ctx, tenantID, request, &applicationPage)
		saveExample(t, request.Query(), "query applications for runtime")

		//THEN
		require.NoError(t, err)
		require.Len(t, applicationPage.Data, len(tenantApplications))
		assert.ElementsMatch(t, tenantApplications, applicationPage.Data)

	})

	t.Run("Applications Query With Consumer Runtime", func(t *testing.T) {
		request := fixApplicationsRequest()
		applicationPage := graphql.ApplicationPage{}

		err = tc.NewOperation(ctx).WithTenant(tenantID).WithConsumer(&jwtbuilder.Consumer{
			ID:   runtime.ID,
			Type: jwtbuilder.RuntimeConsumer,
		}).Run(request, &applicationPage)

		//THEN
		require.NoError(t, err)
		require.Len(t, applicationPage.Data, len(tenantApplications))
		assert.ElementsMatch(t, tenantApplications, applicationPage.Data)

	})
}

func TestApplicationsForRuntimeWithHiddenApps(t *testing.T) {
	//GIVEN
	ctx := context.Background()
	tenantID := testTenants.GetIDByName(t, "TestApplicationsForRuntimeWithHiddenApps")
	expectedApplications := []*graphql.Application{}

	defaultValue := "DEFAULT"
	scenarios := []string{defaultValue, "test-scenario"}

	jsonSchema := map[string]interface{}{
		"type":        "array",
		"minItems":    1,
		"uniqueItems": true,
		"items": map[string]interface{}{
			"type": "string",
			"enum": scenarios,
		},
	}
	var schema interface{} = jsonSchema

	createLabelDefinitionWithinTenant(t, ctx, scenariosLabel, schema, tenantID)

	applications := []struct {
		ApplicationName string
		Scenarios       []string
		Hidden          bool
	}{
		{
			ApplicationName: "first",
			Scenarios:       []string{defaultValue},
			Hidden:          false,
		},
		{
			ApplicationName: "second",
			Scenarios:       []string{"test-scenario"},
			Hidden:          false,
		},
		{
			ApplicationName: "third",
			Scenarios:       []string{"test-scenario"},
			Hidden:          true,
		},
	}

	applicationHideSelectorKey := "applicationHideSelectorKey"
	applicationHideSelectorValue := "applicationHideSelectorValue"

	for _, testApp := range applications {
		applicationInput := fixSampleApplicationRegisterInput(testApp.ApplicationName)
		applicationInput.Labels = &graphql.Labels{scenariosLabel: testApp.Scenarios}
		if testApp.Hidden {
			(*applicationInput.Labels)[applicationHideSelectorKey] = applicationHideSelectorValue
		}
		appInputGQL, err := tc.graphqlizer.ApplicationRegisterInputToGQL(applicationInput)
		require.NoError(t, err)

		createApplicationReq := fixRegisterApplicationRequest(appInputGQL)
		application := graphql.Application{}

		err = tc.RunOperationWithCustomTenant(ctx, tenantID, createApplicationReq, &application)

		require.NoError(t, err)
		require.NotEmpty(t, application.ID)

		defer unregisterApplicationInTenant(t, application.ID, tenantID)
		if !testApp.Hidden {
			expectedApplications = append(expectedApplications, &application)
		}
	}

	//create runtime
	runtimeInput := fixRuntimeInput("runtime")
	(*runtimeInput.Labels)[scenariosLabel] = scenarios
	runtimeInputGQL, err := tc.graphqlizer.RuntimeInputToGQL(runtimeInput)
	require.NoError(t, err)
	registerRuntimeRequest := fixRegisterRuntimeRequest(runtimeInputGQL)
	runtime := graphql.Runtime{}
	err = tc.RunOperationWithCustomTenant(ctx, tenantID, registerRuntimeRequest, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)
	defer unregisterRuntimeWithinTenant(t, runtime.ID, tenantID)

	t.Run("Applications For Runtime Query", func(t *testing.T) {
		//WHEN
		request := fixApplicationForRuntimeRequest(runtime.ID)
		applicationPage := graphql.ApplicationPage{}

		err = tc.RunOperationWithCustomTenant(ctx, tenantID, request, &applicationPage)

		//THEN
		require.NoError(t, err)
		require.Len(t, applicationPage.Data, len(expectedApplications))
		assert.ElementsMatch(t, expectedApplications, applicationPage.Data)
	})

	t.Run("Applications Query With Consumer Runtime", func(t *testing.T) {
		//WHEN
		request := fixApplicationsRequest()
		applicationPage := graphql.ApplicationPage{}

		err = tc.NewOperation(ctx).WithTenant(tenantID).WithConsumer(&jwtbuilder.Consumer{
			ID:   runtime.ID,
			Type: jwtbuilder.RuntimeConsumer,
		}).Run(request, &applicationPage)

		//THEN
		require.NoError(t, err)
		require.Len(t, applicationPage.Data, len(expectedApplications))
		assert.ElementsMatch(t, expectedApplications, applicationPage.Data)
	})
}

func fixSampleApplicationRegisterInput(placeholder string) graphql.ApplicationRegisterInput {
	return graphql.ApplicationRegisterInput{
		Name:         placeholder,
		ProviderName: ptr.String("compass"),
		Webhooks: []*graphql.WebhookInput{{
			Type: graphql.ApplicationWebhookTypeConfigurationChanged,
			URL:  webhookURL},
		},
	}
}

func fixSampleApplicationRegisterInputWithName(placeholder, name string) graphql.ApplicationRegisterInput {
	sampleInput := fixSampleApplicationRegisterInput(placeholder)
	sampleInput.Name = name
	return sampleInput
}

func fixSampleApplicationCreateInputWithIntegrationSystem(placeholder string) graphql.ApplicationRegisterInput {
	sampleInput := fixSampleApplicationRegisterInput(placeholder)
	sampleInput.IntegrationSystemID = &integrationSystemID
	return sampleInput
}

func fixSampleApplicationUpdateInput(placeholder string) graphql.ApplicationUpdateInput {
	return graphql.ApplicationUpdateInput{
		Description:    &placeholder,
		HealthCheckURL: ptr.String(webhookURL),
		ProviderName:   &placeholder,
	}
}

func fixSampleApplicationUpdateInputWithIntegrationSystem(placeholder string) graphql.ApplicationUpdateInput {
	return graphql.ApplicationUpdateInput{
		Description:         &placeholder,
		HealthCheckURL:      ptr.String(webhookURL),
		IntegrationSystemID: &integrationSystemID,
		ProviderName:        ptr.String(placeholder),
	}
}

func unregisterApplicationInTenant(t *testing.T, id string, tenant string) {
	req := gcli.NewRequest(
		fmt.Sprintf(`mutation {
		unregisterApplication(id: "%s") {
			id
		}	
	}`, id))
	require.NoError(t, tc.RunOperationWithCustomTenant(context.Background(), tenant, req, nil))
}

func unregisterApplication(t *testing.T, id string) {
	req := gcli.NewRequest(
		fmt.Sprintf(`mutation {
		unregisterApplication(id: "%s") {
			id
		}	
	}`, id))
	require.NoError(t, tc.RunOperation(context.Background(), req, nil))
}

func deleteAPI(t *testing.T, id string) {
	req := gcli.NewRequest(
		fmt.Sprintf(`mutation {
		deleteAPIDefinition(id: "%s") {
			id
		}	
	}`, id))
	require.NoError(t, tc.RunOperation(context.Background(), req, nil))
}

func deleteEventAPI(t *testing.T, id string) {
	req := gcli.NewRequest(
		fmt.Sprintf(`mutation {
		deleteEventDefinition(id: "%s") {
			id
		}	
	}`, id))
	require.NoError(t, tc.RunOperation(context.Background(), req, nil))
}
