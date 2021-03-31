package testcases

import (
	"context"
	"errors"
	"fmt"

	"github.com/thatchd/thatchd-sample/testsuite"
	"github.com/thatchd/thatchd/pkg/thatchd/testcase"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodAnnotationTestCase asserts that a pod has an annotation with an expected
// value
type PodAnnotationTestCase struct {
	PodName            string
	ExpectedAnnotation string
	ExpectedValue      string
}

var _ testcase.Interface = &PodAnnotationTestCase{}

func (tc *PodAnnotationTestCase) ShouldRun(s interface{}) bool {
	state := s.(testsuite.PodSuiteState)
	podState, ok := state[tc.PodName]
	return ok && podState == testsuite.PodAnnotated
}

func (tc *PodAnnotationTestCase) Run(c client.Client, namespace string) error {
	pod := &v1.Pod{}
	if err := c.Get(context.TODO(), client.ObjectKey{
		Name:      tc.PodName,
		Namespace: namespace,
	}, pod); err != nil {
		return fmt.Errorf("failed to obtain pod: %v", err)
	}

	annotations := pod.Annotations
	if annotations == nil {
		return errors.New("pod has no annotations")
	}

	if value, ok := annotations[tc.ExpectedAnnotation]; !ok || value != tc.ExpectedValue {
		return fmt.Errorf("annotation %s: %s not found in Pod", tc.ExpectedAnnotation, tc.ExpectedValue)
	}

	return nil
}

func NewTestCase(config map[string]string) interface{} {
	return &PodAnnotationTestCase{
		PodName:            config["podName"],
		ExpectedAnnotation: config["expectedAnnotation"],
		ExpectedValue:      config["expectedValue"],
	}
}
