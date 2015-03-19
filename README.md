## Overview

Package dogstatsd provides a Go DogStatsD client. DogStatsD extends StatsD - adding tags and histograms. The documentation for DogStatsD is here: http://docs.datadoghq.com/guides/dogstatsd/

## Get the code

    $ go get github.com/ooyala/go-dogstatsd

## Usage

    // Create the client
    c, err := dogstatsd.New("127.0.0.1:8125", "flubber.", []string{"us-east-1a"})
    if err != nil {
      log.Fatal(err)
    }
    defer c.Close()

    // Prefix every metric with the app name
    c.SetGlobalNamespace("flubber.")
    // Send the EC2 availability zone as a tag with every metric
    c.SetGlobalTags([]string{"us-east-1a"})

    err = c.Gauge("request.duration", 1.2, nil, 1)

    // Post info to datadog event stream
    err = c.Info("cookie alert", "Cookies up for grabs in the kitchen!", nil)

## Development

Run the tests with:

    $ go test -v

## Documentation

Please see: http://godoc.org/github.com/ooyala/go-dogstatsd

## License

go-dogstatsd is released under the [MIT license](http://www.opensource.org/licenses/mit-license.php).
