package goaviatrix

import (
	"fmt"
	"encoding/json"
	"errors"
	"log"
)

type Account struct {
	CID                         	string `form:"CID,omitempty"`
	Action                  		string `form:"action,omitempty"`
	AccountName             		string `form:"account_name,omitempty" json:"account_name,omitempty"`
	AccountPassword         		string `form:"account_password,omitempty" json:"account_password,omitempty"`
	AccountEmail            		string `form:"account_email,omitempty" json:"account_email_addr,omitempty"`
	CloudType               		int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AwsAccountNumber        		string `form:"aws_account_number,omitempty" json:"account_number,omitempty"`
	AwsIam                  		string `form:"aws_iam,omitempty" json:"aws_iam,omitempty"`
	AwsAccessKey            		string `form:"aws_access_key,omitempty" json:"account_access_key,omitempty"`
	AwsSecretKey            		string `form:"aws_secret_key,omitempty" json:"account_secret_access_key,omitempty"`
	AwsRoleApp              		string `form:"aws_role_arn,omitempty" json:"aws_role_arn,omitempty"`
	AwsRoleEc2              		string `form:"aws_role_ec2,omitempty" json:"aws_role_ec2,omitempty"`
	AzureSubscriptionId     		string `form:"azure_subscription_id,omitempty" json:"azure_subscription_id,omitempty"`
	ArmSubscriptionId       		string `form:"arm_subscription_id,omitempty" json:"arm_subscription_id,omitempty"`
	ArmApplicationEndpoint  		string `form:"arm_application_endpoint,omitempty" json:"arm_ad_tenant_id,omitempty"`
	ArmApplicationClientId      	string `form:"arm_application_client_id,omitempty" json:"arm_ad_client_id,omitempty"`
	ArmApplicationClientSecret  	string `form:"arm_application_client_secret,omitempty" json:"arm_ad_client_secret,omitempty"`
	AwsgovAccountNumber       		string `form:"awsgov_account_number,omitempty" json:"awsgov_account_number,omitempty"`
	AwsgovAccessKey       			string `form:"awsgov_access_key,omitempty" json:"awsgov_access_key,omitempty"`
	AwsgovSecretKey       			string `form:"awsgov_secret_key,omitempty" json:"awsgov_secret_key,omitempty"`
	AwsgovCloudtrailBucket      	string `form:"awsgov_cloudtrail_bucket,omitempty" json:"awsgov_cloudtrail_bucket,omitempty"`
	AzurechinaSubscriptionId    	string `form:"azurechina_subscription_id,omitempty" json:"azurechina_subscription_id,omitempty"`
	AwschinaAccountNumber       	string `form:"awschina_account_number,omitempty" json:"awschina_account_number,omitempty"`
	AwschinaAccessKey       		string `form:"awschina_access_key,omitempty" json:"awschinacloud_access_key,omitempty"`
	AwschinaSecretKey       		string `form:"awschina_secret_key,omitempty" json:"awschinacloud_secret_key,omitempty"`
	ArmChinaSubscriptionId      	string `form:"arm_china_subscription_id,omitempty" json:"arm_china_subscription_id,omitempty"`
	ArmChinaApplicationEndpoint 	string `form:"arm_china_application_endpoint,omitempty" json:"arm_china_application_endpoint,omitempty"`
	ArmChinaApplicationClientId     string `form:"arm_china_application_client_id,omitempty" json:"arm_china_application_client_id,omitempty"`
	ArmChinaApplicationClientSecret string `form:"arm_china_application_client_secret,omitempty" json:"arm_china_application_client_secret,omitempty"`
}

type AccountResult struct {
	AccountList        []Account `json:"account_list"`
}

type AccountListResp struct {
	Return  bool   `json:"return"`
	Results AccountResult `json:"results"`
	Reason  string `json:"reason"`
}


func (c *Client) CreateAccount(account *Account) (error) {
	account.CID=c.CID
	account.Action="setup_account_profile"
	resp,err := c.Post(c.baseURL, account)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetAccount(account *Account) (*Account, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_accounts", c.CID)
	resp,err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}
	var data AccountListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	acclist:= data.Results.AccountList
	for i := range acclist {
		log.Printf("[TRACE] %s", acclist[i].AccountName)
		if acclist[i].AccountName == account.AccountName {
			log.Printf("[INFO] Found Aviatrix Account %s", account.AccountName)
			return &acclist[i], nil
		}
	}
	log.Printf("Couldn't find Aviatrix account %s", account.AccountName)
	return nil, ErrNotFound
}

func (c *Client) UpdateAccount(account *Account) (error) {
	account.CID=c.CID
	account.Action="edit_account_profile"
	resp,err := c.Post(c.baseURL, account)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

//There is no Aviatrix account user resource currently to keep things minimal.
//Since creating an Aviatrix account by default creates a user of the same name,
//we'll only be able to configure that default user. Also not creating a new structure
//for Aviatrix user, and marshalling received data to JSON manually, since we
//need Aviatrix user data only for this method.
func (c *Client) UpdateAccountUser(what, username, oldpass, newpass, email string) (error) {
	account := make(map[string]interface{})
	account["CID"] = c.CID
	account["action"] = "edit_account_user"
	account["username"] = username
	account["what"] = what
	if what == "password" {
		account["old_password"] = oldpass
		account["new_password"] = newpass
	}
	if what == "email" {
		account["email"] = email
	}
	log.Printf("[INFO] Parsed Aviatrix account: %#v", account)
	resp,err := c.Post(c.baseURL, account)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteAccount(account *Account) (error) {
	path := c.baseURL + fmt.Sprintf("?action=delete_account_profile&CID=%s&account_name=%s", c.CID, account.AccountName)
	resp,err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}
