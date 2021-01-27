/*
Copyright © 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/cobra/aid"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/cobra/dao"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/helper"
)

var sshKeyValidArgs = []string{"list", "create", "update", "delete"}

// SSHKeyCmd command to display tecli current version
func SSHKeyCmd() *cobra.Command {
	man, err := helper.GetManual("sshKey")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:       man.Use,
		Short:     man.Short,
		Long:      man.Long,
		Example:   man.Example,
		ValidArgs: sshKeyValidArgs,
		Args:      cobra.OnlyValidArgs,
		PreRunE:   sshKeyPreRun,
		RunE:      sshKeyRun,
	}

	usage := `A name to identify the SSH key.`
	cmd.Flags().String("name", "", usage)

	usage = `The content of the SSH private key.`
	cmd.Flags().String("value", "", usage)

	return cmd
}

func sshKeyPreRun(cmd *cobra.Command, args []string) error {
	if err := helper.ValidateCmdArgs(cmd, args, "sshKey"); err != nil {
		return err
	}

	fArg := args[0]
	switch fArg {
	case "list":
		if err := helper.ValidateCmdArgAndFlag(cmd, args, "sshKey", fArg, "organization"); err != nil {
			return err
		}
	case "create", "read", "delete":
		if err := helper.ValidateCmdArgAndFlag(cmd, args, "sshKey", fArg, "organization"); err != nil {
			return err
		}

		if err := helper.ValidateCmdArgAndFlag(cmd, args, "sshKey", fArg, "name"); err != nil {
			return err
		}

		if err := helper.ValidateCmdArgAndFlag(cmd, args, "sshKey", fArg, "value"); err != nil {
			return err
		}
	}

	return nil
}

func sshKeyRun(cmd *cobra.Command, args []string) error {

	token := dao.GetTeamToken(profile)
	client := aid.GetTFEClient(token)

	var sshKey *tfe.SSHKey
	var err error

	fArg := args[0]
	switch fArg {
	case "list":
		list, err := sshKeyList(client, organization)
		if err == nil {
			if len(list.Items) > 0 {
				for _, item := range list.Items {
					fmt.Printf("%v\n", aid.ToJSON(item))
				}
			} else {
				return fmt.Errorf("no ssh key was found")
			}
		}
	case "create":
		options := aid.GetSSHKeysCreateOptions(cmd)
		sshKey, err = sshKeyCreate(client, options)

		if err == nil && sshKey.ID != "" {
			fmt.Println(aid.ToJSON(sshKey))
		}
	case "read":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		sshKey, err := sshKeyRead(client, name)
		if err == nil {
			fmt.Println(aid.ToJSON(sshKey))
		} else {
			return fmt.Errorf("sshKey %s not found\n%v", name, err)
		}
	case "update":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		options := aid.GetSSHKeysUpdateOptions(cmd)
		sshKey, err = sshKeyUpdate(client, name, options)
		if err == nil && sshKey.ID != "" {
			fmt.Println(aid.ToJSON(sshKey))
		} else {
			return fmt.Errorf("unable to update sshKey\n%v", err)
		}
	case "delete":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		err = sshKeyDelete(client, name)
		if err == nil {
			fmt.Printf("sshKey %s deleted successfully\n", name)
		} else {
			return fmt.Errorf("unable to delete sshKey %s\n%v", name, err)
		}
	}

	return nil
}

func sshKeyList(client *tfe.Client, organization string) (*tfe.SSHKeyList, error) {
	return client.SSHKeys.List(context.Background(), organization, tfe.SSHKeyListOptions{})
}

// Create is used to create a new sshKey.
func sshKeyCreate(client *tfe.Client, options tfe.SSHKeyCreateOptions) (*tfe.SSHKey, error) {
	return client.SSHKeys.Create(context.Background(), organization, options)
}

// Read a sshKey by its name.
func sshKeyRead(client *tfe.Client, sshKeyID string) (*tfe.SSHKey, error) {
	return client.SSHKeys.Read(context.Background(), sshKeyID)
}

// Update settings of an existing sshKey.
func sshKeyUpdate(client *tfe.Client, sshKeyID string, options tfe.SSHKeyUpdateOptions) (*tfe.SSHKey, error) {
	return client.SSHKeys.Update(context.Background(), sshKeyID, options)
}

// // Delete a sshKey by its name.
func sshKeyDelete(client *tfe.Client, sshKeyID string) error {
	return client.SSHKeys.Delete(context.Background(), sshKeyID)
}