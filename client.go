package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

//ClientID is the 20 byte id of our client
var ClientID = "DONDESTALABIBLIOTECA"

//ProtoName is the BitTorrent protocol we are using
var ProtoName = "BitTorrent protocol"

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Illegal USAGE!\n USAGE : ./Bittorrent <torrent_file> <output file>")
		return
	}
	torrentFile := os.Args[1]
	fileName := os.Args[2]

	torrent, err := NewTorrent(torrentFile)
	if err != nil {
		log.Fatal("Unable to decode the torrent file\n", err)
	}

	// create a new tracker and receive the list of peers
	hash := torrent.InfoHash()
	iDict := torrent.InfoDict()

	// Tracker connection
	tkInfo := NewTracker(hash, torrent, &iDict)
	peerList, _ := tkInfo.Connect()
	fmt.Printf("%v\n", peerList)

	//Start peer download
	tInfo := TorrentInfo{
		TInfo:        &iDict,
		ClientID:     ClientID,
		ProtoName:    ProtoName,
		ProtoNameLen: len(ProtoName),
		InfoHash:     string(hash[:len(hash)]),
	}

	PeerDownloader := NewPeerDownloader(tInfo, peerList, fileName)

	// keep announcing to tracker at minInterval
	//ticker := time.NewTicker(time.Second * time.Duration(minInterval))
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for _ = range ticker.C {
			tkInfo.Uploaded, tkInfo.Downloaded, tkInfo.Left =
				PeerDownloader.getProgress()
			tkInfo.sendGetRequest("")
		}
	}()
	PeerDownloader.StartDownload()
	ticker.Stop() // ticker is done
	// Send event stopped message to tracker
	tkInfo.Disconnect()
}
