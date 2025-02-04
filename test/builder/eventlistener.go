package builder

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
)

// EventListenerOp is an operation which modifies the EventListener.
type EventListenerOp func(*v1alpha1.EventListener)

// EventListenerSpecOp is an operation which modifies the EventListenerSpec.
type EventListenerSpecOp func(*v1alpha1.EventListenerSpec)

// EventListenerStatusOp is an operation which modifies the EventListenerStatus.
type EventListenerStatusOp func(*v1alpha1.EventListenerStatus)

// EventListenerTriggerOp is an operation which modifies the Trigger.
type EventListenerTriggerOp func(*v1alpha1.EventListenerTrigger)

// EventListenerTriggerValidateOp is an operation which modifies the TriggerValidate.
type EventListenerTriggerValidateOp func(*v1alpha1.TriggerValidate)

// EventListener creates an EventListener with default values.
// Any number of EventListenerOp modifiers can be passed to transform it.
func EventListener(name, namespace string, ops ...EventListenerOp) *v1alpha1.EventListener {
	e := &v1alpha1.EventListener{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	for _, op := range ops {
		op(e)
	}

	return e
}

// EventListenerMeta sets the Meta structs of the EventListener.
// Any number of MetaOp modifiers can be passed.
func EventListenerMeta(ops ...MetaOp) EventListenerOp {
	return func(e *v1alpha1.EventListener) {
		for _, op := range ops {
			switch o := op.(type) {
			case ObjectMetaOp:
				o(&e.ObjectMeta)
			case TypeMetaOp:
				o(&e.TypeMeta)
			}
		}
	}
}

// EventListenerSpec sets the specified spec of the EventListener.
// Any number of EventListenerSpecOp modifiers can be passed to create/modify it.
func EventListenerSpec(ops ...EventListenerSpecOp) EventListenerOp {
	return func(e *v1alpha1.EventListener) {
		for _, op := range ops {
			op(&e.Spec)
		}
	}
}

// EventListenerServiceAccount sets the specified ServiceAccount of the EventListener.
func EventListenerServiceAccount(saName string) EventListenerSpecOp {
	return func(spec *v1alpha1.EventListenerSpec) {
		spec.ServiceAccountName = saName
	}
}

// EventListenerTrigger adds an EventListenerTrigger to the EventListenerSpec Triggers.
// Any number of EventListenerTriggerOp modifiers can be passed to create/modify it.
func EventListenerTrigger(tbName, ttName, apiVersion string, ops ...EventListenerTriggerOp) EventListenerSpecOp {
	return func(spec *v1alpha1.EventListenerSpec) {
		spec.Triggers = append(spec.Triggers, Trigger(tbName, ttName, apiVersion, ops...))
	}
}

// EventListenerTriggerParam adds a param to the EventListenerTrigger
func EventListenerTriggerParam(name, value string) EventListenerTriggerOp {
	return func(trigger *v1alpha1.EventListenerTrigger) {
		trigger.Params = append(trigger.Params,
			pipelinev1.Param{
				Name: name,
				Value: pipelinev1.ArrayOrString{
					StringVal: value,
					Type:      pipelinev1.ParamTypeString,
				},
			})
	}
}

// EventListenerStatus sets the specified status of the EventListener.
// Any number of EventListenerStatusOp modifiers can be passed to create/modify it.
func EventListenerStatus(ops ...EventListenerStatusOp) EventListenerOp {
	return func(e *v1alpha1.EventListener) {
		for _, op := range ops {
			op(&e.Status)
		}
	}
}

// EventListenerConditition sets the specified condition on the EventListenerStatus.
func EventListenerCondition(t apis.ConditionType, status corev1.ConditionStatus, message, reason string) EventListenerStatusOp {
	return func(e *v1alpha1.EventListenerStatus) {
		e.SetCondition(&apis.Condition{
			Type:    t,
			Status:  status,
			Message: message,
			Reason:  reason,
		})
	}
}

// EventListenerConfig sets the EventListenerConfiguration on the EventListenerStatus.
func EventListenerConfig(generatedResourceName string) EventListenerStatusOp {
	return func(e *v1alpha1.EventListenerStatus) {
		e.Configuration.GeneratedResourceName = generatedResourceName
	}
}

// EventListenerAddress sets the EventListenerAddress on the EventListenerStatus
func EventListenerAddress(hostname string) EventListenerStatusOp {
	return func(e *v1alpha1.EventListenerStatus) {
		e.Address = &duckv1alpha1.Addressable{
			Hostname: hostname,
		}
	}
}

// Trigger creates an EventListenerTrigger. Any number of EventListenerTriggerOp
// modifiers can be passed to create/modify it.
func Trigger(tbName, ttName, apiVersion string, ops ...EventListenerTriggerOp) v1alpha1.EventListenerTrigger {
	t := v1alpha1.EventListenerTrigger{
		Binding: v1alpha1.EventListenerBinding{
			Name:       tbName,
			APIVersion: apiVersion,
		},
		Template: v1alpha1.EventListenerTemplate{
			Name:       ttName,
			APIVersion: apiVersion,
		},
	}

	for _, op := range ops {
		op(&t)
	}

	return t
}

// EventListenerTriggerName adds a Name to the Trigger in EventListenerSpec Triggers.
func EventListenerTriggerName(name string) EventListenerTriggerOp {
	return func(trigger *v1alpha1.EventListenerTrigger) {
		trigger.Name = name
	}
}

// EventListenerTriggerValidate adds a TriggerValidate to the Trigger in EventListenerSpec Triggers.
func EventListenerTriggerValidate(ops ...EventListenerTriggerValidateOp) EventListenerTriggerOp {
	return func(trigger *v1alpha1.EventListenerTrigger) {
		validate := &v1alpha1.TriggerValidate{}
		for _, op := range ops {
			op(validate)
		}
		trigger.TriggerValidate = validate
	}
}

// EventListenerTriggerValidateTaskRef adds a TaskRef to the TriggerValidate.
func EventListenerTriggerValidateTaskRef(taskName, apiVersion string, kind pipelinev1.TaskKind) EventListenerTriggerValidateOp {
	return func(validate *v1alpha1.TriggerValidate) {
		validate.TaskRef = pipelinev1.TaskRef{
			Name:       taskName,
			Kind:       kind,
			APIVersion: apiVersion,
		}
	}
}

// EventListenerTriggerValidateServiceAccount adds a service account name to the TriggerValidate.
func EventListenerTriggerValidateServiceAccount(serviceAccount string) EventListenerTriggerValidateOp {
	return func(validate *v1alpha1.TriggerValidate) {
		validate.ServiceAccountName = serviceAccount
	}
}

// EventListenerTriggerValidateParam adds a param name to the TriggerValidate.
func EventListenerTriggerValidateParam(name, value string) EventListenerTriggerValidateOp {
	return func(validate *v1alpha1.TriggerValidate) {
		validate.Params = append(validate.Params,
			pipelinev1.Param{
				Name: name,
				Value: pipelinev1.ArrayOrString{
					StringVal: value,
					Type:      pipelinev1.ParamTypeString,
				},
			},
		)
	}
}
