# MongoDB Atlas CBI Upload

## What it does

This policy will poll the MongoDB Atlas API, download the invoices from the last *n* months, convert
the invoices to the CBI CSV format, and upload them to the CBI Bill Upload API.

## Input Parameters

This policy has the following input parameters required when launching the policy:

- *MongoDB Atlas Organization ID* - The ID of the MongoDB Atlas organization for downloading invoices
- *Months* - The number of months back for downloading invoices
- *Bill Connect ID* - The ID of the CBI Bill Connect for uploading bills (format: `cbi-oi-optima-*`)
- *Email Addresses* - The list of email addresses to notify

## Policy Actions

The following policy actions are taken on any resources found to be out of compliance:

- Upload Bill to Optima endpoint
- Send an email report

## Prerequisites

This policy uses [credentials](https://docs.flexera.com/flexera/EN/Automation/CredentialsIntro.htm)
for connecting to the cloud API -- in order to apply this policy you must have a credential
registered in the system that is compatible with this policy. If there are no credentials listed
when you apply the policy, please contact your cloud admin and ask them to register a credential
that is compatible with this policy. The information below should be consulted when creating the
credentials.

### Credentials

- MongoDB Atlas API key with *Organization Member* and *Organization Read Only* permissions
  - Credential Type: *Digest Auth*
  - Username: value of MongoDB Atlas *public key*
  - Password: value of MongoDB Atlas *private key*
  - Provider: `mongodb_atlas`
- Flexera One client credentials (for Service Account) or refresh token with `csm_bill_upload_admin` role
  - Credential Type: *OAuth2*
  - Grant Type: *Client Credentials* or *Refresh Token*
  - Token URL: `https://login.flexera.com/oidc/token` or `https://login.flexera.eu/oidc/token`
  - When grant type is *Client Credentials*:
    - Client Authentication Method: *Client ID & Secret*
    - Client ID: value of Flexera One Service Account client ID
    - Client Secret: value of Flexera One Service Account client secret
  - When grant type is *Refresh Token*:
    - Client Authentication Method: *Token*
    - Token: value of Flexera One refresh token
  - Provider: `flexera`

## Supported Clouds

- MongoDB Atlas

## Cost

This Policy Template does not launch any instances, and so does not incur any cloud costs.
