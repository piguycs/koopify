<script setup lang="ts">
import { ref } from "vue"
import { RouterLink, useRoute, useRouter } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import Input from "@/components/Input.vue"
import { ApiError } from "@/api/client"
import { useAuthStore } from "@/stores/auth"

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref("")
const password = ref("")
const errorMessage = ref("")
const isSubmitting = ref(false)

const handleSubmit = async () => {
    if (isSubmitting.value) {
        return
    }

    errorMessage.value = ""
    isSubmitting.value = true

    try {
        await authStore.signIn(email.value, password.value)
        const redirect = route.query.redirect
        const nextPath = typeof redirect === "string" && redirect.length > 0 ? redirect : "/catalogue"
        await router.push(nextPath)
    } catch (err) {
        if (err instanceof ApiError) {
            errorMessage.value = err.message
        } else {
            errorMessage.value = "Sign in failed"
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
                    <h1>Sign in</h1>
                    <p class="subtle">Welcome back, tenno.</p>
                </div>

                <form class="auth-form" @submit.prevent="handleSubmit">
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
                        autocomplete="current-password"
                        variant="dark"
                        required
                    />
                    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
                    <Button
                        type="submit"
                        variant="primary"
                        :loading="isSubmitting"
                        :disabled="isSubmitting"
                    >
                        {{ isSubmitting ? "Signing in" : "Sign in" }}
                    </Button>
                </form>

                <p class="hint">
                    No account?
                    <RouterLink to="/sign-up">Create one</RouterLink>
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
