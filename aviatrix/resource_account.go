package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"aws_account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_iam": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_role_app": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_role_ec2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:      d.Get("account_name").(string),
		AccountPassword:  d.Get("account_password").(string),
		AccountEmail:     d.Get("account_email").(string),
		CloudType:        d.Get("cloud_type").(int),
		AwsAccountNumber: d.Get("aws_account_number").(string),
		AwsIam:           d.Get("aws_iam").(string),
		AwsRoleApp:       d.Get("aws_role_app").(string),
		AwsRoleEc2:       d.Get("aws_role_ec2").(string),
		AwsAccessKey:     d.Get("aws_access_key").(string),
		AwsSecretKey:     d.Get("aws_secret_key").(string),
	}

	log.Printf("[INFO] Creating Aviatrix account: %#v", account)
	err := client.CreateAccount(account)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Account: %s", err)
	}
	d.SetId(account.AccountName)
	return nil
	//return resourceAccountRead(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)
	acc, err := client.GetAccount(account)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Aviatrix Account: %s", err)
	}
	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("account_email", acc.AccountEmail)
		d.Set("cloud_type", acc.CloudType)
		d.Set("aws_account_number", acc.AwsAccountNumber)
		//d.Set("aws_iam", acc.AwsIam)
		d.Set("aws_role_app", acc.AwsRoleApp)
		d.Set("aws_role_ec2", acc.AwsRoleEc2)
		d.Set("aws_access_key", acc.AwsAccessKey)
		d.Set("aws_secret_key", acc.AwsSecretKey)
		d.SetId(acc.AccountName)
	}
	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:      d.Get("account_name").(string),
		CloudType:        d.Get("cloud_type").(int),
		AwsAccountNumber: d.Get("aws_account_number").(string),
		AccountPassword:  d.Get("account_password").(string),
		AwsIam:           d.Get("aws_iam").(string),
		AwsRoleApp:       d.Get("aws_role_app").(string),
		AwsRoleEc2:       d.Get("aws_role_ec2").(string),
		AwsAccessKey:     d.Get("aws_access_key").(string),
		AwsSecretKey:     d.Get("aws_secret_key").(string),
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)
	d.Partial(true)
	if d.HasChange("aws_account_number") || d.HasChange("aws_access_key") || d.HasChange("aws_secret_key") || d.HasChange("aws_iam") || d.HasChange("aws_role_app") || d.HasChange("aws_role_ec2") {
		err := client.UpdateAccount(account)
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix Account: %s", err)
		}
		if d.HasChange("aws_account_number") {
			d.SetPartial("aws_account_number")
		}
		if d.HasChange("aws_access_key") {
			d.SetPartial("aws_access_key")
		}
		if d.HasChange("aws_secret_key") {
			d.SetPartial("aws_secret_key")
		}
		if d.HasChange("aws_iam") {
			d.SetPartial("aws_iam")
		}
		if d.HasChange("aws_role_app") {
			d.SetPartial("aws_role_app")
		}
		if d.HasChange("aws_role_ec2") {
			d.SetPartial("aws_role_ec2")
		}
	}
	if d.HasChange("account_password") {
		oldpass, newpass := d.GetChange("account_password")
		err := client.UpdateAccountUser("password", account.AccountName, oldpass.(string), newpass.(string), "")
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix Account User password: %s", err)
		}
		d.SetPartial("account_password")
	}
	if d.HasChange("account_email") {
		err := client.UpdateAccountUser("email", account.AccountName, "", "", d.Get("account_email").(string))
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix Account User email: %s", err)
		}
		d.SetPartial("account_email")
	}

	d.Partial(false)
	//return resourceAccountRead(d, meta)
	return nil
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)

	err := client.DeleteAccount(account)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Account: %s", err)
	}
	return nil
}
