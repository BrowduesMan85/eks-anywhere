package executables

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	"github.com/aws/eks-anywhere/pkg/logger"
	"github.com/aws/eks-anywhere/pkg/templater"
)

//go:embed config/cmk.ini
var cmkConfigTemplate string

const (
	cmkConfigFileName = "cmk_tmp.ini"
	Shared            = "Shared"
)

// Cmk this type will be used once the CloudStack provider is added to the repository
type Cmk struct {
	writer     filewriter.FileWriter
	executable Executable
	config     CmkExecConfig
}

func (c *Cmk) ValidateTemplatePresent(ctx context.Context, domainId string, zoneId string, account string, template v1alpha1.CloudStackResourceRef) error {
	command := newCmkCommand("list templates")
	applyCmkArgs(&command, appendArgs("templatefilter=all"), appendArgs("listall=true"))
	if template.Type == v1alpha1.Id {
		applyCmkArgs(&command, withCloudStackId(template.Value))
	} else {
		applyCmkArgs(&command, withCloudStackName(template.Value))
	}

	applyCmkArgs(&command, withCloudStackZoneId(zoneId))
	if len(domainId) > 0 {
		applyCmkArgs(&command, withCloudStackDomainId(domainId))
		if len(account) > 0 {
			applyCmkArgs(&command, withCloudStackAccount(account))
		}
	}
	result, err := c.exec(ctx, command...)
	if err != nil {
		return fmt.Errorf("error getting templates info: %v", err)
	}
	if result.Len() == 0 {
		return fmt.Errorf("template %s not found", template)
	}

	response := struct {
		CmkTemplates []cmkTemplate `json:"template"`
	}{}
	if err = json.Unmarshal(result.Bytes(), &response); err != nil {
		return fmt.Errorf("failed to parse response into json: %v", err)
	}
	templates := response.CmkTemplates
	if len(templates) > 1 {
		return fmt.Errorf("duplicate templates %s found", template)
	} else if len(templates) == 0 {
		return fmt.Errorf("template %s not found", template)
	}
	return nil
}

func (c *Cmk) ValidateServiceOfferingPresent(ctx context.Context, zoneId string, serviceOffering v1alpha1.CloudStackResourceRef) error {
	command := newCmkCommand("list serviceofferings")
	if serviceOffering.Type == v1alpha1.Id {
		applyCmkArgs(&command, withCloudStackId(serviceOffering.Value))
	} else {
		applyCmkArgs(&command, withCloudStackName(serviceOffering.Value))
	}
	applyCmkArgs(&command, withCloudStackZoneId(zoneId))
	result, err := c.exec(ctx, command...)
	if err != nil {
		return fmt.Errorf("error getting service offerings info: %v", err)
	}
	if result.Len() == 0 {
		return fmt.Errorf("service offering %s not found", serviceOffering)
	}

	response := struct {
		CmkServiceOfferings []cmkServiceOffering `json:"serviceoffering"`
	}{}
	if err = json.Unmarshal(result.Bytes(), &response); err != nil {
		return fmt.Errorf("failed to parse response into json: %v", err)
	}
	offerings := response.CmkServiceOfferings
	if len(offerings) > 1 {
		return fmt.Errorf("duplicate service offering %s found", serviceOffering)
	} else if len(offerings) == 0 {
		return fmt.Errorf("service offering %s not found", serviceOffering)
	}

	return nil
}

func (c *Cmk) ValidateAffinityGroupsPresent(ctx context.Context, domainId string, account string, affinityGroupIds []string) error {
	for _, affinityGroupId := range affinityGroupIds {
		command := newCmkCommand("list affinitygroups")
		applyCmkArgs(&command, withCloudStackId(affinityGroupId))
		// account must be specified with a domainId
		// domainId can be specified without account
		if len(domainId) > 0 {
			applyCmkArgs(&command, withCloudStackDomainId(domainId))
			if len(account) > 0 {
				applyCmkArgs(&command, withCloudStackAccount(account))
			}
		}

		result, err := c.exec(ctx, command...)
		if err != nil {
			return fmt.Errorf("error getting affinity group info: %v", err)
		}
		if result.Len() == 0 {
			return fmt.Errorf(fmt.Sprintf("affinity group %s not found", affinityGroupId))
		}

		response := struct {
			CmkAffinityGroups []cmkAffinityGroup `json:"affinitygroup"`
		}{}
		if err = json.Unmarshal(result.Bytes(), &response); err != nil {
			return fmt.Errorf("failed to parse response into json: %v", err)
		}
		affinityGroup := response.CmkAffinityGroups
		if len(affinityGroup) > 1 {
			return fmt.Errorf("duplicate affinity group %s found", affinityGroupId)
		} else if len(affinityGroup) == 0 {
			return fmt.Errorf("affinity group %s not found", affinityGroupId)
		}
	}
	return nil
}

func (c *Cmk) ValidateZonesPresent(ctx context.Context, zones []v1alpha1.CloudStackZoneRef) ([]v1alpha1.CloudStackResourceIdentifier, error) {
	var zoneIdentifiers []v1alpha1.CloudStackResourceIdentifier
	for _, z := range zones {
		command := newCmkCommand("list zones")
		zone := z.Zone
		if zone.Type == v1alpha1.Id {
			applyCmkArgs(&command, withCloudStackId(zone.Value))
		} else {
			applyCmkArgs(&command, withCloudStackName(zone.Value))
		}
		result, err := c.exec(ctx, command...)
		if err != nil {
			return nil, fmt.Errorf("error getting zones info: %v", err)
		}
		if result.Len() == 0 {
			return nil, fmt.Errorf("zone %s not found", zone.Value)
		}

		response := struct {
			CmkZones []cmkZone `json:"zone"`
		}{}
		if err = json.Unmarshal(result.Bytes(), &response); err != nil {
			return nil, fmt.Errorf("failed to parse response into json: %v", err)
		}
		cmkZones := response.CmkZones
		if len(cmkZones) > 1 {
			return nil, fmt.Errorf("duplicate zone %s found", zone.Value)
		} else if len(zones) == 0 {
			return nil, fmt.Errorf("zone %s not found", zone.Value)
		} else {
			zoneIdentifiers = append(zoneIdentifiers, v1alpha1.CloudStackResourceIdentifier{Name: cmkZones[0].Name, Id: cmkZones[0].Id})
		}
	}
	return zoneIdentifiers, nil
}

func (c *Cmk) ValidateDomainPresent(ctx context.Context, domain string) (v1alpha1.CloudStackResourceIdentifier, error) {
	domainIdentifier := v1alpha1.CloudStackResourceIdentifier{Name: domain, Id: ""}
	command := newCmkCommand("list domains")
	applyCmkArgs(&command, withCloudStackName(domain))
	result, err := c.exec(ctx, command...)
	if err != nil {
		return domainIdentifier, fmt.Errorf("error getting domain info: %v", err)
	}
	if result.Len() == 0 {
		return domainIdentifier, fmt.Errorf("domain %s not found", domain)
	}

	response := struct {
		CmkDomains []cmkDomain `json:"domain"`
	}{}
	if err = json.Unmarshal(result.Bytes(), &response); err != nil {
		return domainIdentifier, fmt.Errorf("failed to parse response into json: %v", err)
	}
	domains := response.CmkDomains
	if len(domains) > 1 {
		return domainIdentifier, fmt.Errorf("duplicate domain %s found", domain)
	} else if len(domains) == 0 {
		return domainIdentifier, fmt.Errorf("domain %s not found", domain)
	}

	domainIdentifier.Id = domains[0].Id
	domainIdentifier.Name = domains[0].Name

	return domainIdentifier, nil
}

func (c *Cmk) ValidateNetworkPresent(ctx context.Context, domainId string, zoneRef v1alpha1.CloudStackZoneRef, zones []v1alpha1.CloudStackResourceIdentifier, account string, multipleZone bool) error {
	command := newCmkCommand("list networks")
	if zoneRef.Network.Type == v1alpha1.Id {
		applyCmkArgs(&command, withCloudStackId(zoneRef.Network.Value))
	}
	if multipleZone {
		applyCmkArgs(&command, withCloudStackNetworkType(Shared))
	}
	// account must be specified within a domainId
	// domainId can be specified without account
	if len(domainId) > 0 {
		applyCmkArgs(&command, withCloudStackDomainId(domainId))
		if len(account) > 0 {
			applyCmkArgs(&command, withCloudStackAccount(account))
		}
	}
	var zoneId string
	var err error
	if zoneRef.Zone.Type == v1alpha1.Id {
		zoneId = zoneRef.Zone.Value
	} else {
		zoneId, err = getZoneIdByName(zones, zoneRef.Zone.Value)
		if err != nil {
			return fmt.Errorf("error getting zone id by name: %v", err)
		}
	}
	applyCmkArgs(&command, withCloudStackZoneId(zoneId))
	result, err := c.exec(ctx, command...)
	if err != nil {
		return fmt.Errorf("error getting network info: %v", err)
	}
	if result.Len() == 0 {
		if multipleZone {
			return fmt.Errorf("%s network %s not found in zoneRef %s", Shared, zoneRef.Network.Value, zoneRef.Zone.Value)
		} else {
			return fmt.Errorf("network %s not found in zoneRef %s", zoneRef.Network.Value, zoneRef.Zone.Value)
		}
	}

	response := struct {
		CmkNetworks []cmkNetwork `json:"network"`
	}{}
	if err = json.Unmarshal(result.Bytes(), &response); err != nil {
		return fmt.Errorf("failed to parse response into json: %v", err)
	}
	networks := response.CmkNetworks

	// filter by network name -- cmk does not support name= filter
	if zoneRef.Network.Type == v1alpha1.Name {
		networks = []cmkNetwork{}
		for _, net := range response.CmkNetworks {
			if net.Name == zoneRef.Network.Value {
				networks = append(networks, net)
			}
		}
	}

	if len(networks) > 1 {
		return fmt.Errorf("duplicate network %s found", zoneRef.Network.Value)
	} else if len(networks) == 0 {
		if multipleZone {
			return fmt.Errorf("%s network %s not found in zoneRef %s", Shared, zoneRef.Network.Value, zoneRef.Zone.Value)
		} else {
			return fmt.Errorf("network %s not found in zoneRef %s", zoneRef.Network.Value, zoneRef.Zone.Value)
		}
	}
	return nil
}

func getZoneIdByName(zones []v1alpha1.CloudStackResourceIdentifier, zoneName string) (string, error) {
	for _, zoneIdentifier := range zones {
		if zoneName == zoneIdentifier.Name {
			return zoneIdentifier.Id, nil
		}
	}
	return "", fmt.Errorf("zoneId not found for zone %s", zoneName)
}

func (c *Cmk) ValidateAccountPresent(ctx context.Context, account string, domainId string) error {
	command := newCmkCommand("list accounts")
	applyCmkArgs(&command, withCloudStackName(account), withCloudStackDomainId(domainId))
	result, err := c.exec(ctx, command...)
	if err != nil {
		return fmt.Errorf("error getting accounts info: %v", err)
	}
	if result.Len() == 0 {
		return fmt.Errorf("account %s not found", account)
	}

	response := struct {
		CmkAccounts []cmkAccount `json:"account"`
	}{}
	if err = json.Unmarshal(result.Bytes(), &response); err != nil {
		return fmt.Errorf("failed to parse response into json: %v", err)
	}
	accounts := response.CmkAccounts
	if len(accounts) > 1 {
		return fmt.Errorf("duplicate account %s found", account)
	} else if len(accounts) == 0 {
		return fmt.Errorf("account %s not found", account)
	}
	return nil
}

func NewCmk(executable Executable, writer filewriter.FileWriter, config CmkExecConfig) *Cmk {
	return &Cmk{
		writer:     writer,
		executable: executable,
		config:     config,
	}
}

// ValidateCloudStackConnection Calls `cmk sync` to ensure that the endpoint and credentials + domain are valid
func (c *Cmk) ValidateCloudStackConnection(ctx context.Context) error {
	command := newCmkCommand("sync")
	buffer, err := c.exec(ctx, command...)
	if err != nil {
		return fmt.Errorf("error validating cloudstack connection for cmk config %s: %v", buffer.String(), err)
	}
	logger.MarkPass("Connected to CloudStack server")
	return nil
}

func (c *Cmk) exec(ctx context.Context, args ...string) (stdout bytes.Buffer, err error) {
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed get environment map: %v", err)
	}
	configFile, err := c.buildCmkConfigFile()
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed cmk validations: %v", err)
	}
	argsWithConfigFile := append([]string{"-c", configFile}, args...)

	return c.executable.Execute(ctx, argsWithConfigFile...)
}

// TODO: Add support for passing in domain from Deployment Config Spec
type CmkExecConfig struct {
	CloudStackApiKey        string // Api Key for CloudMonkey to access CloudStack Cluster
	CloudStackSecretKey     string // Secret Key for CloudMonkey to access CloudStack Cluster
	CloudStackManagementUrl string // Management Endpoint Url for CloudMonkey to access CloudStack Cluster
	CloudMonkeyVerifyCert   bool   // boolean indicating if CloudMonkey should verify the cert presented by the CloudStack Management Server
}

func (c *Cmk) buildCmkConfigFile() (configFile string, err error) {
	t := templater.New(c.writer)
	writtenFileName, err := t.WriteToFile(cmkConfigTemplate, c.config, cmkConfigFileName)
	if err != nil {
		return "", fmt.Errorf("error creating file for cmk config: %v", err)
	}
	configFile, err = filepath.Abs(writtenFileName)
	if err != nil {
		return "", fmt.Errorf("failed to generate absolute filepath for generated config file at %s", writtenFileName)
	}

	return configFile, nil
}

type cmkTemplate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Zonename string `json:"zonename"`
}

type cmkServiceOffering struct {
	CpuNumber int    `json:"cpunumber"`
	CpuSpeed  int    `json:"cpuspeed"`
	Memory    int    `json:"memory"`
	Id        string `json:"id"`
	Name      string `json:"name"`
}

type cmkAffinityGroup struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type cmkZone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type cmkNetwork struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type cmkDomain struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type cmkAccount struct {
	RoleType string `json:"roletype"`
	Domain   string `json:"domain"`
	Id       string `json:"id"`
	Name     string `json:"name"`
}