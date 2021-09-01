---
title: "Manage cluster with GitOps"
linkTitle: "Manage cluster with GitOps"
weight: 30
date: 2017-01-05
description: >
  Use Flux to manage clusters with GitOps
---

## GitOps Support (optional)

EKS-A supports a [GitOps](https://www.weave.works/technologies/gitops/) workflow for the management of your cluster.

When you create a cluster with GitOps enabled, EKS-A will automatically commit your cluster configuration to the provided GitHub repository and install a GitOps toolkit on your cluster which watches that committed configuration file.
You can then manage the scale of the cluster by making changes to the version controlled cluster configuration file and committing the changes.
Once a change is detected by the GitOps controller running in your cluster, the scale of the cluster will be adjusted to match the committed configuration file.

If you'd like to learn more about GitOps and the associated best practices, [check out this introduction from Weaveworks](https://www.weave.works/technologies/gitops/).

>**_NOTE:_** Installing a GitOps controller needs to be done during cluster creation.
In the event that GitOps installation fails, EKS-A cluster creation will continue.

### Supported Cluster Properties

Currently, you can manage a subset of cluster properties with GitOps:

- `Cluster.workerNodeGroupConfigurations[0].count`
- `Cluster.workerNodeGroupConfigurations[0].machineGroupRef.name`

For a VsphereMachineConfig associated with worker nodes via the workerNodeGroups.machineGroupRef.name, you may update the following fields with GitOps:

- `VsphereMachineConfig.diskGiB`
- `VsphereMachineConfig.numCPUs`
- `VsphereMachineConfig.memoryMiB`
- `VsphereMachineConfig.template`

Any other changes to the cluster configuration in the git repository will be ignored.

## Getting Started with EKS-A GitOps

In order to use GitOps to manage cluster scaling, you need a couple of things:

- A GitHub account
- A cluster configuration file with a `GitOpsConfig`, referenced with a `gitOpsRef` in your Cluster spec
- A [Personal Access Token (PAT) for the GitHub account](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token), with permissions to create, clone, and push to a repo

### Create a GitHub Personal Access Token

[Create a PAT](https://github.com/settings/tokens/new) to access your provided github repository.
It must be scoped for all `repo` permissions.

>**_NOTE:_** GitOps configuration only works with hosted github.com and will not work on self hosted GitHub Enterprise instances.

This PAT should have at least the following permissions:

![Github PAT permissions](/images/ss5.png)

>**_NOTE:_** The PAT must belong to the `owner` of the `repository` or, if using an organization as the `owner`, the creator of the `PAT` must have repo permission in that organization.

You can provide the PAT to EKS-A in one of two ways:

* (default) Set the `$GITHUB_TOKEN` environment variable to the value of your PAT

   ```
   export GITHUB_TOKEN=ghp_MyValidPersonalAccessTokenWithRepoPermissions
   ```

* A Github Token file with read permissions.
  The file should be plain text with only the token text inside.
  **Do not commit and push this file to a public GitHub repo.**

  Provide the full path to the file in your GitOps configuration in the field `authTokenPath` of the `flux` `gitops` spec.
  ```
  apiVersion: anywhere.eks.amazonaws.com/v1alpha1
  kind: GitOpsConfig
  metadata:
    name: my-gitops-config
  spec:
    flux:
      github:
        personal: true
        repository: mygithubrepository
        owner: mygithubusername
        authTokenPath: ./GH_TOKEN.txt
  ```

### Create GitOps configuration repo

If you have an existing repo you can set that as your repository name in the configuration.
If you specify a repo in your `GitOpsConfig` which does not exist EKS-A will create it for you.
If you would like to create a new repo you can [click here](https://github.new) to create a new repo.

Your repo will contain your cluster specification file(s).
If you have multiple files you should store them in subfolders and specify the [path in your configuration]({{< relref "#__path__-optional">}}).

Example repo structure:
```
clusters
├── dev
│   └── eksa-system
│       └── eksa-cluster.yaml
├── prod
│   └── eksa-system
│       └── eksa-cluster.yaml
└── stage
    └── eksa-system
        └── eksa-cluster.yaml
```

### Example GitOps cluster configuration

```yaml
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: Cluster
metadata:
  name: mynewgitopscluster
spec:
... # collapsed cluster spec fields
# Below added for gitops support
  gitOpsRef:
    kind: GitOpsConfig
    name: my-cluster-name
---
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: GitOpsConfig
metadata:
  name: my-cluster-name
spec:
  flux:
    github:
      personal: true
      repository: mygithubrepository
      owner: mygithubusername
```

### Create a GitOps enabled cluster

Generate your cluster configuration and add the GitOps configuration.
For a full spec reference see the [Cluster Spec reference]({{< relref "../../reference/clusterspec/gitops" >}}).

>**_NOTE:_** After your cluster is created the cluster configuration will automatically be commited to your git repo.

1. Create an EKS-A cluster with GitOps enabled.

    ```bash
    CLUSTER_NAME=gitops
    eksctl anywhere create cluster -f ${CLUSTER_NAME}.yaml
    ```

### Test GitOps controller

After your cluster is created you can test the GitOps contoller by modifying the cluster specification.

1. Clone your git repo and modify the cluster specification.
   The default path for the cluster file is:

    ```
    clusters/$CLUSTER_NAME/eksa-system/eksa-cluster.yaml
    ```

1. Modify the `workerNodeGroupsConfigurations[0].count` field with your desired changes.

1. Commit the file to your git repository

    ```bash
    git add eksa-cluster.yaml
    git commit -m 'Scaling nodes for test'
    git push origin main
    ```

1. The flux controller will automatically make the required changes.

   If you updated your node count you can use this command to see the current node state.
    ```bash
    kubectl get nodes 
    ```
   