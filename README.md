dic
===

#the golang dependency injection contaner

#### This project is not finished yet (almost started).

### We need a proper Go dependency injection container.

To manage the responsibility of dependency creation, each Go application should have a service locator that is responsible for construction and lookup of dependencies.

Asking for dependencies solves the issue of hard coding, but it also means that the dic needs to be passed throughout the application. Passing the injector breaks the Law of Demeter. 

To remedy this, you could inject to a service only the dependency it needs.

### Benefits

- More **Reusable** Code
- More **Testable** Code
- More **Readable** Code

## Example...

see [example](./example/example.go) folder

``` go
cnt := container.New()
cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)
cnt.Register("transport.postfix", PostfixNew, "different args than sendmail")
cnt.Alias("transport", "transport.sendmail")
cnt.Register("mailer", MailerNew, reference.New("transport"), "[golangit] ")
// and then in your code
cnt.Get("mailer").(Mailer).Send("liuggio")
```

### Easy multiple env.

``` go
cnt := container.New()
if env == "prod" {
	cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)
} else {
   cnt.Register("transport.stub", echoMailNew)
}
```

### Create a not public service

Is possible to create a service usable only as dependency

``` go
def := definition.New(SendmailNew).setPublic(false)
cnt.Register("transport.sendmail", def)
```

### Create a not static service

All the service are served statically, this means that the service
is executed only once.

``` go
def := definition.New(SendmailNew).setStatic(false)
cnt.Register("transport.sendmail", def)
```

### *Injecting using tagging (annotation)

Not implemented

``` go
type mailer struct {
	Sender     MailSender `dic:Mailer type:service`
	MailPrefix string `dic:mailPrefix type:parameter`
}

cnt.Map("mailer", mailer)
```


## Test

This library has been developed with [ginkgo a BDD Testing Framework for Go](http://onsi.github.io/ginkgo),

Install the dependencies:

	go get -v -t ./...
  	go get -v github.com/onsi/ginkgo
  	go get -v github.com/onsi/gomega
  	go install -v github.com/onsi/ginkgo/ginkgo

Run all  the bdd tests:

	ginkgo -r

## License 

[MIT License](LICENSE)

## Thanks to

@mikespook for bind function
https://bitbucket.org/mikespook/golib/src/27c65cdf8a772c737c9f4d14c0099bb82ee7fa35/funcmap/funcmap.go?at=default

## ToDo

1. Paramters
2. tagging
3. aliasing
4. Cli for debugging
5. Setter injection, Property injection

1. Dependency that returns multiples values.