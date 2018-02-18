Monzo Auth Service
==================

Monzo Auth service is a small service which allows
desktop applications to obtain an OAuth access token in a platform agnostic way.

Client Usage
------------

This section covers how to use the service to get request tokens.

### Request a token:

````http request
GET https://mzauth.breakerofthings.solutions/request
````

Response will contain a token request/session ID:

````json
{"request_id": "TOKEN_REQ_ID", "login_url": "a_url_here"}
````

The user should be directed to open the url given by `login_url` to authenticate.

### Poll for the Token

You can then poll the service to see if the token is available:

````http request
GET https://mzauth.breakerofthings.solutions/token?id=YOUR_TOKEN_ID
````

If the user has authenticated and been redirected to get their token, you should get the following:

````json
{
"token": "your_token",
"expires": "token_expiry"
}
````

**Important**: Once a token has been presented, it is destroyed from the auth service.