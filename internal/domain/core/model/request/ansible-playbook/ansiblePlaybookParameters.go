package ansibleplaybook

import (
	_ "embed"

	"github.com/go-playground/validator/v10"
)

type AnsiblePlaybookParameters struct {

	// Playbooks is the ansible's playbooks list to be executed
	Playbooks []string `json:"playbooks" validate:"required"`

	// // AskVaultPassword ask for vault password
	// AskVaultPassword bool

	// Check don't make any changes; instead, try to predict some of the changes that may occur
	Check bool `json:"check,omitempty" validate:"boolean"`

	// Diff when changing (small) files and templates, show the differences in those files; works great with --check
	Diff bool `json:"diff,omitempty" validate:"boolean"`

	// Dependencies is a list of role and collection dependencies
	Dependencies *AnsiblePlaybookDependencies `json:"dependencies,omitempty"`

	// ExtraVars is a map of extra variables used on ansible-playbook execution
	ExtraVars map[string]interface{} `json:"extra_vars,omitempty"`

	// ExtraVarsFile is a list of files used to load extra-vars
	ExtraVarsFile []string `json:"extra_vars_file,omitempty"`

	// FlushCache is the flush cache flag for ansible-playbook
	FlushCache bool `json:"flush_cache,omitempty" validate:"boolean"`

	// ForceHandlers run handlers even if a task fails
	ForceHandlers bool `json:"force_handlers,omitempty" validate:"boolean"`

	// Forks specify number of parallel processes to use (default=50)
	Forks int `json:"forks,omitempty" validate:"number"`

	// Inventory specify inventory host path
	Inventory string `json:"inventory,omitempty" validate:"required"`

	// Limit is selected hosts additional pattern
	Limit string `json:"limit,omitempty"`

	// ListHosts outputs a list of matching hosts
	ListHosts bool `json:"list_hosts,omitempty" validate:"boolean"`

	// ListTags is the list tags flag for ansible-playbook
	ListTags bool `json:"list_tags,omitempty" validate:"boolean"`

	// ListTasks is the list tasks flag for ansible-playbook
	ListTasks bool `json:"list_tasks,omitempty" validate:"boolean"`

	// // ModulePath repend colon-separated path(s) to module library (default=~/.ansible/plugins/modules:/usr/share/ansible/plugins/modules)
	// ModulePath string `json:"module_path,omitempty"`

	// SkipTags only run plays and tasks whose tags do not match these values
	SkipTags string `json:"skip_tags,omitempty"`

	// StartAtTask start the playbook at the task matching this name
	StartAtTask string `json:"start_at_task,omitempty"`

	// // Step one-step-at-a-time: confirm each task before running
	// Step bool

	// SyntaxCheck is the syntax check flag for ansible-playbook
	SyntaxCheck bool `json:"syntax_check,omitempty" validate:"boolean"`

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
	Version bool `json:"version,omitempty" validate:"boolean"`

	// Parameters defined on `Connections Options` section within ansible-playbook's man page, and which defines how to connect to hosts.

	// // AskPass defines whether user's password should be asked to connect to host
	// AskPass bool

	// Connection is the type of connection used by ansible-playbook
	Connection string `json:"connection,omitempty" validate:"alphanum"`

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
	Become bool `json:"become,omitempty" validate:"boolean"`

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

type AnsiblePlaybookDependencies struct {
	// Roles defines how to install roles dependencies
	Roles *AnsiblePlaybookRoleDependencies `json:"roles,omitempty"`
	// Collections defines how to install collections dependencies
	Collections *AnsiblePlaybookCollectionDependencies `json:"collections,omitempty"`
}

type AnsiblePlaybookRoleDependencies struct {

	// Roles is a list of roles to install
	Roles []string `json:"roles,omitempty"`

	// ApiKey represent the API key to use to authenticate against the galaxy server. Same as --token
	ApiKey string `json:"api_key,omitempty"`

	// // Force represents whether to force overwriting an existing role or role file.
	// Force bool

	// // ForceWithDeps represents whether to force overwriting an existing role, role file, or dependencies.
	// ForceWithDeps bool

	// // IgnoreCerts represent the flag to ignore SSL certificate validation errors
	// IgnoreCerts bool

	// IgnoreErrors represents whether to continue processing even if a role fails to install.
	IgnoreErrors bool `json:"ignore_errors,omitempty" validate:"boolean"`

	// // KeepSCMMeta represent the flag to use tar instead of the scm archive option when packaging the role.
	// KeepSCMMeta bool

	// NoDeps represents whether to install dependencies.
	NoDeps bool `json:"no_deps,omitempty" validate:"boolean"`

	// RoleFile represents the path to a file containing a list of roles to install.
	RoleFile string `json:"role_file,omitempty"`

	// // RolesPath represents the path where roles should be installed on the local filesystem.
	// RolesPath string

	// Server represent the flag to specify the galaxy server to use
	Server string `json:"server,omitempty"`

	// Timeout represent the time to wait for operations against the galaxy server, defaults to 60s
	Timeout string `json:"timeout,omitempty" validate:"numeric"`

	// Token represent the token to use to authenticate against the galaxy server. Same as --api-key
	Token string `json:"token,omitempty"`

	// Verbose verbose mode enabled
	Verbose bool `json:"verbose,omitempty" validate:"boolean"`

	// // Verbose verbose mode -v enabled
	// VerboseV bool

	// // Verbose verbose mode -vv enabled
	// VerboseVV bool

	// // Verbose verbose mode -vvv enabled
	// VerboseVVV bool

	// // Verbose verbose mode -vvvv enabled
	// VerboseVVVV bool

	// // Version show program's version number, config file location, configured module search path, module location, executable location and exit
	// Version bool
}

type AnsiblePlaybookCollectionDependencies struct {

	// Collections is a list of collections to install.
	Collections []string `json:"collections,omitempty"`

	// APIKey is the Ansible Galaxy API key.
	APIKey string `json:"api_key,omitempty"`

	// // ClearResponseCache clears the existing server response cache.
	// ClearResponseCache bool

	// // DisableGPGVerify disables GPG signature verification when installing collections from a Galaxy server.
	// DisableGPGVerify bool

	// ForceWithDeps forces overwriting an existing collection and its dependencies.
	ForceWithDeps bool `json:"force_with_deps,omitempty" validate:"boolean"`

	// // IgnoreSignatureStatusCode suppresses this argument. It may be specified multiple times.
	// IgnoreSignatureStatusCode bool

	// // IgnoreSignatureStatusCodes is a space separated list of status codes to ignore during signature verification.
	// IgnoreSignatureStatusCodes string

	// // Keyring is the keyring used during signature verification.
	// Keyring string

	// // NoCache does not use the server response cache.
	// NoCache bool

	// // Offline installs collection artifacts (tarballs) without contacting any distribution servers.
	// Offline bool

	// Pre includes pre-release versions. Semantic versioning pre-releases are ignored by default.
	Pre bool `json:"pre,omitempty" validate:"boolean"`

	// // RequiredValidSignatureCount is the number of signatures that must successfully verify the collection.
	// RequiredValidSignatureCount int

	// // Signature is an additional signature source to verify the authenticity of the MANIFEST.json.
	// Signature string

	// Timeout is the time to wait for operations against the galaxy server, defaults to 60s.
	Timeout string `json:"timeout,omitempty"`

	// Token is the Ansible Galaxy API key.
	Token string `json:"token,omitempty"`

	// // Upgrade upgrades installed collection artifacts. This will also update dependencies unless –no-deps is provided.
	// Upgrade bool

	// // IgnoreCerts ignores SSL certificate validation errors.
	// IgnoreCerts bool

	// // Force forces overwriting an existing role or collection.
	// Force bool

	// IgnoreErrors ignores errors during installation and continue with the next specified collection.
	IgnoreErrors bool `json:"ignore_errors,omitempty" validate:"boolean"`

	// // NoDeps doesn’t download collections listed as dependencies.
	// NoDeps bool

	// // CollectionsPath is the path to the directory containing your collections.
	// CollectionsPath string

	// RequirementsFile is a file containing a list of collections to be installed.
	RequirementsFile string `json:"requirements_file,omitempty"`

	// Server is the Galaxy API server URL.
	Server string `json:"server,omitempty"`

	// Verbose verbose mode enabled
	Verbose bool `json:"verbose,omitempty" validate:"boolean"`

	// // Version show program's version number, config file location, configured module search path, module location, executable location and exit
	// Version bool
}

func (params *AnsiblePlaybookParameters) Validate() error {
	validate := validator.New()
	return validate.Struct(params)
}
