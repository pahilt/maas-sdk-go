# Demo MAAS RP

This is a sample application implemented with MAAS SDK. At minimum, command line options 
`client-id`, `client-secret`, `discovery` and `redirect` have to be provided.

## Building

The demo depends only on the SDK. make sure you have the sdk ("github.com/miracl/maas-sdk-go") available and build with "go build".


## Configuration options.

The demo application is configured only with command line options

* `-client-id string`  RP client ID, obtained from Miracl. Required.

* `-client-secret` RP client secret, obtained from Miracl. Required.

* `-discovery` MAAS authorization server discovery URL. Required.

* `-redirect` URL for the authorization server to redirects back. It should be the one registered at MAAS. Required.

* `-addr` IP and address to bind and listoen on, in the form "IP:PORT". Default is ":8002" (empty IP binds all available interfaces).

* `-static-dir` Location of the demo `static` directory. Default is "static", relative to executable.

* `-codepad-dir` Location of the demo `codepad` directory. Default is "codepad", relative to executable.

* `-templates-dir` Location of the demo `templates` directory. Default is "templates", relative to executable.

