import { defineStore } from "pinia"
import { apiClient, ApiError } from "@/api/client"

export type UserResponse = {
    id: number
    displayName: string
    email: string
    admin: boolean
    requestedDeletionAt?: string | null
    deletionScheduledAt?: string | null
}

type LoginResponse = {
    token: string
}

export type UpdateUserPayload = {
    displayName?: string
    email?: string
}

const TOKEN_KEY = "koopify_token"

export const useAuthStore = defineStore("auth", {
    state: () => ({
        token: localStorage.getItem(TOKEN_KEY) ?? "",
        currentUser: null as UserResponse | null,
    }),
    getters: {
        isAuthenticated: (state) => Boolean(state.token),
    },
    actions: {
        setToken(token: string) {
            this.token = token
            localStorage.setItem(TOKEN_KEY, token)
        },
        clearToken() {
            this.token = ""
            localStorage.removeItem(TOKEN_KEY)
        },
        async signIn(email: string, password: string) {
            const resp = await apiClient.post<LoginResponse>("/public_api/v1/login", {
                email,
                password,
            })
            this.setToken(resp.token)
            await this.fetchCurrentUser()
        },
        async signUp(displayName: string, email: string, password: string) {
            await apiClient.post<UserResponse>("/public_api/v1/sign_up", {
                displayName,
                email,
                password,
            })
            await this.signIn(email, password)
        },
        async fetchCurrentUser() {
            if (!this.token) {
                return null
            }

            try {
                const user = await apiClient.get<UserResponse>("/api/v1/users/me", {
                    authToken: this.token,
                })
                this.currentUser = user
                return user
            } catch (err) {
                if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
                    this.signOut()
                }
                throw err
            }
        },
        signOut() {
            this.currentUser = null
            this.clearToken()
        },
        async updateCurrentUser(update: UpdateUserPayload) {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            try {
                const user = await apiClient.patch<UserResponse>("/api/v1/users/me", update, {
                    authToken: this.token,
                })
                this.currentUser = user
                return user
            } catch (err) {
                if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
                    this.signOut()
                }
                throw err
            }
        },
        async requestDeletion() {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            try {
                const user = await apiClient.post<UserResponse>("/api/v1/users/me/deletion", undefined, {
                    authToken: this.token,
                })
                this.currentUser = user
                return user
            } catch (err) {
                if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
                    this.signOut()
                }
                throw err
            }
        },
        async cancelDeletion() {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            try {
                const user = await apiClient.delete<UserResponse>("/api/v1/users/me/deletion", {
                    authToken: this.token,
                })
                this.currentUser = user
                return user
            } catch (err) {
                if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
                    this.signOut()
                }
                throw err
            }
        },

        async listUsers(): Promise<UserResponse[]> {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            return await apiClient.get<UserResponse[]>("/api/v1/users", {
                authToken: this.token,
            })
        },

        async getUserById(id: number): Promise<UserResponse> {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            return await apiClient.get<UserResponse>(`/api/v1/users/${id}`, {
                authToken: this.token,
            })
        },

        async updateUserAdmin(id: number, admin: boolean): Promise<UserResponse> {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            return await apiClient.patch<UserResponse>(`/api/v1/users/${id}/admin`, { admin }, {
                authToken: this.token,
            })
        },

        async triggerPasswordReset(id: number): Promise<{ message: string }> {
            if (!this.token) {
                throw new ApiError("Not authenticated", 401)
            }

            return await apiClient.post<{ message: string }>(`/api/v1/users/${id}/reset_password`, undefined, {
                authToken: this.token,
            })
        },
    },
})
