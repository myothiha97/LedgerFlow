# Architecture vs Infrastructure

## Core idea

Architecture means how a system is designed.

Infrastructure means the support environment that lets the system run.

```text
architecture = how the system is shaped
infrastructure = what runs and supports the system
```

## Software example

Architecture:

```text
Gin handler -> service -> store -> PostgreSQL
```

This describes the parts of the backend and how responsibilities are split.

Infrastructure:

```text
Docker
PostgreSQL container
ports
volumes
server
CI/CD
cloud provider
logs
secrets
```

This describes the environment and tools used to run the backend.

## Docker vs VM

VMs and containers have different runtime architectures:

```text
VM = full fake computer with its own OS kernel
container = isolated process that shares the host kernel
```

Choosing Docker, VMs, or Kubernetes to run a backend is an infrastructure decision.

```text
VM vs container internals = architecture difference
using VM or container to run services = infrastructure choice
```

## Other engineering examples

IoT architecture:

```text
sensor -> microcontroller -> gateway -> cloud API -> dashboard
```

IoT infrastructure:

```text
Wi-Fi
MQTT broker
cloud server
device provisioning
power supply
monitoring
```

Electrical architecture:

```text
battery -> power regulator -> control board -> motor driver -> motor
```

Electrical infrastructure:

```text
wiring
grounding
power distribution
test bench
enclosure
cooling
safety systems
```

Robotics architecture:

```text
perception -> planning -> control -> actuators
```

Robotics infrastructure:

```text
robot OS setup
calibration tools
networking
charging station
logs
simulation environment
```

## Important note

The boundary can blur.

Networking may be infrastructure in a normal web app, but part of architecture in an IoT system.

Use this rule:

```text
If it describes the main design and responsibility split, call it architecture.
If it describes the support system needed to run it, call it infrastructure.
```
