package main

import (
	"context"
	"runtime"
	"time"

	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	stub "github.com/snarlysodboxer/ingress-controller-controller/pkg/stub"

	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel) // TODO make this configurable
	printVersion()

	sdk.ExposeMetricsPort()
	metrics, err := stub.RegisterOperatorMetrics()
	if err != nil {
		logrus.Errorf("failed to register operator specific metrics: %v", err)
	}
	handler := stub.NewHandler(metrics)

	resource := "v1"
	kind := "Service"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("Failed to get watch namespace: %v", err)
	}
	resyncPeriod := time.Duration(20) * time.Second // TODO make this configurable

	selector := "icc-operator=true" // TODO make this configurable
	watchOption := sdk.WithLabelSelector(selector)
	logrus.Infof("Watching %s, %s, %s, %d, %s", resource, kind, namespace, resyncPeriod, selector)
	sdk.Watch(resource, kind, namespace, resyncPeriod, watchOption)

	sdk.Handle(handler)
	sdk.Run(context.TODO())
}
