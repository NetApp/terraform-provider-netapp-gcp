# CVS/GCP Terraform provider examples

This folder contains examples for the resources (ActiveDirectory, Volumes, Snapshots) and
the data sources (Volumes, ActiveDirectory(future)) the provider supports.

Each example is self-contained and doen't require files from other folders (provider.tf
might be an exception). The tf-file named after the resource of interest in the example
contains a short example description a comments. Intended for copy/paste for your resources.

All examples require Terraform 0.13 or higher.

## Overview of examples:

```bash
.
├── active_directory
│   ├── datasource - Queries a region for existing Active Directory connection
│   └── minimal - Creates an Active Directory connection
├── snapshot - Creates a snapshot for an volume
├── volume_replication - Creates two volumes, one as secondary for volume replication relationship 
└── volume - contains examples for volume manipulation
    ├── advanced - More complex example using variables and outputs
    ├── datasource - Shows how to query an existing volume for later TF use
    ├── minimal - Most simple volume example. Can be considered as "Hello World"
    ├── volume-batch - Creates multiple volumes out of a CSV file
    ├── volume-nfs - Creates an NFS volume
    └── volume-smb - Creates an SMB volume
```

## Usage
Change into respective folder and change provider.tf to suite your GCP project. Next look for the "locals" block in the *.tf files and customize it to for your environment. Then run "terraform init" followed by "terraform apply".
