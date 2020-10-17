# Appctl

### Unified developer experience across all your projects

Are you too lazy to write your own build scripts or your own Makefiles? Do you even know what `make` is? Do you suffer from having several projects written with different technologies and you always struggle to remember how to build/test/package each one?

This project aims to provide a simple and unified development experice across all your projects

## What does `appctl` do?

Appctl executes your most common workflows in your development process through a simple and unified CLI across all your projects

Actually this tool does nothing, your typical build tools do the job. You will use npm for nodejs apps, maven for java, docker for container images, ... Whatever tool, but `appctl` will invoke them for you.
You may even forget about all of that but only remember to execute `appctl build`

## Try it

To install from source, only pre-requisite is golang >= 1.13
```
make install
```

Or just download and install it with the `install.sh` script
```
curl -sfL https://raw.githubusercontent.com/famartinrh/appctl/master/install.sh | sh -
```

There are some examples of applications using `appctl` under the folder `examples`

#### Build, test and create a container image of Quarkus app (Java)

```
cd examples/simple-app

appctl status

appctl build

docker images | grep simple-app

cat app.yaml
```

#### Microservices application using nodejs

```
cd examples/microservices-example

appctl status

appctl build

docker images | grep appctl-

cat app.yaml

cat api-gateway/app.yaml
```

<!-- ## How does it work? -->



<!-- https://www.gnu.org/software/make/manual/make.html -->





