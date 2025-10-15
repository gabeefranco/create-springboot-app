package springinitializr

type AppPackaging string
type AppLanguage string
type AppProjectType string
type DependencyCategory string

const Gradle = AppProjectType("gradle-project")
const Maven = AppProjectType("maven-project")

const Java = AppLanguage("java")
const Kotlin = AppLanguage("kotlin")
const Groovy = AppLanguage("groovy")

const Jar = AppPackaging("jar")
const War = AppPackaging("war")

const Core = DependencyCategory("core")
const Other = DependencyCategory("other")

type Dependency struct {
	ID          string
	Name        string
	Description string
	Category    DependencyCategory
}

type InitializrOptions struct {
	JavaVersion  string
	Packaging    AppPackaging
	Language     AppLanguage
	ProjectType  AppProjectType
	ArtifactID   string
	GroupID      string
	Name         string
	Description  string
	BootVersion  string
	Version      string
	PackageName  string
	Dependencies []Dependency
}

type Capabilities struct {
	JavaVersions   []string
	PackagingTypes []AppPackaging
	ProjectTypes   []AppProjectType
	Languages      []AppLanguage
	BootVersions   []string
	Dependencies   []Dependency
}

type simpleCapabilityDto struct {
	Values []struct {
		ID string `json:"id"`
	} `json:"values"`
}

type dependencyDto struct {
	Values []struct {
		Name   string `json:"name"`
		Values []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"values"`
	} `json:"values"`
}

type capabilitiesDto struct {
	JavaVersions   simpleCapabilityDto `json:"javaVersion"`
	PackagingTypes simpleCapabilityDto `json:"packaging"`
	ProjectTypes   simpleCapabilityDto `json:"type"`
	Languages      simpleCapabilityDto `json:"language"`
	BootVersions   simpleCapabilityDto `json:"bootVersion"`
	Dependencies   dependencyDto       `json:"dependencies"`
}

const DefaultBaseURL = "https://start.spring.io"

type Instance struct {
	BaseURL string
}

func DefaultInstance() Instance {
	return Instance{
		BaseURL: DefaultBaseURL,
	}
}
