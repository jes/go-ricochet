# GoRicochet [![Build Status](https://travis-ci.org/s-rah/go-ricochet.svg?branch=master)](https://travis-ci.org/s-rah/go-ricochet) [![Go Report Card](https://goreportcard.com/badge/github.com/s-rah/go-ricochet)](https://goreportcard.com/report/github.com/s-rah/go-ricochet) [![Coverage Status](https://coveralls.io/repos/github/s-rah/go-ricochet/badge.svg?branch=master)](https://coveralls.io/github/s-rah/go-ricochet?branch=master)

GoRicochet is an experimental implementation of the [Ricochet Protocol](https://ricochet.im)
in Go.

## Features

* A simple API that you can use to build Automated Ricochet Applications
* A suite of regression tests that test protocol compliance.

## Building an Automated Ricochet Application

Below is a simple echo bot, which responds to any chat message. You can also find this code under `examples/echobot`

    package main

    import (
      "github.com/s-rah/go-ricochet/application"
      "github.com/s-rah/go-ricochet/utils"
      "log"
      "time"
    )

    func main() {
      echobot := new(application.RicochetApplication)
      pk, err := utils.LoadPrivateKeyFromFile("./private_key")

      if err != nil {
        log.Fatalf("error reading private key file: %v", err)
      }

      l, err := application.SetupOnion("127.0.0.1:9051", "tcp4", "", pk, 9878)

      if err != nil {
        log.Fatalf("error setting up onion service: %v", err)
      }

      echobot.Init(pk, new(application.AcceptAllContactManager))
      echobot.OnChatMessage(func(rai *application.RicochetApplicationInstance, id uint32, timestamp time.Time, message string) {
        log.Printf("message from %v - %v", rai.RemoteHostname, message)
        rai.SendChatMessage(message)
      })
      log.Printf("echobot listening on %s", l.Addr().String())
      echobot.Run(l)
    }

Each automated ricochet service can extend of the `StandardRicochetService`. From there
certain functions can be extended to fully build out a complete application.

## Security and Usage Note

This project is experimental and has not been independently reviewed. If you are
looking for a quick and easy way to use ricochet please check out [Ricochet IM](https://ricochet.im).
