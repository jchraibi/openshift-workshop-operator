package workshop

import (
	"context"

	securityv1 "github.com/openshift/api/security/v1"
	openshiftv1alpha1 "github.com/redhat/openshift-workshop-operator/pkg/apis/openshift/v1alpha1"
	deployment "github.com/redhat/openshift-workshop-operator/pkg/deployment"
	"github.com/redhat/openshift-workshop-operator/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciling Vault
func (r *ReconcileWorkshop) reconcileVault(instance *openshiftv1alpha1.Workshop, users int) (reconcile.Result, error) {
	enabled := instance.Spec.Infrastructure.Vault.Enabled

	if enabled {
		if result, err := r.addVaultServer(instance, users); err != nil {
			return result, err
		}

		if result, err := r.addVaultAgentInjector(instance, users); err != nil {
			return result, err
		}

		// Installed
		if instance.Status.Vault != util.OperatorStatus.Installed {
			instance.Status.Vault = util.OperatorStatus.Installed
			if err := r.client.Status().Update(context.TODO(), instance); err != nil {
				logrus.Errorf("Failed to update Workshop status: %s", err)
				return reconcile.Result{}, err
			}
		}
	}

	//Success
	return reconcile.Result{}, nil
}

func (r *ReconcileWorkshop) addVaultServer(instance *openshiftv1alpha1.Workshop, users int) (reconcile.Result, error) {
	labels := map[string]string{
		"app":                       "vault",
		"app.kubernetes.io/name":    "vault",
		"app.kubernetes.io/part-of": "vault",
		"component":                 "server",
	}

	namespace := deployment.NewNamespace(instance, "vault")
	if err := r.client.Create(context.TODO(), namespace); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Project", namespace.Name)
	}

	extraconfigFromValues := map[string]string{
		"extraconfig-from-values.hcl": `disable_mlock = true
ui = true

listener "tcp" {
	tls_disable = 1
	address = "[::]:8200"
	cluster_address = "[::]:8201"
}
storage "file" {
	path = "/vault/data"
}
`,
	}

	configMap := deployment.NewConfigMap(instance, "vault-config", namespace.Name, labels, extraconfigFromValues)
	if err := r.client.Create(context.TODO(), configMap); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s ConfigMap", configMap.Name)
	}

	// Create Service Account
	serviceAccount := deployment.NewServiceAccount(instance, "vault", namespace.Name)
	if err := r.client.Create(context.TODO(), serviceAccount); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Service Account", serviceAccount.Name)
	}

	serviceAccountUser := "system:serviceaccount:" + namespace.Name + ":" + serviceAccount.Name

	privilegedSCCFound := &securityv1.SecurityContextConstraints{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "privileged"}, privilegedSCCFound); err != nil {
		return reconcile.Result{}, err
	}
	if !util.StringInSlice(serviceAccountUser, privilegedSCCFound.Users) {
		privilegedSCCFound.Users = append(privilegedSCCFound.Users, serviceAccountUser)
		if err := r.client.Update(context.TODO(), privilegedSCCFound); err != nil {
			return reconcile.Result{}, err
		} else if err == nil {
			logrus.Infof("Updated %s SCC", privilegedSCCFound.Name)
		}
	}

	// Create ClusterRole Binding
	clusterRoleBinding := deployment.NewClusterRoleBindingForServiceAccount(instance, "vault-server-binding", namespace.Name,
		serviceAccount.Name, "system:auth-delegator", "ClusterRole")
	if err := r.client.Create(context.TODO(), clusterRoleBinding); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Cluster Role Binding", clusterRoleBinding.Name)
	}

	// Create Service
	internalService := deployment.NewService(instance, "vault-internal", namespace.Name, labels, []string{"http", "internal"}, []int32{8200, 8201})
	if err := r.client.Create(context.TODO(), internalService); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Service", internalService.Name)
	}

	service := deployment.NewService(instance, "vault", namespace.Name, labels, []string{"http", "internal"}, []int32{8200, 8201})
	if err := r.client.Create(context.TODO(), service); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Service", service.Name)
	}

	// Create Stateful
	stateful := deployment.NewVaultStateful(instance, "vault", namespace.Name, labels)
	if err := r.client.Create(context.TODO(), stateful); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Stateful", stateful.Name)
	}

	//Success
	return reconcile.Result{}, nil
}

func (r *ReconcileWorkshop) addVaultAgentInjector(instance *openshiftv1alpha1.Workshop, users int) (reconcile.Result, error) {
	labels := map[string]string{
		"app":                       "vault",
		"app.kubernetes.io/name":    "vault-agent-injector",
		"app.kubernetes.io/part-of": "vault",
		"component":                 "webhook",
	}

	namespace := deployment.NewNamespace(instance, "vault")
	if err := r.client.Create(context.TODO(), namespace); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Project", namespace.Name)
	}

	// Create Service Account
	serviceAccount := deployment.NewServiceAccount(instance, "vault-agent-injector", namespace.Name)
	if err := r.client.Create(context.TODO(), serviceAccount); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Service Account", serviceAccount.Name)
	}

	serviceAccountUser := "system:serviceaccount:" + namespace.Name + ":" + serviceAccount.Name

	privilegedSCCFound := &securityv1.SecurityContextConstraints{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "privileged"}, privilegedSCCFound); err != nil {
		return reconcile.Result{}, err
	}

	if !util.StringInSlice(serviceAccountUser, privilegedSCCFound.Users) {
		privilegedSCCFound.Users = append(privilegedSCCFound.Users, serviceAccountUser)
		if err := r.client.Update(context.TODO(), privilegedSCCFound); err != nil {
			return reconcile.Result{}, err
		} else if err == nil {
			logrus.Infof("Updated %s SCC", privilegedSCCFound.Name)
		}
	}
	// Create Cluster Role
	clusterRole := deployment.NewClusterRole(instance, "vault-agent-injector", namespace.Name, deployment.VaultAgentInjectorRules())
	if err := r.client.Create(context.TODO(), clusterRole); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Cluster Role", clusterRole.Name)
	}

	clusterRoleBinding := deployment.NewClusterRoleBindingForServiceAccount(instance, "vault-agent-injector", namespace.Name,
		"vault-agent-injector", clusterRole.Name, "ClusterRole")
	if err := r.client.Create(context.TODO(), clusterRoleBinding); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Cluster Role Binding", clusterRoleBinding.Name)
	}

	// Create Service
	service := deployment.NewServiceWithTarget(instance, "vault-agent-injector", namespace.Name, labels,
		[]string{"http"}, []int32{443}, []int32{8080})
	if err := r.client.Create(context.TODO(), service); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Service", service.Name)
	}

	// Create Deployment
	ocpDeployment := deployment.NewVaultAgentInjectorDeployment(instance, "vault-agent-injector", namespace.Name, labels)
	if err := r.client.Create(context.TODO(), ocpDeployment); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Deployment", ocpDeployment.Name)
	}

	// Create
	webhooks := deployment.VaultAgentInjectorWebHook(namespace.Name)
	mutatingWebhookConfiguration := deployment.NewMutatingWebhookConfiguration(instance, "vault-agent-injector-cfg", labels, webhooks)
	if err := r.client.Create(context.TODO(), mutatingWebhookConfiguration); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		logrus.Infof("Created %s Mutating Webhook Configuration", mutatingWebhookConfiguration.Name)
	}

	//Success
	return reconcile.Result{}, nil
}
