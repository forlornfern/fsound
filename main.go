package main

import (
	"fsound/cmd"

	"charm.land/log/v2"
	vlc "github.com/adrg/libvlc-go/v3"
)

func main() {
	if err := vlc.Init("--no-video"); err != nil {
		log.Fatal(err)
	}
	cmd.Exec()
	vlc.Release()
}
