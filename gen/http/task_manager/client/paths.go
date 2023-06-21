// Code generated by goa v3.11.0, DO NOT EDIT.
//
// HTTP request path constructors for the taskManager service.
//
// Command:
// $ goa gen admin-panel/design

package client

// TaskListTaskManagerPath returns the URL path to the taskManager service taskList HTTP endpoint.
func TaskListTaskManagerPath() string {
	return "/lpnode/lpnode_admin_panel/taskManager/list"
}

// TaskDeployTaskManagerPath returns the URL path to the taskManager service taskDeploy HTTP endpoint.
func TaskDeployTaskManagerPath() string {
	return "/lpnode/lpnode_admin_panel/taskManager/deploy"
}

// UnDeployTaskManagerPath returns the URL path to the taskManager service unDeploy HTTP endpoint.
func UnDeployTaskManagerPath() string {
	return "/lpnode/lpnode_admin_panel/taskManager/undeploy"
}

// TaskCreateTaskManagerPath returns the URL path to the taskManager service taskCreate HTTP endpoint.
func TaskCreateTaskManagerPath() string {
	return "/lpnode/lpnode_admin_panel/taskManager/create"
}
