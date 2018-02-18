Monzo Auth Service
==================

Monzo Auth service is a small service which allows
desktop applications to obtain an OAuth access token in a platform agnostic way.

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