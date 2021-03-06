package scenario_test

import (
	"context"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/99designs/gqlgen/graphql"
	lbl_mock "github.com/kyma-incubator/compass/components/director/internal/domain/label/automock"
	pkg_mock "github.com/kyma-incubator/compass/components/director/internal/domain/package/automock"
	pkg_auth_mock "github.com/kyma-incubator/compass/components/director/internal/domain/packageinstanceauth/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/consumer"

	"github.com/kyma-incubator/compass/components/director/pkg/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasScenario(t *testing.T) {
	t.Run("could not extract consumer information, should return error", func(t *testing.T) {
		// GIVEN
		directive := scenario.NewDirective(nil, nil, nil, nil)
		// WHEN
		res, err := directive.HasScenario(context.TODO(), nil, nil, "", "")
		// THEN
		require.Error(t, err)
		assert.EqualError(t, err, consumer.NoConsumerError.Error())
		assert.Equal(t, res, nil)
	})

	t.Run("consumer is of type user, should proceed with next resolver", func(t *testing.T) {
		// GIVEN
		directive := scenario.NewDirective(nil, nil, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.User})
		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, "", "")
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})

	t.Run("consumer is of type application, should proceed with next resolver", func(t *testing.T) {
		// GIVEN
		directive := scenario.NewDirective(nil, nil, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.Application})
		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, "", "")
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})

	t.Run("consumer is of type integration system, should proceed with next resolver", func(t *testing.T) {
		// GIVEN
		directive := scenario.NewDirective(nil, nil, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.IntegrationSystem})
		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, "", "")
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})

	t.Run("could not extract tenant from context, should return error", func(t *testing.T) {
		// GIVEN
		directive := scenario.NewDirective(nil, nil, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.Runtime})
		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, "", "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), apperrors.NewCannotReadTenantError().Error())
		assert.Equal(t, res, nil)
	})

	t.Run("runtime requests non-existent application", func(t *testing.T) {
		// GIVEN
		const (
			idField       = "id"
			tenantID      = "42"
			applicationID = "24"
		)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatDoesntExpectCommit()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "Application",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{idField: applicationID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		notFoundErr := apperrors.NewNotFoundError(resource.Label, model.ScenariosKey)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, applicationID, model.ScenariosKey).Return(nil, notFoundErr)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationID, idField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, notFoundErr)
		assert.Equal(t, res, nil)
	})

	t.Run("runtime requests package instance auth creation for non-existent package", func(t *testing.T) {
		// GIVEN
		const (
			packageIDField = "packageID"
			tenantID       = "42"
			packageID      = "24"
		)

		pkgRepo := &pkg_mock.PackageRepository{}
		defer pkgRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatDoesntExpectCommit()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, nil, pkgRepo, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{packageIDField: packageID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		notFoundErr := apperrors.NewNotFoundErrorWithType(resource.Package)
		pkgRepo.On("GetByID", ctxWithTx, tenantID, packageID).Return(nil, notFoundErr)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationIDByPackage, packageIDField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, notFoundErr)
		assert.Equal(t, res, nil)
	})

	t.Run("runtime requests package instance auth deletion for non-existent system auth ID", func(t *testing.T) {
		// GIVEN
		const (
			pkgAuthIDField = "authID"
			tenantID       = "42"
			pkgAuthID      = "24"
		)

		pkgAuthRepo := &pkg_auth_mock.Repository{}
		defer pkgAuthRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatDoesntExpectCommit()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, nil, nil, pkgAuthRepo)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{pkgAuthIDField: pkgAuthID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		notFoundErr := apperrors.NewNotFoundErrorWithType(resource.PackageInstanceAuth)
		pkgAuthRepo.On("GetByID", ctxWithTx, tenantID, pkgAuthID).Return(nil, notFoundErr)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationIDByPackageInstanceAuth, pkgAuthIDField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, notFoundErr)
		assert.Equal(t, res, nil)
	})

	t.Run("runtime is in formation with application in application query", func(t *testing.T) {
		// GIVEN
		const (
			idField       = "id"
			tenantID      = "42"
			runtimeID     = "23"
			applicationID = "24"
		)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "Application",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{idField: applicationID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, applicationID, model.ScenariosKey).Return(mockedLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedLabel, nil)

		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, scenario.GetApplicationID, idField)
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})

	t.Run("runtime is NOT in formation with application in application query", func(t *testing.T) {
		// GIVEN
		const (
			idField       = "id"
			tenantID      = "42"
			runtimeID     = "23"
			applicationID = "24"
		)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, nil, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "Application",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{idField: applicationID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedAppLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		mockedRuntimeLabel := &model.Label{Value: []interface{}{"TEST"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, applicationID, model.ScenariosKey).Return(mockedAppLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedRuntimeLabel, nil)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationID, idField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, scenario.ErrMissingScenario)
		assert.Equal(t, res, nil)
	})

	t.Run("runtime is in formation with owning application in request package instance auth flow ", func(t *testing.T) {
		// GIVEN
		const (
			packageIDField = "packageID"
			tenantID       = "42"
			packageID      = "24"
			runtimeID      = "23"
			applicationID  = "22"
		)

		pkgRepo := &pkg_mock.PackageRepository{}
		defer pkgRepo.AssertExpectations(t)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, pkgRepo, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{packageIDField: packageID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedPkg := &model.Package{ApplicationID: applicationID}
		pkgRepo.On("GetByID", ctxWithTx, tenantID, packageID).Return(mockedPkg, nil)

		mockedLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, mockedPkg.ApplicationID, model.ScenariosKey).Return(mockedLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedLabel, nil)

		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, scenario.GetApplicationIDByPackage, packageIDField)
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})
	t.Run("runtime is NOT in formation with owning application in request package instance auth flow ", func(t *testing.T) {
		// GIVEN
		const (
			packageIDField = "packageID"
			tenantID       = "42"
			packageID      = "24"
			runtimeID      = "23"
			applicationID  = "22"
		)

		pkgRepo := &pkg_mock.PackageRepository{}
		defer pkgRepo.AssertExpectations(t)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, pkgRepo, nil)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{packageIDField: packageID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedPkg := &model.Package{ApplicationID: applicationID}
		pkgRepo.On("GetByID", ctxWithTx, tenantID, packageID).Return(mockedPkg, nil)

		mockedAppLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		mockedRuntimeLabel := &model.Label{Value: []interface{}{"TEST"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, applicationID, model.ScenariosKey).Return(mockedAppLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedRuntimeLabel, nil)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationIDByPackage, packageIDField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, scenario.ErrMissingScenario)
		assert.Equal(t, res, nil)
	})

	t.Run("runtime is in formation with owning application in delete package instance auth flow", func(t *testing.T) {
		// GIVEN
		const (
			pkgAuthIDField = "authID"
			tenantID       = "42"
			pkgAuthID      = "24"
			runtimeID      = "23"
			applicationID  = "22"
			packageID      = "21"
		)

		pkgAuthRepo := &pkg_auth_mock.Repository{}
		defer pkgAuthRepo.AssertExpectations(t)

		pkgRepo := &pkg_mock.PackageRepository{}
		defer pkgRepo.AssertExpectations(t)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, pkgRepo, pkgAuthRepo)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{pkgAuthIDField: pkgAuthID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedPkgAuth := &model.PackageInstanceAuth{PackageID: packageID}
		pkgAuthRepo.On("GetByID", ctxWithTx, tenantID, pkgAuthID).Return(mockedPkgAuth, nil)

		mockedPkg := &model.Package{ApplicationID: applicationID}
		pkgRepo.On("GetByID", ctxWithTx, tenantID, mockedPkgAuth.PackageID).Return(mockedPkg, nil)

		mockedLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, mockedPkg.ApplicationID, model.ScenariosKey).Return(mockedLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedLabel, nil)

		dummyResolver := &dummyResolver{}
		// WHEN
		res, err := directive.HasScenario(ctx, nil, dummyResolver.SuccessResolve, scenario.GetApplicationIDByPackageInstanceAuth, pkgAuthIDField)
		// THEN
		require.NoError(t, err)
		assert.Equal(t, res, mockedNextOutput())
	})
	t.Run("runtime is NOT in formation with owning application in delete package instance auth flow", func(t *testing.T) {
		// GIVEN
		const (
			pkgAuthIDField = "authID"
			tenantID       = "42"
			pkgAuthID      = "24"
			runtimeID      = "23"
			applicationID  = "22"
			packageID      = "21"
		)

		pkgAuthRepo := &pkg_auth_mock.Repository{}
		defer pkgAuthRepo.AssertExpectations(t)

		pkgRepo := &pkg_mock.PackageRepository{}
		defer pkgRepo.AssertExpectations(t)

		lblRepo := &lbl_mock.LabelRepository{}
		defer lblRepo.AssertExpectations(t)

		mockedTx, mockedTransactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceeds()
		defer mockedTx.AssertExpectations(t)
		defer mockedTransactioner.AssertExpectations(t)

		directive := scenario.NewDirective(mockedTransactioner, lblRepo, pkgRepo, pkgAuthRepo)
		ctx := context.WithValue(context.TODO(), consumer.ConsumerKey, consumer.Consumer{ConsumerID: runtimeID, ConsumerType: consumer.Runtime})
		ctx = context.WithValue(ctx, tenant.TenantContextKey, tenant.TenantCtx{InternalID: tenantID})
		rCtx := &graphql.ResolverContext{
			Object:   "PackageInstanceAuth",
			Field:    graphql.CollectedField{},
			Args:     map[string]interface{}{pkgAuthIDField: pkgAuthID},
			IsMethod: false,
		}
		ctx = graphql.WithResolverContext(ctx, rCtx)
		ctxWithTx := persistence.SaveToContext(ctx, mockedTx)

		mockedPkgAuth := &model.PackageInstanceAuth{PackageID: packageID}
		pkgAuthRepo.On("GetByID", ctxWithTx, tenantID, pkgAuthID).Return(mockedPkgAuth, nil)

		mockedPkg := &model.Package{ApplicationID: applicationID}
		pkgRepo.On("GetByID", ctxWithTx, tenantID, mockedPkgAuth.PackageID).Return(mockedPkg, nil)

		mockedAppLabel := &model.Label{Value: []interface{}{"DEFAULT"}}
		mockedRuntimeLabel := &model.Label{Value: []interface{}{"TEST"}}
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.ApplicationLabelableObject, mockedPkg.ApplicationID, model.ScenariosKey).Return(mockedAppLabel, nil)
		lblRepo.On("GetByKey", ctxWithTx, tenantID, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey).Return(mockedRuntimeLabel, nil)
		// WHEN
		res, err := directive.HasScenario(ctx, nil, nil, scenario.GetApplicationIDByPackageInstanceAuth, pkgAuthIDField)
		// THEN
		require.Error(t, err)
		assert.Error(t, err, scenario.ErrMissingScenario)
		assert.Equal(t, res, nil)
	})

}

type dummyResolver struct {
	called bool
}

func (d *dummyResolver) SuccessResolve(_ context.Context) (res interface{}, err error) {
	d.called = true
	return mockedNextOutput(), nil
}

func mockedNextOutput() string {
	return "nextOutput"
}
