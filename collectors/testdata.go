package collectors

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	f bool

	timestampStr = time.Now().Format(time.RFC3339)
	// Build the time from the string to ensure it is the same as a JSON time
	timestamp, _ = time.Parse(time.RFC3339, timestampStr)

	imageType     = ImageTypeAmi
	launchConsole = &f
	timeRange     = TimeRangeLasthour

	clusterID   = "12345"
	clusterID2  = &clusterID
	role        = "test"
	name        = "collector-name"
	clusterType = "collector-type"
	description = "collector-description"
	health      = "collector-health"
	cluster     = &Cluster{
		CreatedAt:    &timestamp,
		UpdatedAt:    &timestamp,
		ID:           clusterID,
		Role:         &role,
		Name:         &role,
		Type:         &clusterType,
		Description:  &description,
		Network:      network,
		Deployments:  []Deployment{*deployment},
		Status:       []Status{*status},
		Health:       &health,
		Registration: registration,
	}

	dhcp     = false
	hostname = "foobar"
	address  = "127.0.0.1"
	mask     = "255.255.255.255"
	gateway  = "0.0.0.0"
	dns      = StringSlice{"1.1.1.1"}

	network = &Network{
		Dhcp:     &dhcp,
		Hostname: &hostname,
		Hosts:    nil,
		Address:  &address,
		Mask:     &mask,
		Gateway:  &gateway,
		Dns:      &dns,
		Ntp:      nil,
		Proxy:    nil,
	}

	deploymentID          = "98765"
	deploymentName        = "deployment-name"
	deploymentDescription = "deployment-description"
	deploymentChart       = "deployment-chart"
	deploymentVersion     = "1.0"
	deploymentConfig      = Map{
		"foo": "bar",
	}
	deployment = &Deployment{
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
		ID:          deploymentID,
		Role:        nil,
		Name:        &deploymentName,
		Description: &deploymentDescription,
		Chart:       &deploymentChart,
		Version:     &deploymentVersion,
		Config:      &deploymentConfig,
		Status:      nil,
		Endpoints:   nil,
	}

	endpointID          = "87654"
	endpointDescription = ""
	endpointAddress     = "127.0.0.1"
	endpointPort        = 8080
	endpointCredentials = Map{
		"foo": "bar",
	}
	endpoint = &Endpoint{
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
		ID:          endpointID,
		Description: &endpointDescription,
		Address:     &endpointAddress,
		Port:        &endpointPort,
		Credentials: &endpointCredentials,
	}

	statusName = "status-name"
	statusID   = "34567"
	sStatus    = Map{
		"foo": "bar",
	}
	status = &Status{
		Name:      &statusName,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		ID:        statusID,
		Status:    &sStatus,
	}

	registration = &Registration{
		ID:     nil,
		Region: nil,
	}

	awsDetails = &AWSDetails{
		AccountID: "12345",
		Region:    "us-east-2",
	}

	image = &Image{
		Location: "www.example.com",
	}

	password, prk, puk = "password", "prk", "pkk"
	credentials        = &Credentials{
		Password:   &password,
		PrivateKey: &prk,
		PublicKey:  &puk,
	}

	hosts = &Hosts{
		"127.0.0.1": []string{"example.com"},
	}

	hostInput = HostsInput{
		Address:  "127.0.0.1",
		Hostname: "www.example.com",
	}

	osCfg = &OSConfig{
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		ClusterID: clusterID,
		Dhcp:      &dhcp,
		Hostname:  hostname,
		Hosts:     hosts,
		Address:   address,
		Mask:      mask,
		Gateway:   gateway,
		Dns:       &dns,
	}

	chartName        = deploymentChart
	chartDescription = deploymentDescription
	chart            = &Chart{
		Name:        &chartName,
		Description: &chartDescription,
	}

	input = OSConfigInput{
		ClusterID: clusterID,
		Dhcp:      &dhcp,
		Hostname:  &hostname,
		Address:   &address,
		Mask:      &mask,
		Gateway:   &gateway,
		Dns:       dns,
	}

	statusInput = StatusInput{
		DeploymentID: deploymentID,
		Name:         &deploymentName,
		Status: &Map{
			"foo": "bar",
		},
	}

	deploymentInput = DeploymentInput{
		Name:        &deploymentID,
		Description: &deploymentDescription,
		Chart:       &deploymentChart,
		Version:     &deploymentVersion,
		Config:      &deploymentConfig,
	}

	endpointInput = EndpointInput{
		Description: &endpointDescription,
		Address:     &endpointAddress,
		Port:        &endpointPort,
		Credentials: &endpointCredentials,
	}

	updateEndpointArgsInput = &UpdateEndpointArguments{
		ClusterID:     clusterID,
		DeploymentID:  deploymentID,
		EndpointID:    endpointID,
		EndpointInput: endpointInput,
	}

	getClusterResponse = struct {
		Out *Cluster `json:"getCluster"`
	}{
		Out: cluster,
	}

	getAllClustersResponse = struct {
		Out []Cluster `json:"getAllClusters"`
	}{
		Out: []Cluster{*cluster},
	}

	getClusterConfigResponse = struct {
		Out *KubernetesConfig `json:"getClusterConfig"`
	}{
		Out: &KubernetesConfig{
			TypeMeta: v1.TypeMeta{},
		},
	}

	clusterImageArgs = &GetClusterImageArguments{
		ClusterID:     clusterID,
		ImageType:     imageType,
		LaunchConsole: launchConsole,
		AwsDetails:    awsDetails,
	}

	getClusterImageResponse = struct {
		Out *Image `json:"getClusterImage"`
	}{
		Out: image,
	}

	getClusterCredentialsResponse = struct {
		Out *Credentials `json:"getClusterCredentials"`
	}{
		Out: credentials,
	}

	getHostsResponse = struct {
		Out *Hosts `json:"getHosts"`
	}{
		Out: hosts,
	}

	getOSConfigResponse = struct {
		Out *OSConfig `json:"getOSConfig"`
	}{
		Out: osCfg,
	}

	getClusterStatusesResponse = struct {
		Out []Status `json:"getClusterStatuses"`
	}{
		Out: []Status{*status},
	}

	getClusterDeploymentStatusResponse = struct {
		Out *Map `json:"getClusterDeploymentStatus"`
	}{
		Out: &Map{},
	}

	getChartResponse = struct {
		Out *Chart `json:"getChart"`
	}{
		Out: &Chart{},
	}

	getAllChartsResponse = struct {
		Out *ChartList `json:"getAllCharts"`
	}{
		Out: &ChartList{},
	}

	getClusterDeploymentResponse = struct {
		Out *Deployment `json:"getClusterDeployment"`
	}{
		Out: deployment,
	}

	getAllClusterDeploymentsResponse = struct {
		Out []Deployment `json:"getAllClusterDeployments"`
	}{
		Out: []Deployment{*deployment},
	}

	getDeploymentEndpointResponse = struct {
		Out *Endpoint `json:"getDeploymentEndpoint"`
	}{
		Out: endpoint,
	}

	getAllDeploymentEndpointsResponse = struct {
		Out []Endpoint `json:"getAllDeploymentEndpoints"`
	}{
		Out: []Endpoint{*endpoint},
	}

	getAWSRegionsResponse = struct {
		Out []string `json:"getAWSRegions"`
	}{
		Out: []string{"us-west-2"},
	}

	getRoleDeploymentsResponse = struct {
		Out []Deployment `json:"getRoleDeployments"`
	}{
		Out: []Deployment{*deployment},
	}

	getRoleDeploymentResponse = struct {
		Out *Deployment `json:"getRoleDeployment"`
	}{
		Out: deployment,
	}

	getAllCollectorsOverviewResponse = struct {
		Out []CollectorOverview `json:"getAllCollectorsOverview"`
	}{
		Out: []CollectorOverview{{
			Cluster:     *cluster,
			LastSeen:    nil,
			AverageRate: nil,
		}},
	}

	getCollectorMetricsResponse = struct {
		Out *CollectorMetrics `json:"getCollectorMetrics"`
	}{
		Out: &CollectorMetrics{},
	}

	getAggregateRateByCollectorResponse = struct {
		Out *AggregateRateByCollector `json:"getAggregateRateByCollector"`
	}{
		Out: &AggregateRateByCollector{},
	}

	getFlowRateResponse = struct {
		Out *FlowRate `json:"getFlowRate"`
	}{
		Out: &FlowRate{},
	}

	getLogLastSeenMetricsResponse = struct {
		Out *LogLastSeenMetrics `json:"getLogLastSeenMetrics"`
	}{
		Out: &LogLastSeenMetrics{},
	}

	clusterInput = ClusterInput{
		Name:        &name,
		Description: &description,
		Role:        &role,
	}

	createClusterResponse = struct {
		Out *Cluster `json:"createCluster"`
	}{
		Out: cluster,
	}

	updateClusterResponse = struct {
		Out *Cluster `json:"updateCluster"`
	}{
		Out: cluster,
	}

	deleteClusterResponse = struct {
		Out *Deleted `json:"deleteCluster"`
	}{
		Out: &Deleted{
			Type:       "cluster",
			ID:         clusterID,
			Successful: true,
		},
	}

	createOSConfigResponse = struct {
		Out *OSConfig `json:"createOSConfig"`
	}{
		Out: osCfg,
	}

	updateOSConfigResponse = struct {
		Out *OSConfig `json:"updateOSConfig"`
	}{
		Out: osCfg,
	}

	deleteOSConfigResponse = struct {
		Out string `json:"deleteOSConfig"`
	}{
		Out: clusterID,
	}

	addHostResponse = struct {
		Out *Hosts `json:"addHost"`
	}{
		Out: &Hosts{
			"127.0.0.1": []string{"www.example.com"},
		},
	}

	deleteHostResponse = struct {
		Out *Deleted `json:"deleteHost"`
	}{
		Out: &Deleted{
			Type:       "hosts",
			ID:         "127.0.0.1",
			Successful: true,
		},
	}

	createClusterStatusResponse = struct {
		Out *Status `json:"createClusterStatus"`
	}{
		Out: status,
	}

	updateClusterStatusResponse = struct {
		Out *Status `json:"updateClusterStatus"`
	}{
		Out: status,
	}

	deleteClusterStatusResponse = struct {
		Out *Deleted `json:"deleteClusterStatus"`
	}{
		Out: &Deleted{
			Type:       "status",
			ID:         statusID,
			Successful: true,
		},
	}

	createClusterDeploymentResponse = struct {
		Out *Deployment `json:"createClusterDeployment"`
	}{
		Out: deployment,
	}

	updateClusterDeploymentResponse = struct {
		Out *Deployment `json:"updateClusterDeployment"`
	}{
		Out: deployment,
	}

	deleteClusterDeploymentResponse = struct {
		Out *Deleted `json:"deleteClusterDeployment"`
	}{
		Out: &Deleted{
			Type:       "deployment",
			ID:         deploymentID,
			Successful: true,
		},
	}

	createRoleDeploymentResponse = struct {
		Out *Deployment `json:"createRoleDeployment"`
	}{
		Out: deployment,
	}

	updateRoleDeploymentResponse = struct {
		Out *Deployment `json:"updateRoleDeployment"`
	}{
		Out: deployment,
	}

	deleteRoleDeploymentResponse = struct {
		Out *Deleted `json:"deleteRoleDeployment"`
	}{
		Out: &Deleted{
			Type:       "role deployment",
			ID:         deploymentID,
			Successful: true,
		},
	}

	createEndpointResponse = struct {
		Out *Endpoint `json:"createEndpoint"`
	}{
		Out: endpoint,
	}

	updateEndpointResponse = struct {
		Out *Endpoint `json:"updateEndpoint"`
	}{
		Out: endpoint,
	}

	deleteEndpointResponse = struct {
		Out *Deleted `json:"deleteEndpoint"`
	}{
		Out: &Deleted{
			Type:       "endpoint",
			ID:         endpointID,
			Successful: true,
		},
	}
)
