# GRPC Server

This package provides additional functionality on top of the default gRPC implementation.

## CMD

Folder `cmd` provides an example of two gRPC services using this library. This folder also includes examples on how to test the controllers of your gRPC services.

## Errors

This package contains helper functions for using errors in your gRPC services. 

Every error function accepts a message and (multiple) optional details. All available detail types are available in [this](https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails) godoc.

Example usage: 
```go
return grpcerrors.ErrNotFound("requested resource not found")
```

With details:
```go
return nil, grpcerrors.ErrNotFound("requested resource not found", &errdetails.ErrorInfo{
    Reason: "resource not found reason",
    Domain: "grpcerrors",
})
```

These details can be extracted when using a go gRPC client, see the example services for an example. 

Make sure to include the correct proto files inside your proto definition when using the codes, statuses or details. 

```proto
import "google/rpc/code.proto";
import "google/rpc/status.proto";
import "google/rpc/error_details.proto";
```