package com.indefiniteStudies

import com.fasterxml.jackson.databind.SerializationFeature
import com.indefiniteStudies.routes.task.taskRoutes
import io.ktor.serialization.jackson.*
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.*

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        install(ContentNegotiation) {
            jackson {
                enable(SerializationFeature.INDENT_OUTPUT)
            }
        }
        taskRoutes()
    }.start(wait = true)
}