package com.example.kotlin

/**
 * Data class definitions.
 */

data class Request(
    val id: String,
    val method: String,
    val path: String,
    val body: String = ""
)

data class Response(
    val statusCode: Int,
    val body: String,
    val headers: Map<String, String> = emptyMap()
)

/**
 * Unused data class - dead code.
 */
data class ErrorResponse(
    val code: Int,
    val message: String,
    val details: List<String> = emptyList()
)

fun processRequest(req: Request): Response {
    return Response(200, "OK: ${req.path}")
}

/**
 * Unused function - dead code.
 */
fun buildErrorResponse(code: Int, msg: String): ErrorResponse {
    return ErrorResponse(code, msg)
}
