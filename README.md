# Word of Wisdom

## Design and implement “Word of Wisdom” tcp server

- TCP server should be protected from DDOS attacks with the Prof of Work (<https://en.wikipedia.org/wiki/Proof_of_work>), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW Challenge.

## PoW algorithm

The provided code is an excellent choice for a test task for several reasons:

- **Simplicity**: The code is relatively simple and easy to understand, making it suitable for testing a candidate's foundational programming skills and algorithmic understanding.
- **Focus on PoW**: The code focuses on the implementation of the Proof of Work (PoW) algorithm, a fundamental concept in blockchain technology. Testing candidates on PoW demonstrates their understanding of consensus mechanisms and how they contribute to the security of distributed systems.
- **Use of SHA-256**: The code employs the SHA-256 hashing algorithm, which is commonly used in blockchain systems for its cryptographic properties. Evaluating candidates' proficiency with hashing functions is relevant in blockchain development.
- **Concurrency and Context**: The code showcases the use of Go's concurrency features by utilizing goroutines and channels. It also demonstrates handling contexts for timeout scenarios, which is valuable knowledge for scalable and efficient systems.
- **Parameterization**: The code allows the difficulty level to be adjusted, enabling candidates to experiment with different settings and observe how the PoW algorithm's mining process behaves under various conditions.

## How to run

### To run in from source

#### Server

```bash
cd ./cmd/server
go run main.go
```

#### Client

```bash
cd ./cmd/client
go run main.go
```

### Run in Docker

```bash
docker compose up
```
