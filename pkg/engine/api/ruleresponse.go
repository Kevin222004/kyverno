package api

import (
	"fmt"

	pssutils "github.com/kyverno/kyverno/pkg/pss/utils"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/pod-security-admission/api"
)

// PodSecurityChecks details about pod securty checks
type PodSecurityChecks struct {
	// Level is the pod security level
	Level api.Level
	// Version is the pod security version
	Version string
	// Checks contains check result details
	Checks []pssutils.PSSCheckResult
}

// RuleResponse details for each rule application
type RuleResponse struct {
	// name is the rule name specified in policy
	name string
	// ruleType is the rule type (Mutation,Generation,Validation) for Kyverno Policy
	ruleType RuleType
	// message is the message response from the rule application
	message string
	// status rule status
	status RuleStatus
	// stats contains rule statistics
	stats ExecutionStats
	// generatedResources is the list of resources generated by the generate rules of a policy
	generatedResources []*unstructured.Unstructured
	// patchedTarget is the patched resource for mutate.targets
	patchedTarget *unstructured.Unstructured
	// patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
	patchedTargetParentResourceGVR metav1.GroupVersionResource
	// patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
	patchedTargetSubresourceName string
	// podSecurityChecks contains pod security checks (only if this is a pod security rule)
	podSecurityChecks *PodSecurityChecks
	// exceptions are the exceptions applied (if any)
	exceptions []GenericException
	// vapbinding is the validatingadmissionpolicybinding (if any)
	vapBinding *admissionregistrationv1.ValidatingAdmissionPolicyBinding
	// mapbinding is the mutatingadmissionpolicybinding (if any)
	mapBinding *admissionregistrationv1alpha1.MutatingAdmissionPolicyBinding
	// emitWarning enable passing rule message as warning to api server warning header
	emitWarning bool
	// properties are the additional properties from the rule that will be added to the policy report result
	properties map[string]string
}

func NewRuleResponse(name string, ruleType RuleType, msg string, status RuleStatus, properties map[string]string) *RuleResponse {
	emitWarn := false
	if status == RuleStatusError || status == RuleStatusFail || status == RuleStatusWarn {
		emitWarn = true
	}
	return &RuleResponse{
		name:        name,
		ruleType:    ruleType,
		message:     msg,
		status:      status,
		emitWarning: emitWarn,
		properties:  properties,
	}
}

func RuleError(name string, ruleType RuleType, msg string, err error, properties map[string]string) *RuleResponse {
	if err != nil {
		return NewRuleResponse(name, ruleType, fmt.Sprintf("%s: %s", msg, err.Error()), RuleStatusError, properties)
	}
	return NewRuleResponse(name, ruleType, msg, RuleStatusError, properties)
}

func RuleSkip(name string, ruleType RuleType, msg string, properties map[string]string) *RuleResponse {
	return NewRuleResponse(name, ruleType, msg, RuleStatusSkip, properties)
}

func RuleWarn(name string, ruleType RuleType, msg string, properties map[string]string) *RuleResponse {
	return NewRuleResponse(name, ruleType, msg, RuleStatusWarn, properties)
}

func RulePass(name string, ruleType RuleType, msg string, properties map[string]string) *RuleResponse {
	return NewRuleResponse(name, ruleType, msg, RuleStatusPass, properties)
}

func RuleFail(name string, ruleType RuleType, msg string, properties map[string]string) *RuleResponse {
	return NewRuleResponse(name, ruleType, msg, RuleStatusFail, properties)
}

func (r RuleResponse) WithExceptions(exceptions []GenericException) *RuleResponse {
	r.exceptions = exceptions
	return &r
}

func (r RuleResponse) WithVAPBinding(binding *admissionregistrationv1.ValidatingAdmissionPolicyBinding) *RuleResponse {
	r.vapBinding = binding
	return &r
}

func (r RuleResponse) WithMAPBinding(binding *admissionregistrationv1alpha1.MutatingAdmissionPolicyBinding) *RuleResponse {
	r.mapBinding = binding
	return &r
}

func (r RuleResponse) WithPodSecurityChecks(checks PodSecurityChecks) *RuleResponse {
	r.podSecurityChecks = &checks
	return &r
}

func (r RuleResponse) WithPatchedTarget(patchedTarget *unstructured.Unstructured, gvr metav1.GroupVersionResource, subresource string) *RuleResponse {
	r.patchedTarget = patchedTarget
	r.patchedTargetParentResourceGVR = gvr
	r.patchedTargetSubresourceName = subresource
	return &r
}

func (r RuleResponse) WithGeneratedResources(resource []*unstructured.Unstructured) *RuleResponse {
	r.generatedResources = resource
	return &r
}

func (r RuleResponse) WithStats(stats ExecutionStats) RuleResponse {
	r.stats = stats
	return r
}

func (r RuleResponse) WithEmitWarning(emitWarning bool) *RuleResponse {
	r.emitWarning = emitWarning
	return &r
}

func (r RuleResponse) WithProperties(m map[string]string) *RuleResponse {
	r.properties = m
	return &r
}

func (r *RuleResponse) Stats() ExecutionStats {
	return r.stats
}

func (r *RuleResponse) Exceptions() []GenericException {
	return r.exceptions
}

func (r *RuleResponse) ValidatingAdmissionPolicyBinding() *admissionregistrationv1.ValidatingAdmissionPolicyBinding {
	return r.vapBinding
}

func (r *RuleResponse) MutatingAdmissionPolicyBinding() *admissionregistrationv1alpha1.MutatingAdmissionPolicyBinding {
	return r.mapBinding
}

func (r *RuleResponse) IsException() bool {
	return len(r.exceptions) > 0
}

func (r *RuleResponse) PodSecurityChecks() *PodSecurityChecks {
	return r.podSecurityChecks
}

func (r *RuleResponse) PatchedTarget() (*unstructured.Unstructured, metav1.GroupVersionResource, string) {
	return r.patchedTarget, r.patchedTargetParentResourceGVR, r.patchedTargetSubresourceName
}

func (r *RuleResponse) GeneratedResources() []*unstructured.Unstructured {
	return r.generatedResources
}

func (r *RuleResponse) Message() string {
	return r.message
}

func (r *RuleResponse) Name() string {
	return r.name
}

func (r *RuleResponse) RuleType() RuleType {
	return r.ruleType
}

func (r *RuleResponse) Status() RuleStatus {
	return r.status
}

func (r *RuleResponse) EmitWarning() bool {
	return r.emitWarning
}

func (r *RuleResponse) Properties() map[string]string {
	return r.properties
}

// HasStatus checks if rule status is in a given list
func (r *RuleResponse) HasStatus(status ...RuleStatus) bool {
	for _, s := range status {
		if r.status == s {
			return true
		}
	}
	return false
}

// String implements Stringer interface
func (r *RuleResponse) String() string {
	return fmt.Sprintf("rule %s (%s): %v", r.name, r.ruleType, r.message)
}
