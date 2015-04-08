# go-http-request-logger
Go based reverse proxy to capture and HTTP requests

# Usage

To see command line options;

$ ./gohttpreqlog -help


By default, without any parameters, this will bind and listen on localhost:9000 and will reverse proxy all traffic to localhost:9191.

$ ./gohttpreqlog > output.log &
