the go dic
==========

# Golang Dependency Injection Contaner

To manage the responsibility of dependency creation, each Go application should have a service locator that is responsible for construction and lookup of dependencies.

Asking for dependencies solves the issue of hard coding, but it also means that the dic needs to be passed throughout the application. Passing the injector breaks the Law of Demeter. 

To remedy this, you could inject to a service only the dependency it needs.

### Benefits

- More **Reusable** Code
- More **Testable** Code
- More **Readable** Code

## Example...

**see** [example](./example/example.go) folder, you could run with 

**run** `go run example/example.go`
and  `go run example/example_struct.go`

### Injecting deps on structs

``` go
type mailer struct {
	Logger string `dic:"log"`
	Transport Transport `dic:"transport.sendmail"`
}
mailer = new(mailer)
cnt := container.New()

cnt.Register("output_writer", os.Stdout)
cnt.Register("log", log.New, reference.New("output_writer"), "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)

cnt.inject(mailer)
mailer.Transport.Send()
mailer.Logger.Println("hello log")
``` 

### Services

Injection of dependencies into controllers

``` go
cnt := container.New()
cnt.Register("logger", log.New /*, lot of arguments*/)
cnt.Register("template.path", "template/")
cnt.Register("template", abcTemplate.new, reference.New("template.path"), reference.New("logger"))

cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)
cnt.Register("mailer", MailerNew, reference.New("transport"), "[golangit] ")

// injecting into controller
cnt.Register("controller_home", func (logger Logger, mailer Mailer) {/*do something*/}, reference.New("logger"), reference.New("mailer"))

cnt.Register("controller_admin", func (tpl Template, mailer Mailer) {/*do something*/}, reference.New("template"), reference.New("mailer"))
``` 
### Registering Structs and parameters

``` go
type TestStruct struct {
	dbName string
	logger Logger    
}

cnt := container.New()
// Storing pareter
cnt.Register("log.filename", "error.log")
cnt.Register("db.name", "logger-database")
// storing service
cnt.Register("logger", LoggerNew, reference.New("log.filename"))
cnt.Register("context", &TestStruct{}, reference.New("db.name"), reference.New("logger"))

test := cnt.Get("context").(TestStruct)
// now struct has .dbName = 'logger-database' and the logger is an object wht the filename injected
``` 

## API

`Register` stores functions, structs or parameters into the service locator.

`Get` resolves dependencies and return the value.

`Inject` 

Everything is injected with lazy injection.

### Easy multiple env.

``` go
cnt := container.New()
if env == "prod" {
	cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)
} else {
   cnt.Register("transport.sendmail", StubMailNew)
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

### Injecting using tagging (annotation)

``` go
cnt.Register("mail.prefix", "[golangit] ")

type mailer struct {
	Sender     MailSender `dic:"transport.sendmail"`
	MailPrefix string     `dic:"mail.prefix"`
}

mailer := &mailer
cnt.Inject(mailer) //mailer has now dependencies injected
mailer.Sender.Send()
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

[Symfony dependency injection container](http://symfony.com/doc/current/components/dependency_injection)

[Spring service container](http://projects.spring.io/spring-framework)

@mikespook [bind function](https://bitbucket.org/mikespook/golib/src/27c65cdf8a772c737c9f4d14c0099bb82ee7fa35/funcmap/funcmap.go?at=default)

## ToDo

1. alias
2. Cli for debugging
3. godoc
4. improving injection