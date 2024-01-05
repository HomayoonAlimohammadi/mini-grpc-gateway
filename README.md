# Mini gRPC Gateway

Source code for my article: [Protoc Plugins in Go: gRPC-REST Gateway from Scratch](https://medium.com/@homayoonalimohammadi)

![gRPC-REST Gateway Request Flow](gRPC-REST-Gateway.png "gRPC-REST Gateway Request Flow")

## How to run:
- Make sure to have `Go` and `Protoc` installed.
```bash
make generate-backend
make generate
make backend
```
In a new terminal, run:
```bash
make run
```
Try it out with a simple curl:
```bash
curl localhost:8000/api/posts
```