package com.indefiniteStudies.services.db

import com.indefiniteStudies.services.config.ConfigService
import org.jetbrains.exposed.sql.Database

object DBService {
    val connection by lazy {
        Database.connect(ConfigService.dotenv["DATABASE_URL"], driver = ConfigService.dotenv["DATABASE_DRIVER_NAME"], user = ConfigService.dotenv["DATABASE_USER"], password = ConfigService.dotenv["DATABASE_PASSWORD"])
    }
}