package tenantfetcher

import (
	"context"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strconv"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
)

type TenantFieldMapping struct {
	TotalPagesField    string `envconfig:"APP_TENANT_TOTAL_PAGES_FIELD"`
	TotalResultsField  string `envconfig:"APP_TENANT_TOTAL_RESULTS_FIELD"`
	EventsField        string `envconfig:"APP_TENANT_EVENTS_FIELD"`
	TimestampField     string `envconfig:"APP_TENANT_TIMESTAMP_FIELD"`
	NameField          string `envconfig:"default=name,APP_MAPPING_FIELD_NAME"`
	IDField            string `envconfig:"default=id,APP_MAPPING_FIELD_ID"`
	DetailsField       string `envconfig:"default=details,APP_MAPPING_FIELD_DETAILS"`
	DiscriminatorField string `envconfig:"optional,APP_MAPPING_FIELD_DISCRIMINATOR"`
	DiscriminatorValue string `envconfig:"optional,APP_MAPPING_VALUE_DISCRIMINATOR"`
}

// QueryConfig contains the name of query parameters fields and default/start values
type QueryConfig struct {
	PageNumField   string `envconfig:"default=pageNum,APP_QUERY_PAGE_NUM_FIELD"`
	PageSizeField  string `envconfig:"default=pageSize,APP_QUERY_PAGE_SIZE_FIELD"`
	TimestampField string `envconfig:"default=timestamp,APP_QUERY_TIMESTAMP_FIELD"`
	PageStartValue string `envconfig:"default=0,APP_QUERY_PAGE_START"`
	PageSizeValue  string `envconfig:"default=150,APP_QUERY_PAGE_SIZE"`
	TimestampValue string `envconfig:"default=1"`
}

type KubernetesClientConfig struct {
	PollInterval  time.Duration `envconfig:"default=2s"`
	PollTimeout   time.Duration `envconfig:"default=1m"`
	Timeout       time.Duration `envconfig:"default=95s"`
	Namespace     string        `envconfig:"default=compass-system"`
	ConfigMapName string
}

//go:generate mockery -name=TenantStorageService -output=automock -outpkg=automock -case=underscore
type TenantStorageService interface {
	List(ctx context.Context) ([]*model.BusinessTenantMapping, error)
	CreateManyIfNotExists(ctx context.Context, tenantInputs []model.BusinessTenantMappingInput) error
	DeleteMany(ctx context.Context, tenantInputs []model.BusinessTenantMappingInput) error
}

//go:generate mockery -name=EventAPIClient -output=automock -outpkg=automock -case=underscore
type EventAPIClient interface {
	FetchTenantEventsPage(eventsType EventsType, additionalQueryParams QueryParams) (TenantEventsResponse, error)
}

const (
	retryAttempts          = 7
	retryDelayMilliseconds = 100
)

type Service struct {
	queryConfig            QueryConfig
	kubernetesClientConfig KubernetesClientConfig
	transact               persistence.Transactioner
	eventAPIClient         EventAPIClient
	tenantStorageService   TenantStorageService
	providerName           string
	fieldMapping           TenantFieldMapping

	retryAttempts uint
}

func NewService(queryConfig QueryConfig, transact persistence.Transactioner, fieldMapping TenantFieldMapping, providerName string, client EventAPIClient, tenantStorageService TenantStorageService) *Service {
	return &Service{
		transact:             transact,
		fieldMapping:         fieldMapping,
		providerName:         providerName,
		eventAPIClient:       client,
		tenantStorageService: tenantStorageService,
		queryConfig:          queryConfig,

		retryAttempts: retryAttempts,
	}
}

func (s Service) SyncTenants() error {
	ctx := context.Background()
	// create k8s client instance
	k8sClientSet, err := newK8SClientSet(ctx, s.kubernetesClientConfig.PollInterval, s.kubernetesClientConfig.PollTimeout, s.kubernetesClientConfig.Timeout)

	tenantsToCreate, err := s.getTenantsToCreate()
	if err != nil {
		return err
	}
	tenantsToCreate = s.dedupeTenants(tenantsToCreate)

	tenantsToDelete, err := s.getTenantsToDelete()
	if err != nil {
		return err
	}

	deleteTenantsMap := make(map[string]model.BusinessTenantMappingInput)
	for _, ct := range tenantsToDelete {
		deleteTenantsMap[ct.ExternalTenant] = ct
	}

	for i := len(tenantsToCreate) - 1; i >= 0; i-- {
		if _, found := deleteTenantsMap[tenantsToCreate[i].ExternalTenant]; found {
			tenantsToCreate = append(tenantsToCreate[:i], tenantsToCreate[i+1:]...)
		}
	}

	tx, err := s.transact.Begin()
	if err != nil {
		return err
	}
	defer s.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	currentTenants, err := s.tenantStorageService.List(ctx)
	if err != nil {
		return errors.Wrap(err, "while listing tenants")
	}

	currentTenantsMap := make(map[string]bool)
	for _, ct := range currentTenants {
		currentTenantsMap[ct.ExternalTenant] = true
	}

	for i := len(tenantsToCreate) - 1; i >= 0; i-- {
		if currentTenantsMap[tenantsToCreate[i].ExternalTenant] {
			tenantsToCreate = append(tenantsToCreate[:i], tenantsToCreate[i+1:]...)
		}
	}

	tenantsToDelete = make([]model.BusinessTenantMappingInput, 0)
	for _, toDelete := range deleteTenantsMap {
		if currentTenantsMap[toDelete.ExternalTenant] {
			tenantsToDelete = append(tenantsToDelete, toDelete)
		}
	}

	err = s.tenantStorageService.CreateManyIfNotExists(ctx, tenantsToCreate)
	if err != nil {
		return errors.Wrap(err, "while storing new tenants")
	}
	err = s.tenantStorageService.DeleteMany(ctx, tenantsToDelete)
	if err != nil {
		return errors.Wrap(err, "while removing tenants")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func newK8SClientSet(ctx context.Context, interval, pollingTimeout, timeout time.Duration) (*kubernetes.Clientset, error) {
	k8sConfig, err := restclient.InClusterConfig()
	if err != nil {
		log.WithError(err).Warn("Failed to read in cluster Config")
		log.Info("Trying to initialize with local Config")
		home := homedir.HomeDir()
		k8sConfPath := filepath.Join(home, ".kube", "Config")
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", k8sConfPath)
		if err != nil {
			return nil, errors.Errorf("failed to read k8s in-cluster configuration, %s", err.Error())
		}
	}

	k8sConfig.Timeout = timeout

	k8sClientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, errors.Errorf("failed to create k8s core client, %s", err.Error())
	}

	err = wait.PollImmediate(interval, pollingTimeout, func() (bool, error) {
		select {
		case <-ctx.Done():
			return true, nil
		default:
		}
		_, err := k8sClientSet.ServerVersion()
		if err != nil {
			log.Debugf("Failed to access API Server: %s", err.Error())
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	log.Info("Successfully initialized kubernetes client")
	return k8sClientSet, nil
}

func getOrCreateTimestampFromConfigMap(namespace string, configMapName string, k8sClientSet *kubernetes.Clientset, defaultTimestamp int64) string {
	configMapData := make(map[string]string, 0)
	configMapData["timestamp"] = strconv.FormatInt(defaultTimestamp, 10)

	// may not be needed to create a configmap
	newConfigMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: configMapData,
	}

	configMap, err := k8sClientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})

	if k8sErrors.IsNotFound(err) {
		k8sClientSet.CoreV1().ConfigMaps(namespace).Create(&newConfigMap)
		return strconv.FormatInt(defaultTimestamp, 10)
	}

	return configMap.Data["timestamp"]
}

func updateTimestampInConfigMap(namespace string, configMapName string, k8sClientSet *kubernetes.Clientset, timestamp time.Time) error {
	configMap, err := k8sClientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	configMapData := make(map[string]string, 0)
	configMapData["timestamp"] = timestamp.String()
	configMap.Data = configMapData

	k8sClientSet.CoreV1().ConfigMaps(namespace).Update(configMap)

	return nil
}

func (s Service) getTenantsToCreate() ([]model.BusinessTenantMappingInput, error) {
	var tenantsToCreate []model.BusinessTenantMappingInput

	createdTenants, err := s.fetchTenantsWithRetries(CreatedEventsType)
	if err != nil {
		return nil, err
	}
	tenantsToCreate = append(tenantsToCreate, createdTenants...)

	updatedTenants, err := s.fetchTenantsWithRetries(UpdatedEventsType)
	if err != nil {
		return nil, err
	}
	tenantsToCreate = append(tenantsToCreate, updatedTenants...)

	return tenantsToCreate, nil
}

func (s Service) getTenantsToDelete() ([]model.BusinessTenantMappingInput, error) {
	return s.fetchTenantsWithRetries(DeletedEventsType)
}

func (s Service) fetchTenantsWithRetries(eventsType EventsType) ([]model.BusinessTenantMappingInput, error) {
	var tenants []model.BusinessTenantMappingInput
	err := retry.Do(func() error {
		fetchedTenants, err := s.fetchTenants(eventsType)
		if err != nil {
			return err
		}
		tenants = fetchedTenants
		return nil
	}, retry.Attempts(s.retryAttempts), retry.Delay(retryDelayMilliseconds*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s Service) fetchTenants(eventsType EventsType) ([]model.BusinessTenantMappingInput, error) {
	params := QueryParams{
		s.queryConfig.PageNumField:   s.queryConfig.PageStartValue,
		s.queryConfig.PageSizeField:  s.queryConfig.PageSizeValue,
		s.queryConfig.TimestampField: s.queryConfig.TimestampValue,
	}
	firstPage, err := s.eventAPIClient.FetchTenantEventsPage(eventsType, params)
	if err != nil {
		return nil, errors.Wrap(err, "while fetching tenant events page")
	}
	if firstPage == nil {
		return nil, nil
	}

	tenants := make([]model.BusinessTenantMappingInput, 0)
	tenants = append(tenants, s.extractTenantMappings(eventsType, firstPage)...)
	initialCount := gjson.GetBytes(firstPage, s.fieldMapping.TotalResultsField).Int()
	totalPages := gjson.GetBytes(firstPage, s.fieldMapping.TotalPagesField).Int()

	pageStart, err := strconv.ParseInt(s.queryConfig.PageStartValue, 10, 64)
	if err != nil {
		return nil, err
	}
	for i := pageStart + 1; i <= totalPages; i++ {
		params[s.queryConfig.PageNumField] = strconv.FormatInt(i, 10)
		res, err := s.eventAPIClient.FetchTenantEventsPage(eventsType, params)
		if err != nil {
			return nil, errors.Wrap(err, "while fetching tenant events page")
		}

		if res == nil {
			return nil, apperrors.NewInternalError("next page was expected but response was empty")
		}
		if initialCount != gjson.GetBytes(res, s.fieldMapping.TotalResultsField).Int() {
			return nil, apperrors.NewInternalError("total results number changed during fetching consecutive events pages")
		}
		tenants = append(tenants, s.extractTenantMappings(eventsType, res)...)

		// extract last event timestamp -> in configmap
	}

	return tenants, nil
}

func (s Service) extractTenantMappings(eventType EventsType, eventsJSON []byte) []model.BusinessTenantMappingInput {
	bussinessTenantMappings := make([]model.BusinessTenantMappingInput, 0)
	gjson.GetBytes(eventsJSON, s.fieldMapping.EventsField).ForEach(func(key gjson.Result, event gjson.Result) bool {
		detailsType := event.Get(s.fieldMapping.DetailsField).Type
		var details []byte
		if detailsType == gjson.String {
			details = []byte(gjson.Parse(event.Get(s.fieldMapping.DetailsField).String()).Raw)
		} else if detailsType == gjson.JSON {
			details = []byte(event.Get(s.fieldMapping.DetailsField).Raw)
		} else {
			log.Warnf("Invalid event data format: %+v", event)
			return true
		}

		tenant, err := s.eventDataToTenant(eventType, details)
		if err != nil {
			log.Warnf("Error: %s. Could not convert tenant: %s", err.Error(), string(details))
			return true
		}
		bussinessTenantMappings = append(bussinessTenantMappings, *tenant)
		return true
	})
	return bussinessTenantMappings
}

func (s Service) eventDataToTenant(eventType EventsType, eventData []byte) (*model.BusinessTenantMappingInput, error) {
	if eventType == CreatedEventsType && s.fieldMapping.DiscriminatorField != "" {
		discriminator, ok := gjson.GetBytes(eventData, s.fieldMapping.DiscriminatorField).Value().(string)
		if !ok {
			return nil, errors.Errorf("invalid format of %s field", s.fieldMapping.DiscriminatorField)
		}

		if discriminator != s.fieldMapping.DiscriminatorValue {
			return nil, nil
		}
	}

	id, ok := gjson.GetBytes(eventData, s.fieldMapping.IDField).Value().(string)
	if !ok {
		return nil, errors.Errorf("invalid format of %s field", s.fieldMapping.IDField)
	}

	name, ok := gjson.GetBytes(eventData, s.fieldMapping.NameField).Value().(string)
	if !ok {
		return nil, errors.Errorf("invalid format of %s field", s.fieldMapping.NameField)
	}

	return &model.BusinessTenantMappingInput{
		Name:           name,
		ExternalTenant: id,
		Provider:       s.providerName,
	}, nil
}

func (s Service) dedupeTenants(tenants []model.BusinessTenantMappingInput) []model.BusinessTenantMappingInput {
	elms := make(map[string]model.BusinessTenantMappingInput)
	for _, tc := range tenants {
		elms[tc.ExternalTenant] = tc
	}
	tenants = make([]model.BusinessTenantMappingInput, 0, len(elms))
	for _, t := range elms {
		tenants = append(tenants, t)
	}
	return tenants
}
