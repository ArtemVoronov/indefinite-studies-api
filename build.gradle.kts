val ktor_version: String by project
val kotlin_version: String by project
val logback_version: String by project
val exposed_version: String by project
val postgresql_version: String by project
val dotenv_kotlin: String by project
val kotlinx_serialization_version: String by project
val kotlinx_datetime_version: String by project

plugins {
    application
    kotlin("jvm") version "1.6.10"
    kotlin("plugin.serialization") version "1.6.10"
}

group = "com.indefiniteStudies"
version = "0.0.1"
application {
    mainClass.set("com.indefiniteStudies.ApplicationKt")
}

repositories {
    mavenCentral()
    maven { url = uri("https://maven.pkg.jetbrains.space/public/p/ktor/eap") }
}

dependencies {
    implementation("io.ktor:ktor-server-core:$ktor_version")
    implementation("io.ktor:ktor-server-netty:$ktor_version")
    implementation("ch.qos.logback:logback-classic:$logback_version")
    implementation("org.jetbrains.exposed:exposed-core:$exposed_version")
    implementation("org.jetbrains.exposed:exposed-dao:$exposed_version")
    implementation("org.jetbrains.exposed:exposed-jdbc:$exposed_version")
    implementation("org.jetbrains.exposed:exposed-java-time:$exposed_version")
    implementation("org.postgresql:postgresql:$postgresql_version")
    implementation("io.github.cdimascio:dotenv-kotlin:$dotenv_kotlin")
    implementation("io.ktor:ktor-serialization:$ktor_version")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$kotlinx_serialization_version")
    implementation("org.jetbrains.kotlinx:kotlinx-datetime:$kotlinx_datetime_version")


    testImplementation("io.ktor:ktor-server-tests:$ktor_version")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit:$kotlin_version")
}