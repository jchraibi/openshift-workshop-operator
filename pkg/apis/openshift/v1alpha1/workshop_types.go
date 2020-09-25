package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WorkshopSpec defines the desired state of Workshop
// +k8s:openapi-gen=true
type WorkshopSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	User           UserSpec           `json:"user"`
	Source         SourceSpec         `json:"source"`
	Infrastructure InfrastructureSpec `json:"infrastructure"`
}

type UserSpec struct {
	Number   int    `json:"number"`
	Password string `json:"password"`
}

type SourceSpec struct {
	GitURL    string `json:"gitURL"`
	GitBranch string `json:"gitBranch"`
}

type InfrastructureSpec struct {
	ArgoCD             ArgoCDSpec             `json:"argocd"`
	Bookbag            BookbagSpec            `json:"bookbag"`
	CertManager        CertManagerSpec        `json:"certManager"`
	CodeReadyWorkspace CodeReadyWorkspaceSpec `json:"codeReadyWorkspace"`
	Etherpad           EtherpadSpec           `json:"etherpad"`
	Gitea              GiteaSpec              `json:"gitea"`
	Guide              GuideSpec              `json:"guide"`
	IstioWorkspace     IstioWorkspaceSpec     `json:"istioWorkspace"`
	Nexus              NexusSpec              `json:"nexus"`
	Pipeline           PipelineSpec           `json:"pipeline"`
	Project            ProjectSpec            `json:"project"`
	ServiceMesh        ServiceMeshSpec        `json:"serviceMesh"`
	Serverless         ServerlessSpec         `json:"serverless"`
	Vault              VaultSpec              `json:"vault"`
}

type ArgoCDSpec struct {
	Enabled     bool            `json:"enabled"`
	OperatorHub OperatorHubSpec `json:"operatorHub"`
}

type BookbagSpec struct {
	Enabled bool      `json:"enabled"`
	Image   ImageSpec `json:"image"`
}

type CertManagerSpec struct {
	Enabled     bool            `json:"enabled"`
	OperatorHub OperatorHubSpec `json:"operatorHub"`
}

type EtherpadSpec struct {
	Enabled bool `json:"enabled"`
}

type GiteaSpec struct {
	Enabled bool      `json:"enabled"`
	Image   ImageSpec `json:"image"`
}

type NexusSpec struct {
	Enabled bool `json:"enabled"`
}

type PipelineSpec struct {
	Enabled     bool            `json:"enabled"`
	OperatorHub OperatorHubSpec `json:"operatorHub"`
}

type ProjectSpec struct {
	Enabled     bool   `json:"enabled"`
	StagingName string `json:"stagingName"`
}

type ServiceMeshSpec struct {
	Enabled                  bool            `json:"enabled"`
	ServiceMeshOperatorHub   OperatorHubSpec `json:"serviceMeshOperatorHub"`
	ElasticSearchOperatorHub OperatorHubSpec `json:"elasticSearchOperatorHub"`
	JaegerOperatorHub        OperatorHubSpec `json:"jaegerOperatorHub"`
	KialiOperatorHub         OperatorHubSpec `json:"kialiOperatorHub"`
}

type ServerlessSpec struct {
	Enabled     bool            `json:"enabled"`
	OperatorHub OperatorHubSpec `json:"operatorHub"`
}

type GuideSpec struct {
	Enabled                     bool   `json:"enabled"`
	GitRepositoryLabPath        string `json:"gitRepositoryLabPath"`
	GitRepositoryLabReference   string `json:"gitRepositoryLabReference"`
	GitRepositoryGuidePath      string `json:"gitRepositoryGuidePath"`
	GitRepositoryGuideReference string `json:"gitRepositoryGuideReference"`
	GitRepositoryGuideContext   string `json:"gitRepositoryGuideContext"`
	GitRepositoryGuideFile      string `json:"gitRepositoryGuideFile"`
}

type CodeReadyWorkspaceSpec struct {
	Enabled             bool            `json:"enabled"`
	OperatorHub         OperatorHubSpec `json:"operatorHub"`
	OpenshiftOAuth      bool            `json:"openshiftOAuth"`
	PluginRegistryImage ImageSpec       `json:"pluginRegistryImage"`
}

type IstioWorkspaceSpec struct {
	Enabled bool      `json:"enabled"`
	Image   ImageSpec `json:"image"`
}

type OperatorHubSpec struct {
	Channel               string `json:"channel"`
	ClusterServiceVersion string `json:"clusterServiceVersion"`
}

type ImageSpec struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type VaultSpec struct {
	Enabled            bool      `json:"enabled"`
	Image              ImageSpec `json:"image"`
	AgentInjectorImage ImageSpec `json:"agentInjectorImage"`
}

// WorkshopStatus defines the observed state of Workshop
// +k8s:openapi-gen=true
type WorkshopStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ArgoCD               string `json:"argocd"`
	Bookbag              string `json:"bookbag"`
	CertManager          string `json:"certManager"`
	CodeReadyWorkspace   string `json:"codeReadyWorkspace"`
	Etherpad             string `json:"etherpad"`
	Gitea                string `json:"gitea"`
	Guide                string `json:"guide"`
	IstioWorkspace       string `json:"istioWorkspace"`
	Nexus                string `json:"nexus"`
	Pipeline             string `json:"pipeline"`
	Project              string `json:"project"`
	ServiceMesh          string `json:"serviceMesh"`
	Serverless           string `json:"serverless"`
	UsernameDistribution string `json:"usernameDistribution"`
	Vault                string `json:"vault"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workshop is the Schema for the workshops API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Workshop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkshopSpec   `json:"spec,omitempty"`
	Status WorkshopStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkshopList contains a list of Workshop
type WorkshopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workshop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workshop{}, &WorkshopList{})
}
