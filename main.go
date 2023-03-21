package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/FirstProjectFor/FPF_NET/util"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
)

const ProtocolChat = "p2p/chat"

func main() {
	go http.ListenAndServe(":8086", nil)

	ctx := context.Background()
	host8888 := host8888()
	host8888.SetStreamHandler(ProtocolChat, func(stream network.Stream) {
		for {
			lengthData := make([]byte, 8)
			_, err := io.ReadAtLeast(stream, lengthData, 8)
			util.PanicIfNotNil(err)

			dataLength := binary.LittleEndian.Uint64(lengthData)
			if dataLength == 0 {
				continue
			}

			contentData := make([]byte, dataLength)
			_, err = io.ReadAtLeast(stream, contentData, int(dataLength))
			util.PanicIfNotNil(err)

			fmt.Println(stream.ID())
			fmt.Println("data length: ", dataLength)
			fmt.Println("data:", string(contentData))
		}
	})

	host9999 := host9999(ctx, host8888.ID(), host8888.Addrs())

	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)
	go writeData(ctx, host9999, host8888)

	select {}
}

func writeData(ctx context.Context, sourceHost, targetHost host.Host) {
	stream, err := sourceHost.NewStream(ctx, targetHost.ID(), ProtocolChat)
	util.PanicIfNotNil(err)
	for {
		message := util.GenerateData(rand.Int() % 1000)

		dataLength := make([]byte, 8)
		binary.LittleEndian.PutUint64(dataLength, uint64(len(message)))
		_, err := stream.Write(dataLength)
		util.PanicIfNotNil(err)

		writeDataLength := 0
		for {
			if writeDataLength == len(message) {
				break
			}

			length, err := stream.Write(message[writeDataLength:])
			util.PanicIfNotNil(err)

			writeDataLength += length
		}
	}
}

func host8888() host.Host {
	port := 8888

	r := rand.New(rand.NewSource(int64(port)))
	fmt.Println(r.Int())

	privateKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.ECDSA, 2048, r)
	util.PanicIfNotNil(err)
	privateKeyBytes, err := privateKey.Raw()
	util.PanicIfNotNil(err)
	publicKeyBytes, err := pubKey.Raw()
	util.PanicIfNotNil(err)

	privateKeyString := hex.EncodeToString(privateKeyBytes)
	publicKeyString := hex.EncodeToString(publicKeyBytes)

	fmt.Printf("private key:%s\npublic key:%s\n", privateKeyString, publicKeyString)

	multiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	util.PanicIfNotNil(err)
	fmt.Println(multiAddr.String())

	host, err := libp2p.New(libp2p.ListenAddrs(multiAddr), libp2p.Identity(privateKey))
	util.PanicIfNotNil(err)

	fmt.Println(host.ID())

	return host
}

func host9999(ctx context.Context, peerId peer.ID, multiAddress []multiaddr.Multiaddr) host.Host {
	port := 9999

	r := rand.New(rand.NewSource(int64(port)))
	fmt.Println(r.Int())

	privateKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.ECDSA, 2048, r)
	util.PanicIfNotNil(err)
	privateKeyBytes, err := privateKey.Raw()
	util.PanicIfNotNil(err)
	publicKeyBytes, err := pubKey.Raw()
	util.PanicIfNotNil(err)

	privateKeyString := hex.EncodeToString(privateKeyBytes)
	publicKeyString := hex.EncodeToString(publicKeyBytes)

	fmt.Printf("private key:%s\npublic key:%s\n", privateKeyString, publicKeyString)

	multiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	util.PanicIfNotNil(err)
	fmt.Println(multiAddr.String())

	host, err := libp2p.New(libp2p.ListenAddrs(multiAddr), libp2p.Identity(privateKey))
	util.PanicIfNotNil(err)

	fmt.Println(host.ID())

	err = host.Connect(ctx, peer.AddrInfo{
		ID:    peerId,
		Addrs: multiAddress,
	})

	util.PanicIfNotNil(err)
	return host
}
