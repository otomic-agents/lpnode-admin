// Code generated by goa v3.11.0, DO NOT EDIT.
//
// HTTP request path constructors for the lpmonit service.
//
// Command:
// $ goa gen admin-panel/design

package server

// AddScriptLpmonitPath returns the URL path to the lpmonit service add_script HTTP endpoint.
func AddScriptLpmonitPath() string {
	return "/lpnode/lpnode_admin_panel/monit/script_add"
}

// ListScriptLpmonitPath returns the URL path to the lpmonit service list_script HTTP endpoint.
func ListScriptLpmonitPath() string {
	return "/lpnode/lpnode_admin_panel/monit/task_list"
}

// DeleteScriptLpmonitPath returns the URL path to the lpmonit service delete_script HTTP endpoint.
func DeleteScriptLpmonitPath() string {
	return "/lpnode/lpnode_admin_panel/monit/task_del"
}

// RunScriptLpmonitPath returns the URL path to the lpmonit service run_script HTTP endpoint.
func RunScriptLpmonitPath() string {
	return "/lpnode/lpnode_admin_panel/monit/script_run"
}

// RunResultLpmonitPath returns the URL path to the lpmonit service run_result HTTP endpoint.
func RunResultLpmonitPath() string {
	return "/lpnode/lpnode_admin_panel/monit/run_result"
}
