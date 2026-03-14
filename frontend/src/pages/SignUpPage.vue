<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import { ApiError, apiClient } from "@/api/client"
import { useAuthStore } from "@/stores/auth"

const authStore = useAuthStore()
const router = useRouter()

const displayName = ref("")
const email = ref("")
const password = ref("")
const errorMessage = ref("")
const isSubmitting = ref(false)
const policyMessage = ref("")
const minLength = ref<number | null>(null)

type PasswordPolicy = {
    minLength: number
    message: string
}

onMounted(async () => {
    try {
        const policy = await apiClient.get<PasswordPolicy>("/public_api/v1/password_policy")
        minLength.value = policy.minLength
        policyMessage.value = policy.message
    } catch {
        minLength.value = null
        policyMessage.value = ""
    }
})

const handleSubmit = async () => {
    if (isSubmitting.value) {
        return
    }

    errorMessage.value = ""

    if (minLength.value !== null && password.value.length < minLength.value) {
        errorMessage.value = policyMessage.value || "Password is too short"
        return
    }

    isSubmitting.value = true

    try {
        await authStore.signUp(displayName.value, email.value, password.value)
        await router.push("/catalog")
    } catch (err) {
        if (err instanceof ApiError) {
            errorMessage.value = err.message
        } else {
            errorMessage.value = "Sign up failed"
        }
    } finally {
        isSubmitting.value = false
    }
}
</script>

<template>
    <AppLayout>
        <section class="auth">
            <div class="auth-card">
                <div class="auth-header">
                    <p class="eyebrow">Access</p>
                    <h1>Sign up</h1>
                    <p class="subtle">Claim your locker.</p>
                </div>

                <form class="auth-form" @submit.prevent="handleSubmit">
                    <label class="field">
                        <span>Display name</span>
                        <input v-model="displayName" type="text" autocomplete="name" required />
                    </label>
                    <label class="field">
                        <span>Email</span>
                        <input v-model="email" type="email" autocomplete="email" required />
                    </label>
                    <label class="field">
                        <span>Password</span>
                        <input v-model="password" type="password" autocomplete="new-password" required />
                    </label>
                    <p v-if="policyMessage" class="policy">{{ policyMessage }}</p>
                    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
                    <button class="primary" type="submit" :disabled="isSubmitting">
                        {{ isSubmitting ? "Creating account" : "Create account" }}
                    </button>
                </form>

                <p class="hint">
                    Already have an account?
                    <RouterLink to="/sign-in">Sign in</RouterLink>
                </p>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.auth {
    display: flex;
    justify-content: center;
}

.auth-card {
    width: min(420px, 100%);
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 28px;
}

.auth-header h1 {
    font-family: "Rajdhani", sans-serif;
    margin: 8px 0 6px;
    font-size: 30px;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 3px;
    font-size: 11px;
    font-weight: 600;
    opacity: 0.7;
    margin: 0;
}

.subtle {
    color: var(--muted);
    margin: 0;
}

.auth-form {
    display: flex;
    flex-direction: column;
    gap: 14px;
    margin-top: 20px;
}

.field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    font-size: 13px;
    color: var(--muted);
}

.field input {
    background: var(--panel-dark);
    border: 1px solid var(--border);
    padding: 12px 14px;
    color: var(--text);
    font-family: inherit;
}

.error {
    margin: 0;
    font-size: 13px;
    color: #f38b8b;
}

.policy {
    margin: 0;
    font-size: 12px;
    color: var(--muted);
}

.primary {
    border: 1px solid transparent;
    border-radius: 0;
    padding: 12px 16px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s ease;
    background: #2a1c16;
    color: var(--text);
    border-color: rgba(245, 140, 70, 0.6);
}

.primary:disabled {
    opacity: 0.6;
    cursor: default;
}

.primary:hover:not(:disabled) {
    transform: translateY(-1px);
}

.hint {
    margin: 18px 0 0;
    font-size: 13px;
    color: var(--muted);
}

.hint a {
    color: var(--text);
    text-decoration: none;
    margin-left: 4px;
}
</style>
