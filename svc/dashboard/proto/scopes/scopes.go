// Code generated by lab.weave.nl/devops/proto-istio-auth-generator, DO NOT EDIT.
package dashboardscopes

// DashboardDeployService contains the scope for DeployService in the Dashboard service
const DashboardDeployService = "deploy_service"	

// DashboardDeleteService contains the scope for DeleteService in the Dashboard service
const DashboardDeleteService = "delete_service"	

// DashboardListNamespaces contains the scope for ListNamespaces in the Dashboard service
const DashboardListNamespaces = "list_namespaces"	

// DashboardListServices contains the scope for ListServices in the Dashboard service
const DashboardListServices = "list_services"	

// GetAllScopes retrieve all available scopes
func GetAllScopes() []string{
	return []string{
			DashboardDeployService,
			DashboardDeleteService,
			DashboardListNamespaces,
			DashboardListServices,
	}
}
