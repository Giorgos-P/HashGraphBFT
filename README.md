# Hashgraph

   Implementation of HashGraph Byzantine Fault tolerance algorithm in Golang programming language with Zero MQ.

## Contents
  - [About](#about)
  - [Modules](#modules)
  - [Installation](#installation)
  - [Execution](#execution)
  - [References](#references)

## About
Implementation of Hashgraph Byzantine Fault tolerance algorithm in Go programming language with ZeroMQ framework.

<p>Based on: <br>
"The Hashgraph Protocol: Efficient Asynchronous BFT for High-Throughput Distributed Ledgers" <br>
by Leemon Baird and Atul Luykx.<p>

## Modules

[Functions](#functions)  | [Protocols](#protocols)
------------- | -------------
[Send Gossip](#send-gossip)  | [Divide Rounds](#)
[Manage Client Request](#manage-client-request)  | [Decide Fame](#)
[Manage Incoming Gossip](#manage-incoming-gossip)  | [Find Order](#)

### Send Gossip
<p> Send Gossip function implements the gossip protocol.<br>
Chooses randomly another server and sends all the event we believe it does not know.
</p>

### Manage Client Request
<p> It is responsible to receive the client transaction and insert it in the graph.<p>

### Manage Incoming Gossip
 * Receives the gossip messages from the other servers
 * Checks whether or not to insert the events in the graph
 * Calls the protocols

### Protocols
The protocols are implemented based on the paper[1]


## Installation
### Golang
If you have not already installed **Golang** follow the instructions [here](https://golang.org/doc/install).

### Clone Repository
```bash
cd ~/go/src/
git clone https://github.com/Giorgos-P/HashGraphBFT.git
```
### Setup
<p>In bin folder create the following paths:<br>
bin/keys<br>
bin/logs/client<br>
bin/logs/error<br>
bin/logs/out<br>
</p>

## Execution
### Hashgraph Client
You can find the instructions for our BFT client [here](https://github.com/Giorgos-P/HashGraph_Client).

### Manually
To install the program and generate the keys run:
```bash
go install HashgraphBFT
HashgraphBFT generate_keys <N>      // For key generation
```
Open <N> different terminals and in each terminal run:
```bash
HashgraphBFT <ID> <N> <Clients> <Scenario> <Remote>
```

ID : is the server id [0 to (N-1)]<br>
N : is the total number of servers<br>
Clients : is the total number of clients<br>
Scenario : is the scenarion number [0 to 2]<br>
Remote : is [0 or 1] 0=local execution, 1 is remote execution (we have to set the correct ip address of the remote servers in ip.go)

### Script
Adjust the script (HashgraphBFT/scripts/run.sh) and run:
```bash
bash ~/go/src/HashgraphBFT/scripts/run.sh
```

### Kill Processes
When you are done and want to kill the processes run:
```bash
bash ~/go/src/HashgraphBFT/scripts/killHash.sh
```

### Scenarios
* 0 = Normal : no malicious nodes occurs
* 1 = IDLE : malicious node do not send messages
* 2 = Sleep : malicious nodes sleeo for some time and then restarts executing
* 3 = Fork : malicious nodes send messages with arbitrary content


## References
- [[1]The Hashgraph Protocol: Efficient Asynchronous BFT for High-Throughput Distributed Ledgers](https://hedera.com/hh-ieee_coins_paper-200516.pdf).


[Go To TOP](#hashgraph)