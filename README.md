dic
===

the golang dependency injection contaner

### This project is not finished yet (almost started).


# We need a proper Go dependency injection.

To manage the responsibility of dependency creation, each Go application should have a service locator that is responsible for construction and lookup of dependencies.

Asking for dependencies solves the issue of hard coding, but it also means that the dic needs to be passed throughout the application. Passing the injector breaks the Law of Demeter. 

To remedy this, you could inject to a service only the dependency it needs.

### Benefits

- More Reusable Code
- More Testable Code
- More Readable Code


## Implementation 

Please contribute to the [RFC a proper DIC](https://github.com/golangit/dic/issues/1).

## License 

[MIT License](LICENSE)