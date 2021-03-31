/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/thatchd/thatchd-sample/testcases"
	"github.com/thatchd/thatchd-sample/testsuite"
	"github.com/thatchd/thatchd-sample/testworkers"
	"github.com/thatchd/thatchd/pkg/thatchd/manager"
	"github.com/thatchd/thatchd/pkg/thatchd/strategy"
)

// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testcases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testcases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testsuites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testsuites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testworkers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testing.thatchd.io,resources=testworkers/status,verbs=get;update;patch

func addToScheme(scheme *runtime.Scheme) error {
	return nil
}

func main() {
	strategyProviders := map[string]strategy.StrategyProvider{
		"PodsSuite":           testsuite.NewPodsSuiteProvider(),
		"PodAnnotation":       strategy.NewProviderFunction(testcases.NewTestCase),
		"PodAnnotationWorker": strategy.NewProviderFunction(testworkers.NewTestWorker),
	}

	manager.Run(addToScheme, strategyProviders)
}
