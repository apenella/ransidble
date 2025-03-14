package mapper

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
)

// AnsiblePlaybookParametersMapper is responsible for mapping ansible playbook parameters
type AnsiblePlaybookParametersMapper struct{}

// NewAnsiblePlaybookParametersMapper creates a new ansible playbook parameters mapper
func NewAnsiblePlaybookParametersMapper() *AnsiblePlaybookParametersMapper {
	return &AnsiblePlaybookParametersMapper{}
}

// ToAnsiblePlaybookParametersEntity maps a request.AnsiblePlaybookParameters to a entity.AnsiblePlaybookParameters
func (m *AnsiblePlaybookParametersMapper) ToAnsiblePlaybookParametersEntity(parameters *request.AnsiblePlaybookParameters) *entity.AnsiblePlaybookParameters {

	if parameters == nil {
		return &entity.AnsiblePlaybookParameters{}
	}

	return &entity.AnsiblePlaybookParameters{
		Playbooks:         append([]string{}, parameters.Playbooks...),
		Check:             parameters.Check,
		Diff:              parameters.Diff,
		Requirements:      m.toAnsiblePLaybookParametersRequirementsEntity(parameters.Requirements),
		ExtraVars:         m.toAnsiblePlaybookParametersExtraVarsEntity(parameters.ExtraVars),
		ExtraVarsFile:     append([]string{}, parameters.ExtraVarsFile...),
		FlushCache:        parameters.FlushCache,
		ForceHandlers:     parameters.ForceHandlers,
		Forks:             parameters.Forks,
		Inventory:         parameters.Inventory,
		Limit:             parameters.Limit,
		ListHosts:         parameters.ListHosts,
		ListTags:          parameters.ListTags,
		ListTasks:         parameters.ListTasks,
		SkipTags:          parameters.SkipTags,
		StartAtTask:       parameters.StartAtTask,
		SyntaxCheck:       parameters.SyntaxCheck,
		Tags:              parameters.Tags,
		VaultID:           parameters.VaultID,
		VaultPasswordFile: parameters.VaultPasswordFile,
		Verbose:           parameters.Verbose,
		Version:           parameters.Version,
		Connection:        parameters.Connection,
		SCPExtraArgs:      parameters.SCPExtraArgs,
		SFTPExtraArgs:     parameters.SFTPExtraArgs,
		SSHCommonArgs:     parameters.SSHCommonArgs,
		SSHExtraArgs:      parameters.SSHExtraArgs,
		Timeout:           parameters.Timeout,
		User:              parameters.User,
		Become:            parameters.Become,
		BecomeMethod:      parameters.BecomeMethod,
		BecomeUser:        parameters.BecomeUser,
	}
}

// ToAnsiblePLaybookParametersRequirementsEntity maps a request.AnsiblePlaybookParametersDependencies to a entity.AnsiblePlaybookParametersDependencies
func (m *AnsiblePlaybookParametersMapper) toAnsiblePLaybookParametersRequirementsEntity(dependencies *request.AnsiblePlaybookRequirements) *entity.AnsiblePlaybookRequirements {

	if dependencies == nil {
		return &entity.AnsiblePlaybookRequirements{}
	}

	return &entity.AnsiblePlaybookRequirements{
		Roles:       m.toAnsiblePLaybookParametersRolesRequirementsEntity(dependencies.Roles),
		Collections: m.toAnsiblePLaybookParametersCollectionsRequirementsEntity(dependencies.Collections),
	}
}

// toAnsiblePLaybookParametersRolesRequirementsEntity
func (m *AnsiblePlaybookParametersMapper) toAnsiblePLaybookParametersRolesRequirementsEntity(parameters *request.AnsiblePlaybookRoleRequirements) *entity.AnsiblePlaybookRoleRequirements {

	if parameters == nil {
		return &entity.AnsiblePlaybookRoleRequirements{}
	}

	return &entity.AnsiblePlaybookRoleRequirements{
		Roles:        append([]string{}, parameters.Roles...),
		APIKey:       parameters.APIKey,
		IgnoreErrors: parameters.IgnoreErrors,
		NoDeps:       parameters.NoDeps,
		RoleFile:     parameters.RoleFile,
		Server:       parameters.Server,
		Timeout:      parameters.Timeout,
		Token:        parameters.Token,
		Verbose:      parameters.Verbose,
	}
}

// toAnsiblePLaybookParametersCollectionsRequirementsEntity
func (m *AnsiblePlaybookParametersMapper) toAnsiblePLaybookParametersCollectionsRequirementsEntity(paremeters *request.AnsiblePlaybookCollectionRequirements) *entity.AnsiblePlaybookCollectionRequirements {

	if paremeters == nil {
		return &entity.AnsiblePlaybookCollectionRequirements{}
	}

	return &entity.AnsiblePlaybookCollectionRequirements{
		Collections:      append([]string{}, paremeters.Collections...),
		APIKey:           paremeters.APIKey,
		ForceWithDeps:    paremeters.ForceWithDeps,
		Pre:              paremeters.Pre,
		Timeout:          paremeters.Timeout,
		Token:            paremeters.Token,
		IgnoreErrors:     paremeters.IgnoreErrors,
		RequirementsFile: paremeters.RequirementsFile,
		Server:           paremeters.Server,
		Verbose:          paremeters.Verbose,
	}
}

// ToAnsiblePlaybookParametersExtraVarsEntity copies the content of a request.AnsiblePlaybookParametersExtraVars to a entity.AnsiblePlaybookParametersExtraVars
func (m *AnsiblePlaybookParametersMapper) toAnsiblePlaybookParametersExtraVarsEntity(extraVars map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range extraVars {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[key] = m.copyMap(v)
		case []interface{}:
			copy[key] = m.copySlice(v)
		default:
			copy[key] = value
		}
	}

	return copy
}

// copyMap copies the content of a map[string]interface{} to a new map[string]interface{}
func (m *AnsiblePlaybookParametersMapper) copyMap(original map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range original {

		switch v := value.(type) {
		case map[string]interface{}:
			copy[key] = m.copyMap(v)
		case []interface{}:
			copy[key] = m.copySlice(v)
		default:
			copy[key] = value
		}
	}
	return copy
}

// copySlice copies the content of a []interface{} to a new []interface{}
func (m *AnsiblePlaybookParametersMapper) copySlice(original []interface{}) []interface{} {
	copy := make([]interface{}, len(original))
	for i, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[i] = m.copyMap(v)
		case []interface{}:
			copy[i] = m.copySlice(v)
		default:
			copy[i] = value
		}
	}
	return copy
}
