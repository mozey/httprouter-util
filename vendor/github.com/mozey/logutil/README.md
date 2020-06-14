# logutil

Utils for use with [zerolog](https://github.com/rs/zerolog)

## Quick start

Setup logger with console writer

    consoleWriter := true
    logutil.SetupLogger(consoleWriter)

Log error with stack trace

    err := errors.WithStack(fmt.Errorf("testing"))
    log.Error().Stack().Err(err).Msg("")

When logging to json a stack trace will also be included    

