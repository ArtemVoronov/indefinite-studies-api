package com.indefiniteStudies

import com.indefiniteStudies.plugins.configureRouting
import com.indefiniteStudies.routes.task.taskRoutes
import io.ktor.application.*
import io.ktor.features.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.serialization.*

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        install(ContentNegotiation) {
            json()
        }
        taskRoutes()
        configureRouting()
    }.start(wait = true)
}