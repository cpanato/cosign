//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"

	"github.com/sigstore/cosign/v2/cmd/cosign/cli/manifest"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/verify"
	"github.com/spf13/cobra"
)

func Manifest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Provides utilities for discovering images in and performing operations on Kubernetes manifests",
	}

	cmd.AddCommand(
		manifestVerify(),
	)

	return cmd
}

func manifestVerify() *cobra.Command {
	o := &options.VerifyOptions{}

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify all signatures of images specified in the manifest",
		Long: `Verify all signature of images in a Kubernetes resource manifest by checking claims
against the transparency log.`,
		Example: `  cosign manifest verify --key <key path>|<key url>|<kms uri> <path/to/manifest>

  # verify cosign claims and signing certificates on images in the manifest
  cosign manifest verify <path/to/my-deployment.yaml>

  # additionally verify specified annotations
  cosign manifest verify -a key1=val1 -a key2=val2 <path/to/my-deployment.yaml>

  # verify images with public key
  cosign manifest verify --key cosign.pub <path/to/my-deployment.yaml>

  # verify images with public key provided by URL
  cosign manifest verify --key https://host.for/<FILE> <path/to/my-deployment.yaml>

  # verify images with public key stored in Azure Key Vault
  cosign manifest verify --key azurekms://[VAULT_NAME][VAULT_URI]/[KEY] <path/to/my-deployment.yaml>

  # verify images with public key stored in AWS KMS
  cosign manifest verify --key awskms://[ENDPOINT]/[ID/ALIAS/ARN] <path/to/my-deployment.yaml>

  # verify images with public key stored in Google Cloud KMS
  cosign manifest verify --key gcpkms://projects/[PROJECT]/locations/global/keyRings/[KEYRING]/cryptoKeys/[KEY] <path/to/my-deployment.yaml>

  # verify images with public key stored in Hashicorp Vault
  cosign manifest verify --key hashivault://[KEY] <path/to/my-deployment.yaml>`,
		Args:             cobra.ExactArgs(1),
		PersistentPreRun: options.BindViper,
		RunE: func(cmd *cobra.Command, args []string) error {
			annotations, err := o.AnnotationsMap()
			if err != nil {
				return err
			}
			v := &manifest.VerifyManifestCommand{
				VerifyCommand: verify.VerifyCommand{
					RegistryOptions:              o.Registry,
					CertVerifyOptions:            o.CertVerify,
					CheckClaims:                  o.CheckClaims,
					KeyRef:                       o.Key,
					CertRef:                      o.CertVerify.Cert,
					CertGithubWorkflowTrigger:    o.CertVerify.CertGithubWorkflowTrigger,
					CertGithubWorkflowSha:        o.CertVerify.CertGithubWorkflowSha,
					CertGithubWorkflowName:       o.CertVerify.CertGithubWorkflowName,
					CertGithubWorkflowRepository: o.CertVerify.CertGithubWorkflowRepository,
					CertGithubWorkflowRef:        o.CertVerify.CertGithubWorkflowRef,
					CertChain:                    o.CertVerify.CertChain,
					IgnoreSCT:                    o.CertVerify.IgnoreSCT,
					SCTRef:                       o.CertVerify.SCT,
					Sk:                           o.SecurityKey.Use,
					Slot:                         o.SecurityKey.Slot,
					Output:                       o.Output,
					RekorURL:                     o.Rekor.URL,
					Attachment:                   o.Attachment,
					Annotations:                  annotations,
					LocalImage:                   o.LocalImage,
					Offline:                      o.CommonVerifyOptions.Offline,
					TSACertChainPath:             o.CommonVerifyOptions.TSACertChainPath,
					IgnoreTlog:                   o.CommonVerifyOptions.IgnoreTlog,
					MaxWorkers:                   o.CommonVerifyOptions.MaxWorkers,
				},
			}

			if o.CommonVerifyOptions.MaxWorkers == 0 {
				return fmt.Errorf("please set the --max-worker flag to a value that is greater than 0")
			}

			return v.Exec(cmd.Context(), args)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
