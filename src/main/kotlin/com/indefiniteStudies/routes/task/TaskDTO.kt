package com.indefiniteStudies.routes.task

import com.indefiniteStudies.model.entities.TaskState
import kotlinx.serialization.*

@Serializable
data class TaskDTO(val id: Int? = null, val name: String, val state: TaskState) {
}