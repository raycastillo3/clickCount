# Apache Datahub

~ I built *Apache Datahub* a click-tracking software that allows ***Data Analysts***, ***Marketing Teams***, and ***Web Developers*** to keep track of clicks on a website. This is a backend project.

Tracking user clicks on an e-commerce site is crucial for engineers as it provides invaluable insights into user behavior, enabling the optimization of website design and functionality. This data helps identify popular products, enhance user experience, and refine the overall site navigation, leading to increased customer satisfaction and higher conversion rates. By understanding how users interact with the platform, engineers can make informed decisions to tailor the online shopping experience and boost business success.

Apache Datahub has two main components. The primary component will be referred to as the "database", and it keeps track of a global count of `items`, `add to cart`, and `buy` button clicks. It also hosts the GRPC server, which web servers can use to get and set values for `items`, `add to cart`, and `buy` button clicks.

The other component, referred to as the "webserver", consists of three important parts. First, it serves a simple webpage with `item`, `add to cart`, and `buy buttons`, which allows us to test clicks from visitors/ users. Second, it hosts an HTTP API which the webpage reaches out to in order to update click counts. Finally, it runs a GRPC client, which allows the client to get and update the global count of the central server.

This project supports many webservers against the database. The webserver will maintain a local cache of click counts and periodically sync these values with the database, to precisely collect clicks with low load.

----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

### To build the binary:

```
go build -o clickCountApp -race -v .
```

### To run the database:

```
./clickCountApp database
```

You can specify an `--rpc-addr <addr>` flag to set the RPC server address the
database listens for requests on to something other than ":8080".

### To run a webserver:

```
./clickCountApp webserver
```

You can specify an `--rpc-addr <addr>` flag to set the address to connect to
the database on (it defaults to "localhost:8080").

You can specify an `--http-addr <addr>` (or `-a`) flag to set the address to
serve the website on. This is used to run multiple web servers, like:

```
./clickCountApp webserver -a :3001
./clickCountApp webserver -a :3002
./clickCountApp webserver -a :3003
```

Note that you can also run the database and webserver without the `build` step
by running the `main.go` file directly:

## Runing `main.go` file directly:

```
go run main.go database
go run main.go webserver -a :3002
```

To access the webserver's frontend, go to http://localhost:3000 in a browser (if
you used `--http-addr` or `-a` to specify a different one, use that port
instead of `3000`).

## Generating Protobufs

You will need to install the protobuf compiler:

https://grpc.io/docs/protoc-installation/

As well as some additional dependencies:

```
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

To re-generate the protobuf code after changing `pb/clickCountApp.proto`, run:

```
go generate ./pb/...
```
----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

### Example of how the front end would look like: 
<img width="1201" alt="Screen Shot 2023-11-28 at 12 54 22 AM" src="https://github.com/raycastillo3/clickCountApp/assets/60750068/89fca51e-c48d-4a54-a301-4c7abba5fe9b">

