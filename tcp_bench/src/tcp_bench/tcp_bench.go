package main

import (
	"bufio"
	uuid "code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type TokenInfo struct {
	SessionName SessionId
	Token       string
	TimeStamp   time.Time
}

type Timing struct {
	Sent     map[string]time.Time
	Received map[string]time.Time
}

type SessionId string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// listen for TCP connections
	ln, err := net.Listen("tcp", ":11111")
	if err != nil {
		os.Exit(1)
	}

	// chanel for regular pings
	pingChan := make(chan bool)
	replyChan := make(chan TokenInfo)
	sentChan := make(chan TokenInfo)
	dumpChan := make(chan SessionId)

	// handle reply's
	go handleReply(replyChan, sentChan, dumpChan)

	// hand out pings
	go ping(pingChan, 500*time.Millisecond)
	for {
		conn, err := ln.Accept()
		if err != nil {
			os.Exit(1)
		}
		fmt.Println("New connection")
		go handleConnection(conn, pingChan, replyChan, sentChan, dumpChan)
	}
}

func handleConnection(conn net.Conn, pingChannel chan bool, replyChan chan TokenInfo, sentChan chan TokenInfo, dumpChan chan SessionId) {
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return
	}

	// if this is the start of a communication
	var sessionName SessionId
	msg = strings.TrimSpace(msg)
	if msg == "start" {
		sessionName = SessionId(time.Now().Format(time.RFC3339))
	} else if strings.HasPrefix(msg, "start") {
		sessionName = SessionId(msg[5:])
	} else {
		conn.Close()
		return
	}
	sessionName = SessionId(strings.TrimSpace(string(sessionName)))
	conn.Write([]byte(fmt.Sprintf("%v\n", sessionName)))
	defer fmt.Printf("Dropping session %v\n", sessionName)

	// keep reading from this connection while we can!
	dataChan := make(chan string)
	finishedChan := make(chan bool)
	timeoutChan := make(chan bool)
	go keepReading(conn, reader, finishedChan, dataChan, timeoutChan)
	for {
		select {
		case <-pingChannel:

			// send random token
			go func() {
				token := uuid.NewRandom().String()
				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
				_, err := conn.Write([]byte(token + "\n"))
				if err != nil {
					timeoutChan <- true
				}
				fmt.Printf("S: %v\n", token)
				sentChan <- TokenInfo{SessionName: sessionName, Token: token, TimeStamp: time.Now()}
			}()
		case msg = <-dataChan:
			go func() {
				split := strings.Split(msg, "@")
				name := SessionId(split[0])
				token := split[1]
				fmt.Printf("R: %v\n", token)
				replyChan <- TokenInfo{SessionName: name, Token: token, TimeStamp: time.Now()}
			}()
		case <-finishedChan:
			dumpChan <- sessionName
			conn.Close()
			return
		case <-timeoutChan:
			conn.Close()
			return
		}
	}
}

func keepReading(conn net.Conn, reader *bufio.Reader, finished chan bool, data chan string, timeout chan bool) {
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		msg, err := reader.ReadString('\n')
		if err != nil {
			timeout <- true
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "close" {
			finished <- true
			return
		}
		data <- strings.TrimSpace(msg)
	}
}

func handleReply(replyChan chan TokenInfo, sentChan chan TokenInfo, dumpChan chan SessionId) {

	// a map of timings for each session
	var sessions = make(map[SessionId]Timing)
	for {
		select {
		case reply := <-replyChan:
			if session, ok := sessions[reply.SessionName]; !ok {
				fmt.Printf("ERROR: an unknown session token has arived: %v\n", reply.SessionName)
			} else {
				if _, ok := session.Sent[reply.Token]; !ok {
					fmt.Printf("ERROR: received garbage token: %v\n", reply.Token)
				}
				session.Received[reply.Token] = reply.TimeStamp
			}
		case sent := <-sentChan:
			if session, ok := sessions[sent.SessionName]; !ok {
				fmt.Printf("Adding session: %v\n", sent.SessionName)
				sessions[sent.SessionName] = Timing{
					Sent:     make(map[string]time.Time),
					Received: make(map[string]time.Time),
				}
				sessions[sent.SessionName].Sent[sent.Token] = sent.TimeStamp
			} else {
				session.Sent[sent.Token] = sent.TimeStamp
			}
		case session := <-dumpChan:

			sessionResults, ok := sessions[session]
			if !ok {
				fmt.Printf("ERROR: cannot dump invalid session: %v\n", session)
			}

			// write JSON file
			type result struct {
				Id           string
				Timesent     string
				Timereceived string
				Elapsed      string
			}
			type results struct {
				Results []result
			}
			allResults := make([]result, len(sessionResults.Sent))
			i := 0
			for tok, sentResult := range sessionResults.Sent {
				receiveResult, ok := sessionResults.Received[tok]
				elapsed := ""
				if ok {
					elapsed = strconv.FormatFloat(receiveResult.Sub(sentResult).Seconds()*1000, 'f', 6, 64)
				}
				allResults[i] = result{
					Id:           tok,
					Timesent:     sentResult.Format(time.RFC3339),
					Timereceived: receiveResult.Format(time.RFC3339),
					Elapsed:      elapsed,
				}
				i++
			}
			fmt.Printf("Sent %v, received %v\n", len(sessionResults.Sent), len(sessionResults.Received))
			b, err := json.Marshal(results{Results: allResults})
			if err != nil {
				fmt.Printf("ERROR: error making json for %v\n", session)
			}
			err = ioutil.WriteFile(fmt.Sprintf("./web/%v.json", session), b, 0644)
			if err != nil {
				fmt.Printf("ERROR: error writing json file for %v\n", session)
			}
			fmt.Printf("DUMPED %v\n", session)
		}
	}
}

func ping(pingChan chan bool, interval time.Duration) {
	for {
		time.Sleep(interval)
		pingChan <- true
	}
}
