# Thatchd Sample Project

This is a simple [Thatchd](https://github.com/thatchd/thatchd)-powered testing operator.
It reconciles the components to test that a Pod has an expected annotation.

## Components

| Resource | Component | Logic |
| -------- | --------- | ----- |
| `TestSuite` | `PodsSuite` | Reconciles a state that represents Pod readiness: a JSON object where key = Pod name, value = Pod readiness |
| `TestCase` | `PodAnnotation` | Asserts that a given Pod contains an expected annotation. Dispatched when the Pod state is `PodAnnotated` |
| `TestWorker` | `PodAnnotationWorker` | Sets an annotation on a Pod. Dispatched when the Pod state is `PodReady` |

## Try it

The sample project can run locally or deployed on a cluster. This section describes
how to try it locally.

### Pre-requisites

* Admin access to a Kubernetes cluster

### Set up

Clone the repo and install resources in the cluster
```sh
git clone https://github.com/thatchd/thatchd-sample.git
cd thatchd-sample
make install
```

Start running the operator
```sh
make run
```

### Create CRs

#### TestSuite

> See the source code of the example TestSuite reconciler:
>
> 👓 [testsuite/podsuite.go](testsuite/podsuite.go)

Create the TestSuite CR with the `PodsSuite` provider

```yaml
apiVersion: testing.thatchd.io/v1alpha1
kind: TestSuite
metadata:
  name: test-pods
spec:
  initialState: '{}'
  stateStrategy:
    provider: PodsSuite
```

Once created, Thatchd will reconcile the status with a list of Pods in the namespace.
Go ahead and create a simple Pod. Thatchd will reconcile the `status` field accordingly

```yaml
status:
  currentState: |-
    {
      "my-pod": true
    }
```

> ℹ You can use any Go type as test state, leveraging the language type information

#### TestCase

> See the source code of the example TestCase implementation:
>
> 👓 [testcases/podannotation.go](testcases/podannotation.go)

The example test case will be dispatched when a specific pod is annotated according
to the TestSuite state. Create a TestCase CR to verify that the `foo: bar` annotation
is set on the `test-success` Pod

```yaml
apiVersion: testing.thatchd.io/v1alpha1
kind: TestCase
metadata:
  name: testcase-success
spec:
  strategy:
    configuration:
      expectedAnnotation: foo
      expectedValue: bar
      podName: test-success
    provider: PodAnnotation
```

> ℹ️ The `configuration` field in the CR allows to reuse logic in multiple test cases

The test case won't be dispatched yet as the Pod hasn't been created

#### TestWorker

> See the source code of the example TestWorker implementation:
>
> 👓 [testworkers/podannotation.go](testworkers/podannotation.go)

The example test worker will be dispatched when a specific pod is ready, and will
annotate the pod with the configured annotation. Create a TestWorker CR to annotate
the `test-success` Pod with `foo: bar`

```yaml
apiVersion: testing.thatchd.io/v1alpha1
kind: TestWorker
metadata:
  name: testworker-success
spec:
  strategy:
    configuration:
      annotation: foo
      value: bar
      podName: test-success
    provider: PodAnnotationWorker
```

#### Test subject: `test-success` Pod

Create the Pod called `test-success`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-success
  labels:
    app: hello-openshift
spec:
  containers:
    - name: hello-openshift
      image: openshift/hello-openshift
      ports:
        - containerPort: 8080
```

Once the Pod is ready, the TestWoker will be dispatched, and quickly executed,
annotating the Pod and setting the suite status. When the Pod status is set
to annotated, the TestCase will be dispatched and executed, verifying the
annotation and setting the status to `Finished`