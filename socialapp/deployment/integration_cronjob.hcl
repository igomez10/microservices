job "docs" {
    type = "batch"
    periodic {
        cron             = "*/1 * * * * *"
        prohibit_overlap = false
    }
    datacenters = ["dc1"]
    group "docs" {
        task "docs" {
        driver = "docker"
    
        config {
            image = "igomeza/socialapptests"
        }
        resources {
            cpu    = 500
            memory = 1024
        }
        }
    }
}
