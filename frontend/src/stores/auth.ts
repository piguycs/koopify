import { defineStore } from "pinia"
import { apiClient, ApiError } from "@/api/client"

export type UserResponse = {
    id: number
    displayName: string
    email: string
    admin: boolean
}

type LoginResponse = {
    token: string
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
    },
})
