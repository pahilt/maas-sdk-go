# maas-sdk-go

[![Master Build Status](https://secure.travis-ci.org/miracl/maas-sdk-go.png?branch=master)](https://travis-ci.org/miracl/maas-sdk-go?branch=master)
[![Master Coverage Status](https://coveralls.io/repos/miracl/maas-sdk-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/miracl/maas-sdk-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/miracl/maas-sdk-go)](https://goreportcard.com/report/github.com/miracl/maas-sdk-go)

* **category**:    SDK
* **copyright**:   2016 MIRACL UK LTD
* **license**:     ASL 2.0 - http://www.apache.org/licenses/LICENSE-2.0
* **link**:        https://github.com/miracl/maas-sdk-go

## Description

GO version of the Software Development Kit (SDK) for MPin-As-A-Service (MAAS).


# Setup

## Install

Get it with `go get github.com/miracl/maas-sdk-go`.

There are two packages there - `maas` (in the root folder) and `demo` (in folder `demo`).

Package `maas` is the SDK. To use it, add `import "github.com/miracl/maas-sdk-go"`.

Package `demo`, located in a subdirectory named correspondingly, is an example application
that uses the SDK.

## Tests

To run tests with `go test` in the root folder.

# MAAS API

## Details and usage

All interaction with API happens through `mass.Client` interface. Each
application needs to construct instance of `maas.Client`.


### Initialization
To start using MAAS API, `maas.Client` can be initialized with the
`maas.NewClient` function . This function accepts a single argument,
`maas.Config`. In general it is needed to populate the ClientID, ClientSecret and
RedirectURI parameters.

```
client, err := maas.NewClient(maas.Config{
        ClientID:     CLIENT_ID,
        ClientSecret: CLIENT_SECRET,
        RedirectURI:  REDIRECT_URL,
    })
```

`CLIENT_ID` and `CLIENT_SECRET` can be obtained from Miracl (unique per
application). `REDIRECT_URI` is URI of your application end-point that will be
responsible obtaining token. It should be the same as registered in Miracl
system for this client ID. `DISCOVERY_URL` is the provider discovery URL.

Please note that this initialization includes (at least one) network call to
discovery endpoint, so it is recommended to do it at a proper time.


### Authorization flow

Authorization flow depends on `mpin.js` browser library. `mpad.js` depends on
`jquery` and `bootstrap`.To use it, load it in the bottom of the `<body>` element of
page responsible for login, by populating AUTH_URL (see below):

```
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
<script src="/static/js/bootstrap.min.js"></script>
...
<script src="<<Insert correct mpad url here>>" data-authurl="{{ .AuthURL }}" data-element="btmpin"></script>
```

Please refer to your distributor-specific documentation to find the correct url for the mpad.js script src

If user is not authorized, use `client.GetAuthRequestURL(state)`
to get authorization request URL and set client internal state. Returned URL
should be used with `mpad.js` - AUTH_URL in the sample. After user
interaction with Miracl system  user will be sent to `RedirectURI` defined at
creation of `maas.Client`.

To complete authorization pass the authorization code (from the query string received on `redirect_uri`)
 to`client.ValidateAuth(code)`. This method will return access and identity token
and `nil` error if user denied authorization and token if authorization succeeded.

The replying party application should take care for additional checks for state at the OIDC
handler - the ValidateAuth method only check OIDC token validity.

User info can be retrieved using the client.GetUserUnfo(accessToken). This method returns
`maas.UserInfo` structure and error.


## Example

Pass `CLIENT_ID`, `CLIENT_SECRET` and `REDIRECT_URI` as command line options to example.
