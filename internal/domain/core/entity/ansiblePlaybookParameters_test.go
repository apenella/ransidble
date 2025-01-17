package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityAnsiblePlaybookParametersValidate(t *testing.T) {
	type fields struct {
		Become        bool
		Check         bool
		Diff          bool
		FlushCache    bool
		ForceHandlers bool
		Forks         int
		Inventory     string
		ListHosts     bool
		ListTags      bool
		ListTasks     bool
		Playbooks     []string
		Requirements  *AnsiblePlaybookRequirements
		SyntaxCheck   bool
		Timeout       int
		Verbose       bool
		Version       bool
	}
	test := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a AnsiblePlaybookParameters entity",
			fields: fields{
				Become:        false,
				Check:         false,
				Diff:          false,
				FlushCache:    false,
				ForceHandlers: false,
				Forks:         5,
				Inventory:     "inventory",
				ListHosts:     false,
				ListTags:      false,
				ListTasks:     false,
				Playbooks:     []string{"playbook.yml"},
				Requirements: &AnsiblePlaybookRequirements{
					Roles: &AnsiblePlaybookRoleRequirements{
						Roles:        []string{"roles"},
						APIKey:       "apikey",
						IgnoreErrors: true,
						NoDeps:       true,
						RoleFile:     "rolefile",
						Server:       "server",
						Timeout:      "30",
						Token:        "token",
						Verbose:      true,
					},
					Collections: &AnsiblePlaybookCollectionRequirements{
						Collections:      []string{"collections"},
						APIKey:           "apikey",
						ForceWithDeps:    true,
						Pre:              true,
						Timeout:          "30",
						Token:            "token",
						IgnoreErrors:     true,
						RequirementsFile: "requirementsfile",
						Server:           "server",
						Verbose:          true,
					},
				},
				SyntaxCheck: false,
				Timeout:     30,
				Verbose:     false,
				Version:     false,
			},
			wantErr: false,
		},
		{
			desc: "Validating a AnsiblePlaybookParameters entity with empty playbooks",
			fields: fields{
				Check:         false,
				Diff:          false,
				FlushCache:    false,
				ForceHandlers: false,
				Forks:         5,
				Inventory:     "inventory",
				ListHosts:     false,
				ListTags:      false,
				ListTasks:     false,
				SyntaxCheck:   false,
				Verbose:       false,
				Version:       false,
				Timeout:       30,
				Become:        false,
			},
			wantErr: true,
		},
		{
			desc: "Validating a AnsiblePlaybookParameters entity with empty inventory",
			fields: fields{
				Playbooks:     []string{"playbook.yml"},
				Check:         false,
				Diff:          false,
				FlushCache:    false,
				ForceHandlers: false,
				Forks:         5,
				ListHosts:     false,
				ListTags:      false,
				ListTasks:     false,
				SyntaxCheck:   false,
				Verbose:       false,
				Version:       false,
				Timeout:       30,
				Become:        false,
			},
			wantErr: true,
		},
	}

	for _, test := range test {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			params := &AnsiblePlaybookParameters{
				Playbooks:     test.fields.Playbooks,
				Check:         test.fields.Check,
				Diff:          test.fields.Diff,
				FlushCache:    test.fields.FlushCache,
				ForceHandlers: test.fields.ForceHandlers,
				Forks:         test.fields.Forks,
				Inventory:     test.fields.Inventory,
				ListHosts:     test.fields.ListHosts,
				ListTags:      test.fields.ListTags,
				ListTasks:     test.fields.ListTasks,
				SyntaxCheck:   test.fields.SyntaxCheck,
				Verbose:       test.fields.Verbose,
				Version:       test.fields.Version,
				Timeout:       test.fields.Timeout,
				Become:        test.fields.Become,
				Requirements:  test.fields.Requirements,
			}

			err := params.Validate()
			if err != nil {
				assert.Equal(t, test.wantErr, true, err.Error())
			}

			if test.wantErr {
				assert.NotNil(t, err)
			}

		})
	}
}
