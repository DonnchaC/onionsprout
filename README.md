# OnionSprout

OnionSprout is a tool to run publicaly-accessible web services, for example from Raspberry Pi in your home, without a public IP. Additionally because traffic is routed via the Tor network you home location and IP address is kept private.

The tool involved a publicly accessible server, the OnionGateway, which transparently routes TLS
connections to your web server over a Tor onion service. Your onion service will have a public domain name and will be accessible by anyone without using Tor. All connections are end-to-end encrypted from the user to your onion service.



## Installation

The tool can install by cloning the repo and building the go package.

```
git get github.com/DonnchaC/onionsprout
```

To run you will need to have Tor running on the computer where you will host your server. OnionSprout will automatically provison a Tor onion server.

## Usage

Currently this tool sets up a local webserver which will automatically request a LetsEncrypt certifcate from the first connection. The webserver will show a test page.

```
DOMAIN=yoursubdomain.example-gateway.com onionsprout
```

## TODO

The tool should have a CLI to:

 - Register and configure a new subdomain or domain
 - Start a new onion address and store the private key and configuration
 - Configure the tool to terminate TLS connections and forward to the local service.

## Example CLI

The public domain will run the OnionGateway proxy. You can host it yourself or use a service run by a third party. A third party server will have it's own interface for registering your subdomain.

```
$ onionsprout init
Enter public domain: yoursubdomain.example-gateway.com
Enter token: XXXXXXXXXXXXXXXXXXXXXXXXX
Onion address private key (leave blank to generate):
Generating onion address.....
Generated onion address: uih5owv3h5huiyyjf7rnlkimpmh3hz2qdmqjqyvazcbau2lucjh3woyd.onion
Destination server (plaintext request will be forwarded here): localhost:8080
Stored new configutation at /etc/onionsprout/yoursubdomain.example-gateway.com.yml
```

Later OnionSprout can be started as a server:

```
$ onionsprout
Loading onionsprout configurations from /etc/onionsprout/
Started client service onion service on uih5owv3h5huiyyjf7rnlkimpmh3hz2qdmqjqyvazcbau2lucjh3woyd.onion forwarding to localhost:8080
...
```

or

```
onionsprout yoursubdomain.example-gateway.com
Loading onionsprout configurations from /etc/onionsprout/yoursubdomain.example-gateway.com
```