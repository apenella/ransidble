package mapper

import (
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
	"github.com/stretchr/testify/assert"
)

// TestToAnsiblePlaybookParametersEntity tests ToAnsiblePlaybookParametersEntity method
func TestToAnsiblePlaybookParametersEntity(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   *request.AnsiblePlaybookParameters
		expected *entity.AnsiblePlaybookParameters
	}{
		{
			desc:   "Testing to ansible playbook parameters entity with all fields",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookParameters{
				Playbooks: []string{"playbook1", "playbook2"},
				Check:     true,
				Diff:      true,
				Requirements: &request.AnsiblePlaybookRequirements{
					Roles: &request.AnsiblePlaybookRoleRequirements{
						Roles: []string{"role1", "role2"},
					},
					Collections: &request.AnsiblePlaybookCollectionRequirements{
						Collections: []string{"collection1", "collection2"},
					},
				},
				ExtraVars:         map[string]interface{}{"key": "value"},
				ExtraVarsFile:     []string{"extra-vars-file1", "extra-vars-file2"},
				FlushCache:        true,
				ForceHandlers:     true,
				Forks:             10,
				Inventory:         "inventory",
				Limit:             "limit",
				ListHosts:         true,
				ListTags:          true,
				ListTasks:         true,
				SkipTags:          "skip-tags",
				StartAtTask:       "start-at-task",
				SyntaxCheck:       true,
				Tags:              "tags",
				VaultID:           "vault-id",
				VaultPasswordFile: "vault-password-file",
				Verbose:           true,
				Version:           true,
				Connection:        "connection",
				SCPExtraArgs:      "scp-extra-args",
				SFTPExtraArgs:     "sftp-extra-args",
				SSHCommonArgs:     "ssh-common-args",
				SSHExtraArgs:      "ssh-extra-args",
				Timeout:           10,
				User:              "user",
				Become:            true,
				BecomeMethod:      "become-method",
				BecomeUser:        "become-user",
			},
			expected: &entity.AnsiblePlaybookParameters{
				Playbooks: []string{"playbook1", "playbook2"},
				Check:     true,
				Diff:      true,
				Requirements: &entity.AnsiblePlaybookRequirements{
					Roles: &entity.AnsiblePlaybookRoleRequirements{
						Roles: []string{"role1", "role2"},
					},
					Collections: &entity.AnsiblePlaybookCollectionRequirements{
						Collections: []string{"collection1", "collection2"},
					},
				},
				ExtraVars:         map[string]interface{}{"key": "value"},
				ExtraVarsFile:     []string{"extra-vars-file1", "extra-vars-file2"},
				FlushCache:        true,
				ForceHandlers:     true,
				Forks:             10,
				Inventory:         "inventory",
				Limit:             "limit",
				ListHosts:         true,
				ListTags:          true,
				ListTasks:         true,
				SkipTags:          "skip-tags",
				StartAtTask:       "start-at-task",
				SyntaxCheck:       true,
				Tags:              "tags",
				VaultID:           "vault-id",
				VaultPasswordFile: "vault-password-file",
				Verbose:           true,
				Version:           true,
				Connection:        "connection",
				SCPExtraArgs:      "scp-extra-args",
				SFTPExtraArgs:     "sftp-extra-args",
				SSHCommonArgs:     "ssh-common-args",
				SSHExtraArgs:      "ssh-extra-args",
				Timeout:           10,
				User:              "user",
				Become:            true,
				BecomeMethod:      "become-method",
				BecomeUser:        "become-user",
			},
		},
		{
			desc:     "Testing to ansible playbook parameters entity with nil source",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: &entity.AnsiblePlaybookParameters{},
		},
		{
			desc:   "Testing to ansible playbook parameters entity with empty source",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookParameters{},
			expected: &entity.AnsiblePlaybookParameters{
				Playbooks:     []string{},
				Requirements:  &entity.AnsiblePlaybookRequirements{},
				ExtraVars:     map[string]interface{}{},
				ExtraVarsFile: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.ToAnsiblePlaybookParametersEntity(test.source)

			assert.Equal(t, test.expected, res)
		})
	}
}

// TestToAnsiblePLaybookParametersRequirementsEntity tests ToAnsiblePLaybookParametersRequirementsEntity method
func TestToAnsiblePLaybookParametersRequirementsEntity(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   *request.AnsiblePlaybookRequirements
		expected *entity.AnsiblePlaybookRequirements
	}{
		{
			desc:   "Testing to ansible playbook parameters dependencies entity with all fields",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookRequirements{
				Roles: &request.AnsiblePlaybookRoleRequirements{
					Roles: []string{"role1", "role2"},
				},
				Collections: &request.AnsiblePlaybookCollectionRequirements{
					Collections: []string{"collection1", "collection2"},
				},
			},
			expected: &entity.AnsiblePlaybookRequirements{
				Roles: &entity.AnsiblePlaybookRoleRequirements{
					Roles: []string{"role1", "role2"},
				},
				Collections: &entity.AnsiblePlaybookCollectionRequirements{
					Collections: []string{"collection1", "collection2"},
				},
			},
		},
		{
			desc:     "Testing to ansible playbook parameters dependencies entity with nil source",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: &entity.AnsiblePlaybookRequirements{},
		},
		{
			desc:   "Testing to ansible playbook parameters dependencies entity with empty source",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookRequirements{},
			expected: &entity.AnsiblePlaybookRequirements{
				Roles:       &entity.AnsiblePlaybookRoleRequirements{},
				Collections: &entity.AnsiblePlaybookCollectionRequirements{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.toAnsiblePLaybookParametersRequirementsEntity(test.source)

			assert.Equal(t, test.expected, res)
		})
	}
}

// TestToAnsiblePLaybookParametersRolesRequirementsEntity tests ToAnsiblePLaybookParametersRolesRequirementsEntity method
func TestToAnsiblePLaybookParametersRolesRequirementsEntity(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   *request.AnsiblePlaybookRoleRequirements
		expected *entity.AnsiblePlaybookRoleRequirements
	}{
		{
			desc:   "Testing to ansible playbook parameters roles dependencies entity with all fields",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookRoleRequirements{
				Roles:        []string{"role1", "role2"},
				APIKey:       "apikey",
				IgnoreErrors: true,
				NoDeps:       true,
				RoleFile:     "roles.yml",
				Server:       "server",
				Timeout:      10,
				Token:        "token",
				Verbose:      true,
			},
			expected: &entity.AnsiblePlaybookRoleRequirements{
				Roles:        []string{"role1", "role2"},
				APIKey:       "apikey",
				IgnoreErrors: true,
				NoDeps:       true,
				RoleFile:     "roles.yml",
				Server:       "server",
				Timeout:      10,
				Token:        "token",
				Verbose:      true,
			},
		},
		{
			desc:     "Testing to ansible playbook parameters roles dependencies entity with nil source",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: &entity.AnsiblePlaybookRoleRequirements{},
		},
		{
			desc:   "Testing to ansible playbook parameters roles dependencies entity with empty source",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookRoleRequirements{},
			expected: &entity.AnsiblePlaybookRoleRequirements{
				Roles: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.toAnsiblePLaybookParametersRolesRequirementsEntity(test.source)

			assert.Equal(t, test.expected, res)
		})
	}
}

// TestToAnsiblePLaybookParametersCollectionsRequirementsEntity tests ToAnsiblePLaybookParametersCollectionsRequirementsEntity method
func TestToAnsiblePLaybookParametersCollectionsRequirementsEntity(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   *request.AnsiblePlaybookCollectionRequirements
		expected *entity.AnsiblePlaybookCollectionRequirements
	}{
		{
			desc:   "Testing to ansible playbook parameters collections dependencies entity with all fields",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookCollectionRequirements{
				Collections:      []string{"collection1", "collection2"},
				APIKey:           "apikey",
				ForceWithDeps:    true,
				Pre:              true,
				Timeout:          10,
				Token:            "token",
				IgnoreErrors:     true,
				RequirementsFile: "requirements.yml",
				Server:           "server",
				Verbose:          true,
			},
			expected: &entity.AnsiblePlaybookCollectionRequirements{
				Collections:      []string{"collection1", "collection2"},
				APIKey:           "apikey",
				ForceWithDeps:    true,
				Pre:              true,
				Timeout:          10,
				Token:            "token",
				IgnoreErrors:     true,
				RequirementsFile: "requirements.yml",
				Server:           "server",
				Verbose:          true,
			},
		},
		{
			desc:     "Testing to ansible playbook parameters collections dependencies entity with nil source",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: &entity.AnsiblePlaybookCollectionRequirements{},
		},
		{
			desc:   "Testing to ansible playbook parameters collections dependencies entity with empty source",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: &request.AnsiblePlaybookCollectionRequirements{},
			expected: &entity.AnsiblePlaybookCollectionRequirements{
				Collections: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.toAnsiblePLaybookParametersCollectionsRequirementsEntity(test.source)

			assert.Equal(t, test.expected, res)
		})
	}

}

// TestToAnsiblePlaybookParametersExtraVarsEntity tests ToAnsiblePlaybookParametersExtraVarsEntity method
func TestToAnsiblePlaybookParametersExtraVarsEntity(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   map[string]interface{}
		expected map[string]interface{}
	}{
		{
			desc:   "Testing to ansible playbook parameters extra vars entity",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"string1", "string2"},
				"key4": map[string]string{"key": "value"},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"string1", "string2"},
				"key4": map[string]string{"key": "value"},
			},
		},
		{
			desc:     "Testing to ansible playbook parameters extra vars entity when recieve a nil map",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: map[string]interface{}{},
		},
		{
			desc:     "Testing to ansible playbook parameters extra vars entity when recieve an empty map",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.toAnsiblePlaybookParametersExtraVarsEntity(test.source)

			assert.Equal(t, test.expected, res)
		})
	}

}

// TestCopyMap tests copySlice method
func TestCopyMap(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   map[string]interface{}
		expected map[string]interface{}
	}{
		{
			desc:   "Testing copy map",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"string1", "string2"},
				"key4": map[string]string{"key": "value"},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"string1", "string2"},
				"key4": map[string]string{"key": "value"},
			},
		},
		{
			desc:     "Testing copy map when recieve a nil map",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: map[string]interface{}{},
		},
		{
			desc:     "Testing copy map when recieve an empty map",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.copyMap(test.source)
			assert.Equal(t, test.expected, res)
		})
	}

}

// TestCopySlice tests copySlice method
func TestCopySlice(t *testing.T) {
	tests := []struct {
		desc     string
		mapper   *AnsiblePlaybookParametersMapper
		source   []interface{}
		expected []interface{}
	}{
		{
			desc:   "Testing copy slice",
			mapper: NewAnsiblePlaybookParametersMapper(),
			source: []interface{}{
				"string",
				1,
				[]string{"string1", "string2"},
				map[string]string{"key": "value"},
			},
			expected: []interface{}{
				"string",
				1,
				[]string{"string1", "string2"},
				map[string]string{"key": "value"},
			},
		},
		{
			desc:     "Testing copy slice when recieve a nil slice",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   nil,
			expected: []interface{}{},
		},
		{
			desc:     "Testing copy slice when recieve an empty slice",
			mapper:   NewAnsiblePlaybookParametersMapper(),
			source:   []interface{}{},
			expected: []interface{}{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.copySlice(test.source)
			assert.Equal(t, test.expected, res)
		})
	}

}
