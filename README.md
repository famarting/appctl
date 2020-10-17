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

## How does it work?

Appctl basically executes predefined Makefiles in your projects, Makefiles that just run whatever commands are needed to perform some task in your development process. Appctl defines this concepts:

### Recipes
A recipe is a list of tasks to be executed sequentially by Appctl. Developers can define recipes in the `app.yaml` file, other kind of recipes are also defined in `Templates`.

### Tasks
One task can execute one or more pre-defined recipes from one `Template`

### Templates

A template defines a list of recipes. But this recipes are special, they consist of the execution of a Makefile, i.e: [This is the docker template](docs/catalog/v1/docker/). As you can see the `index.json` file defines a list of "recipes" and each one points to a Makefile stored alongside the index.json file. [This](examples/simple-app/app.yaml) is an example showing the usage of recipes defined in templates.

There is a special template called `appctl` which invokes other appctl recipes on applications found in subdirectories of the current application directory. An example of this can be found on the [microservices example](examples/microservices-example/app.yaml)

### Applications

Managed with `app.yaml` files that can be created with `appctl init -n <your-app-name>`. Applications can use one or more templates and inherit all of their recipes, [here](examples/microservices-example/api-gateway/app.yaml) an example of an application inheriting all the recipes from the "nodejs" and "docker" templates. You can see the inherited recipes executing `appctl status` into the application directory.

`appctl status` will display a list with the available recipes that the current application has.

A recipe can be executed in an application by running:
```
appctl <recipe-name>
```

### Templates catalog

Templates allow to manage recipes and Makefiles and they are made available to you with an online [catalog](docs/catalog/v1). Because the catalog is just some [static files](https://famartinrh.github.io/appctl/catalog/v1/docker/) available through an http-server, currently the catalog is hosted in github pages and managed in this repository (under the `docs/catalog` folder)

Appctl will download the templates used in your applications into the `~/.appctl/templates/` directory