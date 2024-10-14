terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.13"
    }
  }
}

provider "docker" {}

# Create Docker network
resource "docker_network" "clickhouse_network" {
  name = "clickhouse_network"
}

# Pull ClickHouse Docker image
resource "docker_image" "clickhouse_image" {
  name = "clickhouse/clickhouse-server:latest"
}

# Create ClickHouse container
resource "docker_container" "clickhouse" {
  name  = "clickhouse"
  image = docker_image.clickhouse_image.name
  ports {
    internal = 8123
    external = 8123
  }
  env = [
    "CLICKHOUSE_DB=analytics"
  ]
  networks_advanced {
    name = docker_network.clickhouse_network.name
  }
  
}
resource "docker_image" "cubejs" {
  name = "cubejs/cube"
}

resource "docker_container" "cubejs" {
  image = docker_image.cubejs.latest
  name  = "cubejs"
  ports {
    internal = 4000
    external = 4000
  }
  env = [
   "CUBEJS_DB_TYPE=clickhouse",
    "CUBEJS_DB_HOST=docker_container.clickhouse.network_data[0].gateway",
    "CUBEJS_DB_PORT=8123",
    "CUBEJS_DB_NAME=analytics",
    "CUBEJS_DB_USER=default",
    "CUBEJS_DB_PASS=",
    "CUBEJS_API_SECRET=",
    "CUBEJS_DEV_MODE=true"
  ]
  depends_on = [docker_container.clickhouse]
}