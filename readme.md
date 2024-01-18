
```markdown
# Redis Stack Setup for Golang Projects

This guide will walk you through setting up a Redis stack for two Golang projects, including a browser UI for Redis to monitor data. We will use Docker to install and run the Redis stack.

## Prerequisites

Make sure you have Docker installed on your machine.

- [Docker](https://www.docker.com/)

## Step 1: Clone the Repository

```bash
git clone https://github.com/Naveenbhargav-Mallampati/sarvatask.git
cd sarvatask
```

## Step 2: Run Redis Stack via Docker

```bash
docker-compose up -d
```

This will start the Redis server and the associated UI.

## Step 3: Run Golang Projects

Navigate to the root folder containing the two Golang modules.

### `filesarva`

```bash
cd filesarva
./filesarva.exe
```

### `raftinstance`

```bash
cd raftinstance
./raftinstance -state_dir=$TMPDIR/1 -raft :8080 -api :9090
```

In a new shell:

```bash
cd raftinstance
./raftinstance -state_dir $TMPDIR/3 -raft :8082 -api :9092 -join :8080
```

This command adds another Raft node to the cluster.

## Step 4: Upload a File

Use your preferred method to send a POST request to `http://127.0.0.1:8000/upload` with form data. Ensure the key is set as `file` and attach the file to be uploaded.

### Expected Output

You should receive a response indicating the successful upload, along with the JSON returned by the Raft instance.

### please find the demo video in the link:
 - [Demos drive link](https://drive.google.com/drive/folders/1EipDiF6BUFA-fE8t5fzJRsisciWLUp-K?usp=sharing)

## Conclusion

Your Redis stack is now set up, and your Golang projects can interact with it. Explore the Redis UI and monitor the data using the provided commands.

```