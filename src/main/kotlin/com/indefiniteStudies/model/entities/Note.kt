package com.indefiniteStudies.model.entities

import org.jetbrains.exposed.sql.Column
import org.jetbrains.exposed.sql.Table
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

object Note : Table("notes") {
    val id: Column<Int> = integer("id").autoIncrement()
    val text: Column<String> = text("text")
    val topic: Column<String> = varchar("topic", 512)
    val tagId: Column<Int> = integer("tag_id")
    val userId: Column<Int> = integer("user_id")
    val state: Column<String> = varchar("state", 256)
    val createDate: Column<LocalDateTime> = datetime("create_date")
    val lastUpdateDate: Column<LocalDateTime> = datetime("last_update_date")
}