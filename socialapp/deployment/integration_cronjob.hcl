job "tests" {
    type = "batch"
    periodic {
        cron             = "*/1 * * * * *"
        prohibit_overlap = false
    }
    datacenters = ["dc1"]
    group "docs" {
        count = 8
        task "docs" {
        driver = "docker"
    
        config {
            image = "igomeza/socialapptests"
        }
        resources {
            cpu    = 100
            memory = 256
        }
        }
    }
}