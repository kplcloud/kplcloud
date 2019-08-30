/**
 * @Time : 2019-06-25 17:52
 * @Author : solacowa@gmail.com
 * @File : types
 * @Software: GoLand
 */

package jenkins

import "encoding/xml"

type DefaultParameterValue struct {
	Class string `json:"_class"`
	Value string `json:"value"`
}

type HudsonModelTextParameterDefinition struct {
	Text         string `xml:",chardata"`
	Name         string `xml:"name"`
	Description  string `xml:"description"`
	DefaultValue string `xml:"defaultValue"`
	Trim         string `xml:"trim"`
}

type ParameterDefinitions struct {
	Class                              string                               `json:"_class"`
	DefaultParameterValue              DefaultParameterValue                `json:"defaultParameterValue"`
	Description                        string                               `json:"description"`
	Name                               string                               `json:"name"`
	Type                               string                               `json:"type"`
	HudsonModelTextParameterDefinition []HudsonModelTextParameterDefinition `xml:"hudson.model.TextParameterDefinition"`
}

type Action struct {
	Class                string                 `json:"_class"`
	ParameterDefinitions []ParameterDefinitions `json:"parameterDefinitions"`
}

type Build struct {
	Id     string `json:"id"`
	Class  string `json:"_class"`
	Number int    `json:"number"`
	URL    string `json:"url"`

	FullDisplayName string `json:"fullDisplayName"`
	Description     string `json:"description"`

	Timestamp         int `json:"timestamp"`
	Duration          int `json:"duration"`
	EstimatedDuration int `json:"estimatedDuration"`

	Building bool   `json:"building"`
	KeepLog  bool   `json:"keepLog"`
	Result   string `json:"result"`

	Artifacts []Artifact `json:"artifacts"`
	Actions   []Action   `json:"actions"`

	//	ChangeSet ScmChangeSet `json:"changeSet"`
}

type HealthReport struct {
	Description   string `json:"description"`
	IconClassName string `json:"iconClassName"`
	IconURL       string `json:"iconUrl"`
	Score         int    `json:"score"`
}

type Property struct {
	Class string `json:"_class"`
}

type GitUserRemoteConfig struct {
	Url           string `xml:"url"`
	CredentialsId string `xml:"credentialsId"`
}

type SubmoduleCfg struct {
	Text  string `xml:",chardata"`
	Class string `xml:"class,attr"`
}

type HudsonPluginsGitUserRemoteConfig struct {
	Text          string `xml:",chardata"`
	URL           string `xml:"url"`
	CredentialsId string `xml:"credentialsId"`
}

type UserRemoteConfigs struct {
	GitUserRemoteConfig GitUserRemoteConfig `xml:"hudson.plugins.git.UserRemoteConfig"`
	Text                string              `xml:",chardata"`
	//HudsonPluginsGitUserRemoteConfig GitUserRemoteConfig `xml:"hudson.plugins.git.UserRemoteConfig"`
}

type BranchSpec struct {
	Name string `xml:"name"`
}

type HudsonPluginsGitBranchSpec struct {
	Text string `xml:",chardata"`
	Name string `xml:"name"`
}

type Branches struct {
	BranchSpec HudsonPluginsGitBranchSpec `xml:"hudson.plugins.git.BranchSpec"`
	//HudsonPluginsGitBranchSpec HudsonPluginsGitBranchSpec `xml:"hudson.plugins.git.BranchSpec"`
}

type Scm struct {
	ConfigVersion                     string            `xml:"configVersion"`
	SubmoduleCfg                      SubmoduleCfg      `xml:"submoduleCfg"`
	DoGenerateSubmoduleConfigurations string            `xml:"doGenerateSubmoduleConfigurations"`
	UserRemoteConfigs                 UserRemoteConfigs `xml:"userRemoteConfigs"`
	Branches                          Branches          `xml:"branches"`
	Class                             string            `xml:"class,attr"`
	Plugin                            string            `xml:"plugin,attr"`
	Text                              string            `xml:",chardata"`
	Extensions                        string            `xml:"extensions"`
}

type View struct {
	Name            string   `json:"name"`
	URL             string   `json:"url"`
	XMLName         xml.Name `xml:"hudson.model.MyView"`
	Text            string   `xml:",chardata"`
	FilterExecutors string   `xml:"filterExecutors"`
	FilterQueue     string   `xml:"filterQueue"`
	Properties      struct {
		Text  string `xml:",chardata"`
		Class string `xml:"class,attr"`
	} `xml:"properties"`
}

type Job struct {
	Actions               []Action       `json:"actions" xml:"actions"`
	Buildable             bool           `json:"buildable" xml:"buildable"`
	Builds                []Build        `json:"builds" xml:"builds"`
	Color                 string         `json:"color" xml:"color"`
	ConcurrentBuild       bool           `json:"concurrentBuild" xml:"concurrentBuild"`
	Description           string         `json:"description" xml:"description"`
	DisplayName           string         `json:"displayName" xml:"displayName"`
	DisplayNameOrNull     interface{}    `json:"displayNameOrNull" xml:"displayNameOrNull"`
	DownstreamProjects    []interface{}  `json:"downstreamProjects" xml:"downstreamProjects"`
	FirstBuild            Build          `json:"firstBuild" xml:"firstBuild"`
	FullDisplayName       string         `json:"fullDisplayName" xml:"fullDisplayName"`
	FullName              string         `json:"fullName" xml:"fullName"`
	HealthReport          []HealthReport `json:"healthReport" xml:"healthReport"`
	InQueue               bool           `json:"inQueue" xml:"inQueue"`
	KeepDependencies      bool           `json:"keepDependencies" xml:"keepDependencies"`
	LabelExpression       interface{}    `json:"labelExpression" xml:"labelExpression"`
	LastBuild             Build          `json:"lastBuild" xml:"lastBuild"`
	LastCompletedBuild    Build          `json:"lastCompletedBuild" xml:"lastCompletedBuild"`
	LastFailedBuild       interface{}    `json:"lastFailedBuild" xml:"lastFailedBuild"`
	LastStableBuild       Build          `json:"lastStableBuild" xml:"lastStableBuild"`
	LastSuccessfulBuild   Build          `json:"lastSuccessfulBuild" xml:"lastSuccessfulBuild"`
	LastUnstableBuild     interface{}    `json:"lastUnstableBuild" xml:"lastUnstableBuild"`
	LastUnsuccessfulBuild interface{}    `json:"lastUnsuccessfulBuild" xml:"lastUnsuccessfulBuild"`
	Name                  string         `json:"name" xml:"name"`
	NextBuildNumber       int            `json:"nextBuildNumber" xml:"nextBuildNumber"`
	Property              []Property     `json:"property" xml:"property"`
	QueueItem             interface{}    `json:"queueItem" xml:"queueItem"`
	Scm                   Scm            `json:"scm" xml:"scm"`
	UpstreamProjects      []interface{}  `json:"upstreamProjects" xml:"upstreamProjects"`
	URL                   string         `json:"url" xml:"url"`
}

// mavenjob
type JenkinsPluginsMaveninfoConfigMavenInfoJobConfig struct {
	Text                string `xml:",chardata"`
	Plugin              string `xml:"plugin,attr"`
	MainModulePattern   string `xml:"mainModulePattern"`
	DependenciesPattern string `xml:"dependenciesPattern"`
	AssignName          string `xml:"assignName"`
	NameTemplate        string `xml:"nameTemplate"`
	AssignDescription   string `xml:"assignDescription"`
	DescriptionTemplate string `xml:"descriptionTemplate"`
}

type HudsonModelParametersDefinitionProperty struct {
	Text                 string               `xml:",chardata"`
	ParameterDefinitions ParameterDefinitions `xml:"parameterDefinitions"`
}

type Properties struct {
	Text                                            string                                          `xml:",chardata"`
	JenkinsPluginsMaveninfoConfigMavenInfoJobConfig JenkinsPluginsMaveninfoConfigMavenInfoJobConfig `xml:"jenkins.plugins.maveninfo.config.MavenInfoJobConfig"`
	HudsonModelParametersDefinitionProperty         HudsonModelParametersDefinitionProperty         `xml:"hudson.model.ParametersDefinitionProperty"`
}

type Settings struct {
	Text  string `xml:",chardata"`
	Class string `xml:"class,attr"`
}

type HudsonTasksShell struct {
	Text    string `xml:",chardata"`
	Command string `xml:"command"`
}

type Postbuilders struct {
	Text             string           `xml:",chardata"`
	HudsonTasksShell HudsonTasksShell `xml:"hudson.tasks.Shell"`
}

type RunPostStepsIfResult struct {
	Text          string `xml:",chardata"`
	Name          string `xml:"name"`
	Ordinal       string `xml:"ordinal"`
	Color         string `xml:"color"`
	CompleteBuild string `xml:"completeBuild"`
}

type MavenJobItem struct {
	XMLName                          xml.Name             `xml:"project"`
	Text                             string               `xml:",chardata"`
	Plugin                           string               `xml:"plugin,attr"`
	Description                      string               `xml:"description"`
	KeepDependencies                 string               `xml:"keepDependencies"`
	Properties                       Properties           `xml:"properties"`
	Scm                              Scm                  `xml:"scm"`
	CanRoam                          string               `xml:"canRoam"`
	Disabled                         string               `xml:"disabled"`
	BlockBuildWhenDownstreamBuilding string               `xml:"blockBuildWhenDownstreamBuilding"`
	BlockBuildWhenUpstreamBuilding   string               `xml:"blockBuildWhenUpstreamBuilding"`
	Triggers                         string               `xml:"triggers"`
	ConcurrentBuild                  string               `xml:"concurrentBuild"`
	RootPOM                          string               `xml:"rootPOM"`
	Goals                            string               `xml:"goals"`
	AggregatorStyleBuild             string               `xml:"aggregatorStyleBuild"`
	IncrementalBuild                 string               `xml:"incrementalBuild"`
	IgnoreUpstremChanges             string               `xml:"ignoreUpstremChanges"`
	IgnoreUnsuccessfulUpstreams      string               `xml:"ignoreUnsuccessfulUpstreams"`
	ArchivingDisabled                string               `xml:"archivingDisabled"`
	SiteArchivingDisabled            string               `xml:"siteArchivingDisabled"`
	FingerprintingDisabled           string               `xml:"fingerprintingDisabled"`
	ResolveDependencies              string               `xml:"resolveDependencies"`
	ProcessPlugins                   string               `xml:"processPlugins"`
	MavenValidationLevel             string               `xml:"mavenValidationLevel"`
	RunHeadless                      string               `xml:"runHeadless"`
	DisableTriggerDownstreamProjects string               `xml:"disableTriggerDownstreamProjects"`
	BlockTriggerWhenBuilding         string               `xml:"blockTriggerWhenBuilding"`
	Settings                         Settings             `xml:"settings"`
	GlobalSettings                   Settings             `xml:"globalSettings"`
	Reporters                        string               `xml:"reporters"`
	Publishers                       string               `xml:"publishers"`
	BuildWrappers                    string               `xml:"buildWrappers"`
	Prebuilders                      string               `xml:"prebuilders"`
	Postbuilders                     Postbuilders         `xml:"postbuilders"`
	Builders                         Postbuilders         `xml:"builders"`
	RunPostStepsIfResult             RunPostStepsIfResult `xml:"runPostStepsIfResult"`
}

type ViewProperties struct {
	Text  string `xml:",chardata"`
	Class string `xml:"class,attr"`
}

type ListView struct {
	Name            string         `xml:"name"`
	XMLName         xml.Name       `xml:"hudson.model.MyView"`
	Text            string         `xml:",chardata"`
	FilterExecutors string         `xml:"filterExecutors"`
	FilterQueue     string         `xml:"filterQueue"`
	Properties      ViewProperties `xml:"properties"`
}

type Queue struct {
}

type Artifact struct {
	RelativePath string `xml:"relative_path"`
}

type ComputerObject struct {
}

type Computer struct {
}
