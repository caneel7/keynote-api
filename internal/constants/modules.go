package constants

var ResourceToScope = map[string][]string{
	"tenant":           {"X-Tenant-Id"},
	"subscriptions":    {"X-Tenant-Id"},
	"workspaces":       {"X-Tenant-Id", "X-Workspace-Id"},
	"projects":         {"X-Tenant-Id", "X-Workspace-Id"},
	"users":            {"X-Tenant-Id"},
	"roles":            {"X-Tenant-Id"},
	"audit":            {"X-Tenant-Id"},
	"project_settings": {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"project_members":  {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"project_roles":    {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"nodes":            {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"comments":         {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"attachments":      {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"workflows":        {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"sprints":          {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
	"retros":           {"X-Tenant-Id", "X-Workspace-Id", "X-Project-Id"},
}