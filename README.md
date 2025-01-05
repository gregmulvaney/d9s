# D9s
A shameless ripoff of the brilliant UI for [derailed/K9s](https://github.com/derailed/k9s) built for Docker. 

## Installation 
Clone the repo and  
```make build```

Currently connecting to the server relies on the `DOCKER_HOST` environment variable. Refer to the [Docker documentation](https://docs.docker.com/engine/daemon/remote-access/) on exposing the host daemon to remote access. The default example is a simple unauthenticated TCP connection, however more secure methods are available.
