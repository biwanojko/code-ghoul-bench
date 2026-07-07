package com.example.kotlin

/**
 * Sealed class hierarchy for result types.
 */
sealed class Result<out T> {
    data class Success<T>(val value: T) : Result<T>()
    data class Failure(val error: String) : Result<Nothing>()
    object Loading : Result<Nothing>()
}

fun <T> handleResult(result: Result<T>): String {
    return when (result) {
        is Result.Success -> "success: ${result.value}"
        is Result.Failure -> "error: ${result.error}"
        is Result.Loading -> "loading..."
    }
}

/**
 * Unused sealed class - dead code.
 */
sealed class Command {
    object Start : Command()
    object Stop : Command()
    data class Send(val message: String) : Command()
}

/**
 * Unused handler - dead code.
 */
fun handleCommand(cmd: Command): Unit = when (cmd) {
    is Command.Start -> println("starting")
    is Command.Stop -> println("stopping")
    is Command.Send -> println("sending: ${cmd.message}")
}
