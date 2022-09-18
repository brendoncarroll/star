# Star

Star is an opinionated Go library for writing CLI applications.

## Goals
- Make testing commands easy.
- Make working with parameters (args and flags) easy.
- No parsing or nullable flag/arg noise in the applications main flow.
- Separate the way a parameter is passed (via flag or positionally) from its use.
- Eliminate the need to use shared package level state.

## Design
Commands in Star are just metadata around functions.
Utilities are provided to create "directory" or "parent" commands which are common in modern CLI apps.
Parent commands take 1 argument and use it to lookup the name of a child command.

Command functions are of type `func(*star.Context) error`

The Context has references to `std{in, out, err}`, as well as methods to access any parameters the command is expecting.
Commands will never be executed without their requested parameters.

Retreiving a parameter from the `star.Context` will always succeed or always panic, regardless of the programs runtime input.
So panicing when accessing parameters is always a logic error, never a user error.
If the user forgets a parameter it will not panic and instead return an error without running any of the application logic.

