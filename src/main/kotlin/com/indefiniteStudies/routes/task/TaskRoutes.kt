package com.indefiniteStudies.routes.task

import com.indefiniteStudies.model.entities.Task
import com.indefiniteStudies.model.entities.TaskState
import com.indefiniteStudies.services.db.DBService
import io.ktor.application.*
import io.ktor.http.*
import io.ktor.request.*
import io.ktor.response.*
import io.ktor.routing.*
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import org.jetbrains.exposed.sql.transactions.transaction
import org.jetbrains.exposed.sql.update


fun Application.taskRoutes() {

    routing {
        get("/tasks") {
            try {
                val limit = call.request.queryParameters["limit"]?.toIntOrNull()?:50
                val offset = call.request.queryParameters["offset"]?.toLongOrNull()?:0
                val result = transaction(DBService.connection) {
                    Task.select { Task.state neq "${TaskState.DELETED.name}" }
                        .limit(limit, offset = offset)
                        .toList()
                }
                    .map { TaskDTO(
                        id = it[Task.id],
                        name = it[Task.name],
                        state = TaskState.valueOf(it[Task.state]),
                    )}
                    .toList()


                val response = TaskListDTO (
                    limit = limit,
                    offset = offset?:0,
                    count = result.size,
                    data = result,
                )

                call.respond(status = HttpStatusCode.OK, response)

            } catch (e : NumberFormatException) {
                call.respondText("Wrong value at 'limit' or 'offset' parameter, please use number", status = HttpStatusCode.BadRequest)
            } catch (e : Exception) {
                call.respondText("Internal server error", status = HttpStatusCode.InternalServerError)
            }
        }

        get("/task/{id}") {
            try {
                val id = Integer.parseInt(call.parameters["id"])
                val task = transaction(DBService.connection) {
                    Task.select {
                        (Task.id eq id) and (Task.state neq "${TaskState.DELETED.name}")
                    }.firstOrNull()
                }
                if (task == null) {
                    call.respondText("Task with ID '$id' not found", status = HttpStatusCode.BadRequest)
                } else {
                    call.respond(
                        status = HttpStatusCode.OK,
                        TaskDTO(
                            id = task[Task.id],
                            name = task[Task.name],
                            state = TaskState.valueOf(task[Task.state]),
                        )
                    )
                }
            } catch (e : NumberFormatException) {
                call.respondText("Missed ID", status = HttpStatusCode.BadRequest)
            } catch (e : Exception) {
                call.respondText("Internal server error", status = HttpStatusCode.InternalServerError)
            }
        }

        post("/task") {
            try {
                val task = call.receive<TaskDTO>()
                val createdValue = transaction(DBService.connection) {
                    Task.insert {
                        it[name] = task.name
                        it[state] = task.state.name
                    }
                }.resultedValues!!.first()

                call.respond(
                    status = HttpStatusCode.Created,
                    TaskDTO(
                        id = createdValue[Task.id],
                        name = createdValue[Task.name],
                        state = TaskState.valueOf(createdValue[Task.state]),
                    )
                )
            } catch (e : Exception) {
                call.respondText("Internal server error", status = HttpStatusCode.InternalServerError)
            }
        }

        put("/task/{id}") {
            try {
                val task = call.receive<TaskDTO>()
                val id = Integer.parseInt(call.parameters["id"])

                transaction(DBService.connection) {
                    Task.update({ (Task.id eq id) and (Task.state neq "${TaskState.DELETED.name}") }) {
                        it[Task.name] = task.name
                        it[Task.state] = task.state.name
                    }
                }

                // TODO: send NOT found if it does not exist
                call.respondText("Task updated successfully", status = HttpStatusCode.OK)
            } catch (e : NumberFormatException) {
                call.respondText("Missed ID", status = HttpStatusCode.BadRequest)
            } catch (e : Exception) {
                call.respondText("Internal server error", status = HttpStatusCode.InternalServerError)
            }


        }

        delete("/task/{id}") {
            try {
                val id = Integer.parseInt(call.parameters["id"])

                transaction(DBService.connection) {
                    Task.update({ (Task.id eq id) and (Task.state neq "${TaskState.DELETED.name}") }) {
                        it[state] = TaskState.DELETED.toString()
                    }
                }

                // TODO: send NOT found if it does not exist
                call.respondText("Task deleted successfully", status = HttpStatusCode.OK)
            } catch (e : NumberFormatException) {
                call.respondText("Missed ID", status = HttpStatusCode.BadRequest)
            } catch (e : Exception) {
                call.respondText("Internal server error", status = HttpStatusCode.InternalServerError)
            }
        }
    }
}
