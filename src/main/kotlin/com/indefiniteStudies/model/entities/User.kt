package com.indefiniteStudies.model.entities

import org.jetbrains.exposed.sql.Column
import org.jetbrains.exposed.sql.Table
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

object User : Table("users") {
    val id: Column<Int> = integer("id").autoIncrement()
    val login: Column<String> = varchar("varchar", 256)
    val email: Column<String> = varchar("email", 512)
    val password: Column<String> = varchar("password", 128)
    val role: Column<String> = varchar("role", 256)
    val state: Column<String> = varchar("state", 256)
    val createDate: Column<LocalDateTime> = datetime("create_date")
    val lastUpdateDate: Column<LocalDateTime> = datetime("last_update_date")
}