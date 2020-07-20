# Image Registry Operator

[![GoDoc](https://godoc.org/github.com/openshift/cluster-image-registry-operator?status.png)](https://godoc.org/github.com/openshift/cluster-image-registry-operator)
[![Licensed under Apache License version 2.0](https://img.shields.io/github/license/openshift/cluster-image-registry-operator.svg?maxAge=2592000)](https://www.apache.org/licenses/LICENSE-2.0)

# Overview

The registry operator manages a singleton instance of the openshift registry.  It manages all configuration of the registry including creating storage.

On initial startup the operator will create a default image-registry resource instance based on configuration detected in the cluster (e.g. what cloud storage type to use based on the cloud provider).

If insufficient information is available to define a complete image-registry resource, the incomplete resource will be defined and the operator will update the resource status with information about what is missing.

The registry operator runs in the openshift-image-registry namespace, and manages the registry instance in that location as well.  All configuration+workload resources for the registry reside in that namespace.

# Configuration

## Registry resource

The image-registry resource offers the following configuration fields:

* ManagementState
  * Managed - the operator will update the registry as configuration resources are updated
  * Unmanaged - the operator will ignore changes to the configuration resources
  * Removed - the operator will remove the registry instance and tear down any storage that the operator provisioned.
* Logging
  * Sets loglevel of the registry instance
* HTTPSecret
  * Value needed by the registry to secure uploads, generated by default.
* Proxy
  * Not currently implemented
  * Defines the Proxy to be used when calling master api, upstream registries, etc
* Storage
  * Storagetype details for configuring registry storage, e.g. S3 bucket coordinates.
  * Normally configured by default
* Requests
  * API Request Limit details
  * Controls how many parallel requests a given registry instance will handle before queuing additional requests
* DefaultRoute
  * Determines whether or not an external route is defined using the default hostname.  If enabled, the route uses re-encrypt encryption.
  * Defaults to false today
* Routes
  * Array of additional routes to create
  * User provides hostname, certificate for the route
* Replicas
  * Replica count for the registry


## Additional config resources

In addition to the image-registry resource, additional config is provided to the operator via separate configmap + secret resources located within the openshift-image-registry namespace:

### image-registry-certificates (configmap)

Provides additional CAs for contacting upstream registries.  Mounted to `/etc/pki/ca-trust/source/anchors` in the registry pod.

### image-registry-private-configuration-user (secret)

Provides credentials needed for storage management/access, overrides the default
credentials used by the operator, if default credentials were found.

For S3 storage it is expected to contain two keys whose values are the AWS access key and secret key that you want to use:
* REGISTRY_STORAGE_S3_ACCESSKEY
* REGISTRY_STORAGE_S3_SECRETKEY

For GCS storage it is expected to contain one key whose value is the contents of a credentials file provided by GCP:
* REGISTRY_STORAGE_GCS_KEYFILE

For Azure storage it is expected to contain one key whose value is an account key:
* REGISTRY_STORAGE_AZURE_ACCOUNTKEY

# Troubleshooting

The registry operator reports status in two places:

A ClusterOperator resource is defined in the cluster scope which reflects the state of the registry operator at a high level.  Retrievable via:

    oc get clusteroperators.config.openshift.io/image-registry -o yaml

The image-registry resource itself has a status section with detailed conditions indicating the state of the managed registry, you can view this via:

    oc get configs.imageregistry.operator.openshift.io/cluster -o yaml


**If you cannot access your registry, check the following:**

Is the registry deployed?  Check for a registry deployment + corresponding pod in the openshift-image-registry namespace:

    oc get deployment image-registry -n openshift-image-registry
    oc get pods -n openshift-image-registry | grep image-registry | grep -v operator

**If there is no registry pod, check the deployment for any error conditions:**

    oc get deployment image-registry -o yaml -n openshift-image-registry

**If there is no registry deployment, check the image-registry resource instance for any error conditions:**

    oc get configs.imageregistry.operator.openshift.io/cluster -o yaml -n openshift-image-registry

**If there is no image-registry resource at all, check if the image-registry operator deployment exists:**

    oc get deployment/cluster-image-registry-operator -n openshift-image-registry

**If the operator deployment exists, check for the corresponding pod and, if it exists, check its logs:**

    oc get pods  -n openshift-image-registry | grep cluster-image-registry-operator
    oc logs cluster-image-registry-operator-5c8bcf89bb-4nr8p -n openshift-image-registry

**If the operator pod does not exist, inspect the deployment to determine why the operator pod was not created:**

    oc get deployment cluster-image-registry-operator -o yaml -n openshift-image-registry

**If the deployment does not exist:**

Something went wrong at the installer/CVO level that it did not deploy the image-registry operator.
