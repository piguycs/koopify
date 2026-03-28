const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ""

type RequestOptions = {
    method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
    body?: unknown
    headers?: Record<string, string>
    authToken?: string
}

type ErrorPayload = {
    error?: {
        code?: string
        message?: string
    }
}

export class ApiError extends Error {
    status: number
    code?: string

    constructor(message: string, status: number, code?: string) {
        super(message)
        this.name = "ApiError"
        this.status = status
        this.code = code
    }
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
        method: options.method ?? "GET",
        headers: {
            "Content-Type": "application/json",
            ...(options.headers ?? {}),
            ...(options.authToken ? { Authorization: `Bearer ${options.authToken}` } : {}),
        },
        body: options.body ? JSON.stringify(options.body) : undefined,
    })

    if (!response.ok) {
        const contentType = response.headers.get("Content-Type") ?? ""
        const payload = contentType.includes("application/json")
            ? ((await response.json().catch(() => undefined)) as ErrorPayload | undefined)
            : undefined

        const message = payload?.error?.message ?? `Request failed: ${response.status}`
        const code = payload?.error?.code
        throw new ApiError(message, response.status, code)
    }

    if (response.status === 204) {
        return undefined as T
    }

    const contentType = response.headers.get("Content-Type") ?? ""
    if (contentType.includes("application/json")) {
        return (await response.json()) as T
    }

    return undefined as T
}

export const apiClient = {
    /// Make a GET request to the `path`
    get: <T>(path: string, options?: RequestOptions) => request<T>(path, options),

    /// Make a POST request to the `path` with the provided `body`
    post: <T>(path: string, body?: unknown, options?: RequestOptions) =>
        request<T>(path, { ...options, method: "POST", body }),

    /// Make a PATCH request to the `path` with the provided `body`
    patch: <T>(path: string, body?: unknown, options?: RequestOptions) =>
        request<T>(path, { ...options, method: "PATCH", body }),

    /// Make a PUT request to the `path` with the provided `body`
    put: <T>(path: string, body?: unknown, options?: RequestOptions) =>
        request<T>(path, { ...options, method: "PUT", body }),

    /// Make a DELETE request to the `path`
    delete: <T>(path: string, options?: RequestOptions) =>
        request<T>(path, { ...options, method: "DELETE" }),
}
