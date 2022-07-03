package com.indefiniteStudies.model.entities

import org.jetbrains.exposed.sql.Column
import org.jetbrains.exposed.sql.Table

object Tag : Table("tags") {
    val id: Column<Int> = integer("id").autoIncrement()
    val name: Column<String> = varchar("varchar", 256)
    val state: Column<String> = varchar("state", 256)
}