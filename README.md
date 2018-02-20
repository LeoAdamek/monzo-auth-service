Monzo Auth Service
==================

Monzo Auth service is a small service which allows
desktop applications to obtain an OAuth access token in a platform agnostic way.

Deployment
----------

Because Monzo OAuth apps are (currently) limited to your own account and consenting collaborators, you'll probably need to spin up your own service.

### Prerequisites:

You will need:
  * The AWS CLI
  * [AWS Sam Local](https://github.com/awslabs/aws-sam-local) installed
  * An AWS Account
  * A Monzo OAuth app [Create one here](https://developers.monzo.com/api)
  
### Preparation:

You will need to edit the following values in `template.yaml`

  * `RETURN_URL` -- Change to your own domain
  * `OAuthClientID` -- Change to your own OAuthClientID
  * `OAuthClientSecret` -- Leave this for now, we'll need to deploy the stack first to create the encryption key needed to set this.
  * The KMS `KeyPolicy` -- Change the administration users to your own AWS account. **This is important -- The AWS API should stop you forgetting this, but if you do, you will be uanble to use your key, and I will.**
  * Create an S3 bucket to deploy your app to, I use `aws-sam-MY_AWS_ACCOUNT_ID`
    `aws s3 mb s3://aws-sam-MY_AWS_ACCOUNT_ID`
  
### Build & Deploy

Build the binaries:

  $ make clean build
  
Upload the package:

  $ sam package --template-file template.yaml --s3-bucket $YOUR_S3_BUCKET --output-template-file packaged.yaml
  
Deploy:

  $ sam deploy --template-file ./packaged.yaml --stack-name monzo-auth-service
  
Wait a few minutes, it may take a while on first run.

Encrypt your OAuth Client Secret:

  $ aws kms encrypt --key-id generated_kms_key_id --plaintext $YOUR_MONZO_OAUTH_SECRET
  
Set the `OAuthClientSecret` to the value of the `CipherTextBlob` and then re-run the package and deploy steps.


Client Usage
------------

This section covers how to use the service to get request tokens.

### Request a token:

````http request
GET https://monzo-auth.adamek.io/new
````

Response will contain a token request/session ID:

````json
{"id": "TOKEN_REQ_ID", "login_url": "a_url_here"}
````

The user should be directed to open the url given by `login_url` to authenticate.

### Poll for the Token

You can then poll the service to see if the token is available:

````http request
GET https://monzo-auth.adamek.io/token?id=TOKEN_REQ_ID
````

If the user has authenticated and been redirected to get their token, you should get the following:

````json
{
"token": "your_token",
"expires": "token_expiry"
}
````

**Important**: Tokens are _not_ stored in the service. The single-use Authorization code is held
until a token reterevial request is received, at which point it is exchanged with Monzo for an access token which is
passed directly to the caller.
