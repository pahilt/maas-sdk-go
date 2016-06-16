# Command line options

## Required:

 - `client-id` - the client id, registered in Miracl OIDC provider
 - `client-secret`- the corresponding client secret
 - `redirect` - the registered redirect URL

## Optional
 - `addr` - Host to bind and port to listen on in the form host:port; the default is
 ":8002" which means bind all available interfaces and listen on port 8002
 - `templates-dir` - Folder holding the templates - absolute or relative to binary;
 the default is "templates", relative to binary and it should be functional.
 - `debug` - Print more information to stdout and stderr.
