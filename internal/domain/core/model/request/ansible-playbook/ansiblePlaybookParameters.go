package ansibleplaybook

type AnsiblePlaybookParameters struct {

	// Playbooks is the ansible's playbooks list to be executed
	Playbook []string `json:"playbooks" validate:"required"`

	// // AskVaultPassword ask for vault password
	// AskVaultPassword bool

	// Check don't make any changes; instead, try to predict some of the changes that may occur
	Check bool `json:"check,omitempty" validate:"bool"`

	// Diff when changing (small) files and templates, show the differences in those files; works great with --check
	Diff bool `json:"diff,omitempty" validate:"bool"`

	// Dependencies is a list of role and collection dependencies
	Dependencies interface{} `json:"dependencies,omitempty"`

	// ExtraVars is a map of extra variables used on ansible-playbook execution
	ExtraVars map[string]interface{} `json:"extra_vars,omitempty"`

	// ExtraVarsFile is a list of files used to load extra-vars
	ExtraVarsFile []string `json:"extra_vars_file,omitempty"`

	// FlushCache is the flush cache flag for ansible-playbook
	FlushCache bool `json:"flush_cache,omitempty" validate:"bool"`

	// ForceHandlers run handlers even if a task fails
	ForceHandlers bool `json:"force_handlers,omitempty" validate:"bool"`

	// Forks specify number of parallel processes to use (default=50)
	Forks string `json:"forks,omitempty" validate:"numeric"`

	// Inventory specify inventory host path
	Inventory string `json:"inventory,omitempty" validate:"required"`

	// Limit is selected hosts additional pattern
	Limit string `json:"limit,omitempty"`

	// ListHosts outputs a list of matching hosts
	ListHosts bool `json:"list_hosts,omitempty" validate:"bool"`

	// ListTags is the list tags flag for ansible-playbook
	ListTags bool `json:"list_tags,omitempty" validate:"bool"`

	// ListTasks is the list tasks flag for ansible-playbook
	ListTasks bool `json:"list_tasks,omitempty" validate:"bool"`

	// ModulePath repend colon-separated path(s) to module library (default=~/.ansible/plugins/modules:/usr/share/ansible/plugins/modules)
	ModulePath string `json:"module_path,omitempty"`

	// SkipTags only run plays and tasks whose tags do not match these values
	SkipTags string `json:"skip_tags,omitempty"`

	// StartAtTask start the playbook at the task matching this name
	StartAtTask string `json:"start_at_task,omitempty"`

	// // Step one-step-at-a-time: confirm each task before running
	// Step bool

	// SyntaxCheck is the syntax check flag for ansible-playbook
	SyntaxCheck bool `json:"syntax_check,omitempty" validate:"bool"`

	// Tags is the tags flag for ansible-playbook
	Tags string `json:"tags,omitempty"`

	// VaultID the vault identity to use
	VaultID string `json:"vault_id,omitempty"`

	// VaultPasswordFile path to the file holding vault decryption key
	VaultPasswordFile string `json:"vault_password_file,omitempty"`

	// Verbose verbose mode enabled
	Verbose bool `json:"verbose,omitempty"`

	// // Verbose verbose mode -v enabled
	// VerboseV bool

	// // Verbose verbose mode -vv enabled
	// VerboseVV bool

	// // Verbose verbose mode -vvv enabled
	// VerboseVVV bool

	// // Verbose verbose mode -vvvv enabled
	// VerboseVVVV bool

	// Version show program's version number, config file location, configured module search path, module location, executable location and exit
	Version bool `json:"version,omitempty" validate:"bool"`

	// Parameters defined on `Connections Options` section within ansible-playbook's man page, and which defines how to connect to hosts.

	// // AskPass defines whether user's password should be asked to connect to host
	// AskPass bool

	// // Connection is the type of connection used by ansible-playbook
	// Connection string

	// // PrivateKey is the user's private key file used to connect to a host
	// PrivateKey string

	// SCPExtraArgs specify extra arguments to pass to scp only
	SCPExtraArgs string `json:"scp_extra_args,omitempty"`

	// SFTPExtraArgs specify extra arguments to pass to sftp only
	SFTPExtraArgs string `json:"sftp_extra_args,omitempty"`

	// SSHCommonArgs specify common arguments to pass to sftp/scp/ssh
	SSHCommonArgs string `json:"ssh_common_args,omitempty"`

	// SSHExtraArgs specify extra arguments to pass to ssh only
	SSHExtraArgs string `json:"ssh_extra_args,omitempty"`

	// Timeout is the connection timeout on ansible-playbook. Take care because Timeout is defined ad string
	Timeout int `json:"timeout,omitempty" validate:"numeric"`

	// User is the user to use to connect to a host
	User string `json:"user,omitempty"`

	// Parameters defined on `Privilege Escalation Options` section within ansible-playbook's man page, and which controls how and which user you become as on target hosts.

	// // AskBecomePass is ansble-playbook's ask for become user password flag
	// AskBecomePass bool

	// Become is ansble-playbook's become flag
	Become bool `json:"become,omitempty" validate:"bool"`

	// BecomeMethod is ansble-playbook's become method. The accepted become methods are:
	// 	- ksu        Kerberos substitute user
	// 	- pbrun      PowerBroker run
	// 	- enable     Switch to elevated permissions on a network device
	// 	- sesu       CA Privileged Access Manager
	// 	- pmrun      Privilege Manager run
	// 	- runas      Run As user
	// 	- sudo       Substitute User DO
	// 	- su         Substitute User
	// 	- doas       Do As user
	// 	- pfexec     profile based execution
	// 	- machinectl Systemd's machinectl privilege escalation
	// 	- dzdo       Centrify's Direct Authorize
	BecomeMethod string `json:"become_method,omitempty"`

	// BecomeUser is ansble-playbook's become user
	BecomeUser string `json:"become_user,omitempty"`
}
