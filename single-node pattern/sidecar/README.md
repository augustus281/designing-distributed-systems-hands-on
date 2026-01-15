# Hands On: Deploying the topz Container

The topz sidecar provides a web interface to view resource usage of processes in a container (similar to the `top` command). It works by sharing the PID namespace with an application container.

## Step 1: Build the Application

```bash
docker build -t sidecar-example .
```

## Step 2: Run the Application Container

```bash
APP_ID=$(docker run -d -p 8080:8080 sidecar-example)
```

Verify it's running:

```bash
docker ps
curl http://localhost:8080
```

## Step 3: Build and Run the topz Sidecar

Build the topz sidecar (ARM-compatible):

```bash
docker build -t topz ./topz
```

Run the topz sidecar in the same PID namespace as the application (using port 9090 since app uses 8080):

```bash
APP_ID=$(docker run -d -p 8080:8080 sidecar-example)

docker run --pid=container:${APP_ID} \
  -p 9090:8080 \
  topz
```

## Step 4: View the Results

- Application: http://localhost:8080
- topz sidecar: http://localhost:9090/topz

The topz page shows all processes running in your application container and their resource usage.