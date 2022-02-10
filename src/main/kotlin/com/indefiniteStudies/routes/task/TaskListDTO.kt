package com.indefiniteStudies.routes.task

import kotlinx.serialization.*

@Serializable
data class TaskListDTO(val limit: Int, val offset: Long, val count: Int, val data: List<TaskDTO>) {
}