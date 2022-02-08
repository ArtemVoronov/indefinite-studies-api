package com.indefiniteStudies.model.entities

import org.jetbrains.exposed.sql.Column
import org.jetbrains.exposed.sql.Table

object Task : Table("tasks") {
    val id: Column<Int> = integer("id").autoIncrement()
    val name: Column<String> = varchar("name", 100)
    val state: Column<String> = varchar("state", 100)
}