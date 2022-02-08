package com.indefiniteStudies.routes.task

import com.indefiniteStudies.model.entities.TaskState

data class TaskDTO(val id: Int?, val name: String, val state: TaskState) {
}