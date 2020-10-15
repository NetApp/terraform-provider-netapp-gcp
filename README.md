# Terraform Provider for NetApp Cloud Volumes Service for Google Cloud

This is the repository for the Terraform Provider for NetApp Cloud Volumes Service (CVS) for Google Cloud.  The Provider can be used
with Terraform to work with Cloud Volumes Service for Google Cloud resources.

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

The provider plugin was developed by NetApp.

# Naming Conventions

The APIs for NetApp Cloud Volumes Service for Google Cloud do not require resource names to be unique.  They are considered
as 'labels' and resources are uniquely identified by 'ids'.  However these ids are not
user friendly, and as they are generated on the fly, they make it difficult to track
resources and automate.

This provider assumes that resource names are unique, and enforces it within its scope.
This is not an issue if everything is managed through Terraform, but could raise
conflicts if the rule is not respected outside of Terraform.

For the snapshot resource, you can also use the volume creation token (with or without
the volume name) to ensure uniqueness.

# Using the Provider

The current version of this provider requires Terraform 0.12 or higher to
run.   Though we recommend 0.13 or better.

For version 0.12, you will need to build the provider before being able to use it
(see [the section below](#building-the-provider)).

For version 0.13 and above, you can either build your provider, or use the verified provider in
Terraform registery. If you don't need to customize our provider, using 0.13 is highly recommended.
Please see Terraform website for more information:
https://registry.terraform.io/providers/NetApp/netapp-gcp/latest

If you are using 0.12 and would like to upgrade to 0.13, please see Terraform website for details:
https://www.terraform.io/upgrade-guides/0-13.html

Note that you need to run `terraform init` to fetch the provider before
deploying.

## Provider Documentation

The documentation is available at: https://registry.terraform.io/providers/NetApp/netapp-gcp/latest/docs

The provider is also documented [here][tf-netapp-gcp-docs].

Check the provider documentation for details on
entering your connection information and how to get started with writing
configuration for NetApp GCP resources.

[tf-netapp-gcp-docs]: website/docs/index.html.markdown

### Controlling the provider version

Note that you can also control the provider version. This is controlled by a
`required_provider` block in your Terraform configuration.

The syntax is as follows:

```hcl
terraform {
  required_providers {
    netapp-gcp = {
      source = "NetApp/netapp-gcp"
      version = "20.10.0"
    }
  }
}
```

[Read more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/provider-requirements.html#requiring-providers

# Building The Provider

## Prerequisites

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.11+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

The following go packages are required to build the provider:
```
	github.com/fatih/structs v1.1.0
	github.com/hashicorp/terraform v0.12.28
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
```

Check go.mod for the latest list.

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/terraform-provider-netapp-gcp`:

```sh
mkdir -p $GOPATH
cd $GOPATH
git clone https://github.com/NetApp/terraform-provider-netapp-gcp.git
```

## Running the Build

Note: check go.mod and make necessary updates to use TF 0.13 for instance.

After the clone operation is complete, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/terraform-provider-netapp-gcp
make build
```
Note: go install will move the binary to $GOPATH/bin

## Installing the Local Plugin

With Terraform 0.13 or newer, see the [sanity check](#sanity-check) section under **Walkthrough example**.

With earlier versions of Terraform, after
the build is complete, copy the `terraform-provider-netapp-gcp` binary into
the same path as your `terraform` binary, and run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/netapp/terraform-provider-netapp-gcp/issues
[gh-prs]: https://github.com/netapp/terraform-provider-netapp-gcp/pulls

See [Building the Provider](#building-the-provider) for details on building the provider.

# Testing the Provider

**NOTE:** Testing the provider for NetApp Cloud Volumes Service for Google Cloud is currently a complex operation as it
requires having a NetApp CVS subscription in GCP to test against.
You can then use a .json file to expose your credentials.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the
[`gcp/`](netapp_gcp/) directory for more details. The next section also
describes how you can manage a configuration file of the test environment
variables.

### Using the `.tf-netapp-gcp-devrc.mk` file

The [`tf-netapp-gcp-devrc.mk.example`](tf-netapp-gcp-devrc.mk.example) file contains
an up-to-date list of environment variables required to run the acceptance
tests. Copy this to `$HOME/.tf-netapp-gcp-devrc.mk` and change the permissions to
something more secure (ie: `chmod 600 $HOME/.tf-netapp-gcp-devrc.mk`), and
configure the variables accordingly.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccNetAppGCPVolume"
```

This following example would run all of the acceptance tests matching
`TestAccNetAppGCPSwVolume`. Change this for the specific tests you want to
run.


# Walkthrough example

### Installing go and terraform

```
bash
mkdir tf_na_gcp_cvs
cd tf_na_gcp_cvs

# if you want a private installation, use
export GO_INSTALL_DIR=`pwd`/go_install
mkdir $GO_INSTALL_DIR
# otherwise, go recommends to use
export GO_INSTALL_DIR=/usr/local


curl -O https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
tar -C $GO_INSTALL_DIR -xvf go1.14.1.linux-amd64.tar.gz

export PATH=$PATH:$GO_INSTALL_DIR/go/bin
```

#### Installing Terraform 0.12 
```
curl -O https://releases.hashicorp.com/terraform/0.12.24/terraform_0.12.24_linux_amd64.zip
unzip terraform_0.12.24_linux_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin
```

Note: for mac, subtitute linux_amd64 with darwin_amd64.

#### Installing Terraform 0.13 
```
curl -O https://releases.hashicorp.com/terraform/0.13.4/terraform_0.13.4_linux_amd64.zip
unzip terraform_0.13.4_linux_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin
```

### Installing dependencies

We're using go.mod to manage dependencies, so there is not much to do.
```
# make sure git is installed
which git

export GOPATH=`pwd`
```

### Cloning the NetApp provider repository and building the provider


```
git clone https://github.com/NetApp/terraform-provider-netapp-gcp.git
cd terraform-provider-netapp-gcp
make build
mv $GOPATH/bin/terraform-provider-netapp-gcp $GO_INSTALL_DIR/go/bin
```

The build step will install the provider in the $GOPATH/bin directory.

### Sanity check

#### Sanity check with TF 0.12

```
mv $GOPATH/bin/terraform-provider-netapp-gcp $GO_INSTALL_DIR/go/bin
cd examples/gcp/local_012
terraform init
```

Should do nothing but indicate that `Terraform has been successfully initialized!`

#### Sanity check with TF 0.13

0.13 is using registry by default, so it is a bit more involved to use a local build.

##### Local installation - linux

```
mkdir -p /tmp/terraform/netapp.com/netapp/netapp-gcp/20.10.0/linux_amd64
cp $GOPATH/bin/terraform-provider-netapp-gcp /tmp/terraform/netapp.com/netapp/netapp-gcp/20.10.0/linux_amd64
```

##### Local installation - mac

```
mkdir -p /tmp/terraform/netapp.com/netapp/netapp-gcp/20.10.0/darwin_amd64
cp $GOPATH/bin/terraform-provider-netapp-gcp /tmp/terraform/netapp.com/netapp/netapp-gcp/20.10.0/darwin_amd64
```

##### Check the provider can be loaded
```
cd examples/gcp/local
export TF_CLI_CONFIG_FILE=`pwd`/terraform.rc
terraform init
```

Should do nothing but indicate that `Terraform has been successfully initialized!`

#### Debugging

```
export TF_LOG=TRACE
```
