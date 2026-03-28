<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import Input from "@/components/Input.vue"
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
        await router.push("/catalogue")
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
                    <Input
                        v-model="displayName"
                        label="Display name"
                        type="text"
                        autocomplete="name"
                        variant="dark"
                        required
                    />
                    <Input
                        v-model="email"
                        label="Email"
                        type="email"
                        autocomplete="email"
                        variant="dark"
                        required
                    />
                    <Input
                        v-model="password"
                        label="Password"
                        type="password"
                        autocomplete="new-password"
                        variant="dark"
                        required
                    />
                    <p v-if="policyMessage" class="policy">{{ policyMessage }}</p>
                    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
                    <Button
                        type="submit"
                        variant="primary"
                        :loading="isSubmitting"
                        :disabled="isSubmitting"
                    >
                        {{ isSubmitting ? "Creating account" : "Create account" }}
                    </Button>
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

.error :deep(.input-error) {
    margin: 0;
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
