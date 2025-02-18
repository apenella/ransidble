package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequsetAnsiblePlaybookParametersValidate(t *testing.T) {
	type fields struct {
		Playbooks     []string
		Check         bool
		Diff          bool
		FlushCache    bool
		ForceHandlers bool
		Forks         int
		Inventory     string
		ListHosts     bool
		ListTags      bool
		ListTasks     bool
		SyntaxCheck   bool
		Verbose       bool
		Version       bool
		Timeout       int
		Become        bool
		Requirements  *AnsiblePlaybookRequirements
	}
	test := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a AnsiblePlaybookParameters",
			fields: fields{
				Playbooks:     []string{"playbook.yml"},
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
			wantErr: false,
		},
		{
			desc: "Validating a AnsiblePlaybookParameters with empty playbooks",
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
			desc: "Validating a AnsiblePlaybookParameters with empty inventory",
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
		{
			desc: "Testing validate a AnsiblePlaybookParameters with forks less than 1",
			fields: fields{
				Playbooks: []string{"playbook.yml"},
				Inventory: "inventory",
				Forks:     -1,
				Timeout:   30,
			},
			wantErr: true,
		},
		{
			desc: "Testing validate a AnsiblePlaybookParameters with timeout less than 1",
			fields: fields{
				Playbooks: []string{"playbook.yml"},
				Inventory: "inventory",
				Forks:     5,
				Timeout:   -1,
			},
			wantErr: true,
		},
		{
			desc: "Testing validate a AnsiblePlaybookParameters with roles requirement timeout less than 1",
			fields: fields{
				Playbooks: []string{"playbook.yml"},
				Inventory: "inventory",
				Forks:     5,
				Timeout:   30,
				Requirements: &AnsiblePlaybookRequirements{
					Roles: &AnsiblePlaybookRoleRequirements{
						Roles:   []string{"roles"},
						Timeout: -1,
					},
				},
			},
			wantErr: true,
		},
		{
			desc: "Testing validate a AnsiblePlaybookParameters with collections requirement timeout less than 1",
			fields: fields{
				Playbooks: []string{"playbook.yml"},
				Inventory: "inventory",
				Forks:     5,
				Timeout:   30,
				Requirements: &AnsiblePlaybookRequirements{
					Collections: &AnsiblePlaybookCollectionRequirements{
						Collections: []string{"collections"},
						Timeout:     -1,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range test {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

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
