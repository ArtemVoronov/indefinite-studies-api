package com.indefiniteStudies.services.config

import io.github.cdimascio.dotenv.Dotenv

object ConfigService {
    val dotenv = Dotenv.configure()
        .directory("./.env")
        .load()
}