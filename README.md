Logging Service
================

A logging library around [go-logging](https://github.com/op/go-logging) 
and using [lumberjack.v2](http://gopkg.in/natefinch/lumberjack.v2) for rotational file configuration.

You can use dep to get the dependencies if you want to work on the file

## Installing

### Using *go get*

    $ go get github.com/oaStuff/logservice

After this command *logservice* is ready to use. Its source will be in:

    $GOPATH/src/pkg/github.com/oaStuff/logservice

You can use `go get -u` to update the package.
You could also *dep* and the source would be avialable in the vendor folder of your project.


## Example

Let's have a look at an example 

```go
package main

import (
    "github.com/oaStuff/logservice"
)


func main() {
	logger := logger.Neww(logger.LoggerConfig{Enabled:true, AllowFileLog:true, AllowConsoleLog:true})
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Critical("Critical message")
	logger.Error("Error message")
}
```

## Explanation:
using the above code logging will happen at the console and in a file automatically created for you in a folder
called *logs* within your application folder.

you can allow/disable  file logging as well as console logging. You could also disable logging as
a whole in the application using `logger.LoggerConfig{}`

### NOTE:

you **MUST** call `logger.New()` to get a logger for your application.
Calling `logger.New()` with both `AllowFileLog` and `AllowConsoleLog` in 
`logger.LoggerConfig{}` set to false will automatically disable logging totally.