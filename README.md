# go-godot-poc
A Proof of Concept for a multiplayer system using Godot on the frontend and Go on the backend.

## Getting Set Up

### Install Protoc

Protoc is a protocol buffer tool from Google that we will use to formulate our packets. Install it on your system thusly:

#### Linux/Unix

1. Download the latest release [here](https://github.com/protocolbuffers/protobuf/releases)
2. Unzip the archive
3. Put the `bin/protoc` binary into `usr/local/bin` (or somewhere else on your path)
4. Put the `/includes` folder contents into `usr/local/include` (or somewhere else on your path)
5. Test to ensure it worked: `protoc --version`. When set up (on Mac) the output was `libprotoc 29.3`.

#### Windows

Instructions TBD when somebody tries this on Windows. Try following [this tutorial](https://www.tbat.me/2024/11/09/godot-golang-mmo-part-1) and then update this README with a PR when you figure it out.
