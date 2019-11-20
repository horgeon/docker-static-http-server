# Docker Static HTTP Server
Provide a lightweight image with a basic, static HTTP server.
It intends to provide only basic features. Supports no SSL, no directory listing and custom 404 page.

Mostly based on [Alex Edwards](https://www.alexedwards.net/)'s work.

## Configuration
The server is configurable using environment variables.

| Environment variable name | Default value | Function |
|--|--|--|
| ``HTTP_SERVER_LISTENING_IP`` | ``0.0.0.0`` | IP the server should listen on
| ``HTTP_SERVER_LISTENING_PORT`` | ``80`` | Port the server should listen on
| ``HTTP_SERVER_DIRECTORY`` | ``/static-data`` | Which directory to serve
| ``HTTP_SERVER_PREFIX`` | ``/static/`` | Which URL prefix the web server should consider

The server uses the ``404.html`` at the root of ``HTTP_SERVER_DIRECTORY`` as the custom 404 error page. If it is not present, it sends a generic 404 error message.
