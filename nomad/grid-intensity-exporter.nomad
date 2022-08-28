
job "grid-intensity-exporter" {

  # The "datacenters" parameter specifies the list of datacenters which should
  # be considered when placing this task. This must be provided.
  datacenters = ["dc1"]

  # the exporter job runs as a service with a single instance that
  # can be scraped by prometheus.
  type = "service"

  group "grid-intensity-exporter" {

    count = 1

    network {
      # for testing, we can get away with having a fixed port
      # but in production we'd let nomad allocate a port instead
      port "exporter" {
        static = 8000
        to = 8000 
      }
    }

    task "grid-intensity-exporter" {
      
      driver = "docker"
      
      config {
        args = [
          "exporter"
        ]
        image = "thegreenwebfoundation/grid-intensity:integration-test"
        ports = ["exporter"]
      }

      env {
        GRID_INTENSITY_LOCATION = "GBR"
        GRID_INTENSITY_PROVIDER = "Ember"
        GRID_INTENSITY_NODE = "${node.unique.name}"
        GRID_INTENSITY_REGION = "${node.region}"
      }
    }
  }
}
