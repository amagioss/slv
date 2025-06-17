/*
Copyright 2024.

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

package operator

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/open-policy-agent/cert-controller/pkg/rotator"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/session"
	slvv1 "slv.sh/slv/internal/k8s/api/v1"
	"slv.sh/slv/internal/k8s/internal/controller"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var (
	// Webhook settings
	disableWebhook      = getEnvBool("SLV_DISABLE_WEBHOOKS", false)
	disableCertRotation = getEnvBool("SLV_DISABLE_CERT_ROTATION", false)
	certServiceName     = getEnvOrDefault("SLV_WEBHOOK_SERVICE_NAME", "slv-webhook-service")
	VwhName             = getEnvOrDefault("SLV_WEBHOOK_VWH_NAME", "slv-operator-validating-webhook")

	// Certificate and secret configuration
	secretName     = getEnvOrDefault("SLV_WEBHOOK_SECRET_NAME", "slv-webhook-server-cert")
	caName         = getEnvOrDefault("SLV_WEBHOOK_CA_NAME", "slv-webhook-ca")
	caOrganization = getEnvOrDefault("SLV_WEBHOOK_CA_ORG", "slv")
	certDir        = getEnvOrDefault("SLV_WEBHOOK_CERT_DIR", "/tmp/k8s-webhook-server/serving-certs")

	// Certificate rotation durations
	caCertDuration         = getEnvDuration("SLV_CA_CERT_DURATION", 3*365*24*time.Hour)   // 3 years
	serverCertDuration     = getEnvDuration("SLV_SERVER_CERT_DURATION", 365*24*time.Hour) // 1 year
	rotationCheckFrequency = getEnvDuration("SLV_ROTATION_CHECK_FREQUENCY", 24*time.Hour) // every 24 hours
	lookaheadInterval      = getEnvDuration("SLV_LOOKAHEAD_INTERVAL", 7*24*time.Hour)     // one week before expiration
)

// TODO: Move to utils
func getEnvOrDefault(envKey, defaultValue string) string {
	if value, exists := os.LookupEnv(envKey); exists && value != "" {
		return strings.ToLower(value)
	}
	return defaultValue
}

// TODO: Move to utils
func getEnvBool(envKey string, defaultVal bool) bool {
	val := getEnvOrDefault(envKey, "")
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// TODO: Move to utils
func getEnvDuration(envKey string, defaultVal time.Duration) time.Duration {
	val := getEnvOrDefault(envKey, "")
	if val == "" {
		return defaultVal
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(slvv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func Run() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool

	// Define webhook
	var webhooks []rotator.WebhookInfo
	webhooks = append(webhooks, rotator.WebhookInfo{
		Name: VwhName,
		Type: rotator.Validating,
	})

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set the metrics endpoint is served securely")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	setupLog.Info("initializing SLV operator...")
	setupLog.Info(config.VersionInfo())

	if _, err := session.GetSecretKey(); err != nil {
		setupLog.Error(err, "unable to initialize slv environment")
		os.Exit(1)
	}

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancelation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	tlsOpts := []func(*tls.Config){}
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "736c7673.slv.sh",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	namespace := session.GetK8sNamespace()

	setupFinished := make(chan struct{})
	if !disableCertRotation {
		setupLog.Info("setting up cert rotation")

		if err := rotator.AddRotator(mgr, &rotator.CertRotator{
			SecretKey: types.NamespacedName{
				Namespace: namespace,
				Name:      secretName,
			},
			CertDir:                certDir,
			CAName:                 caName,
			CAOrganization:         caOrganization,
			DNSName:                fmt.Sprintf("%s.%s.svc", certServiceName, namespace),
			IsReady:                setupFinished,
			Webhooks:               webhooks,
			RequireLeaderElection:  enableLeaderElection,
			CaCertDuration:         caCertDuration,
			ServerCertDuration:     serverCertDuration,
			RotationCheckFrequency: rotationCheckFrequency,
			LookaheadInterval:      lookaheadInterval,
			EnableReadinessCheck:   true,
		}); err != nil {
			setupLog.Error(err, "unable to set up cert rotation")
			os.Exit(1)
		}
	} else {
		close(setupFinished)
	}

	if err = (&controller.SLVReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", slvv1.Kind)
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	if !disableWebhook {
		if err := mgr.Add(manager.RunnableFunc((func(ctx context.Context) error {
			if !disableCertRotation {
				setupLog.Info("waiting for certs to be ready before registering webhook")
				<-setupFinished
				setupLog.Info("certs ready, setting up webhook")
			} else {
				setupLog.Info("skipping cert rotation, setting up webhook")
			}
			if err := (&slvv1.SLV{}).SetupWebhookWithManager(mgr); err != nil {
				return fmt.Errorf("failed to set up SLV webhook: %w", err)
			}
			return nil
		}))); err != nil {
			setupLog.Error(err, "unable to register webhook setup hook")
			os.Exit(1)
		}
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

}
