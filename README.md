# promtail golang client

The purpose of this small library is creating a tool for upload logs to grafana loki. The logs looks pretty amazing and awesome, we can have static labels or dynamic labels for filtering while we're debugging. 

![awesome_logs](https://user-images.githubusercontent.com/3860869/178655633-505993ba-28ef-4a18-9639-642ba5b2c401.png)


promtail client for golang services

```golang
// create a promtail client
promtail, pErr := client.NewSimpleClient(host, username, password,
    client.WithStaticLabels(map[string]interface{}{
        "app": "awesome-service",
    }),
    client.WithStreamConverter(promtail.NewRawStreamConv(labels, "=")),
)
if pErr != nil {
    panic(pErr)
}

// now can be used as a Writer for any logger

// setup any logger and use as the output, can be combined with many outputs using io.MultiWriter
output := zerolog.ConsoleWriter{Out: promtail}
output.FormatMessage = func(i interface{}) string {
    _, ok := i.(string)
    if ok {
        return fmt.Sprintf("%-50s", i)
    } else {
        return ""
    }
}
output.FormatLevel = func(i interface{}) string {
    _, ok := i.(string)
    if ok {
        return fmt.Sprintf("level=%-7s", i)
    } else {
        return "level=info"
    }
}
log.Logger = log.Output(output)

log.Info().
    Str("my_label", "some_value").
    Msg("the message")

```

