## Appctl
### Unified developer experience across all your projects

Do you write your own build scripts or do you prefer to use `make`? Are you too lazy to write your own Makefiles? do you even know what `make` is? Do you suffer of having several projects written with different technologies and you always struggle to remember how to build/test/package each one? 

No worries `appctl` knows well about `make` and Makefiles and leverages it for you so you 
have the simplest build chain accross all your projects, works for Java, Go,.. no matter what is the underline technology, it's actually indiferent to `appctl`. You just specify what kind of app you are working with and `appctl` manages it for you so you don't have to remember the exact commands to build each one of your projects.

You could just write Makefiles following your own guidelines in all your projects and you would already have a unified developer experience. But because you are too lazy and you won't do that here is `appctl` to do that for you.

I hope you like it

### Usage

First install the binary
```
make install
```

Then you can use it :)

Examples:

```
appctl build examples/simple-app/
```

or


```
cd examples/simple-app

appctl build
```

https://www.gnu.org/software/make/manual/make.html

