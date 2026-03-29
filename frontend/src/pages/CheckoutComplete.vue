<script setup lang="ts">
import { ref, onMounted } from "vue"
import { RouterLink, useRoute } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import { useAuthStore } from "@/stores/auth"
import { ApiError } from "@/api/client"

const route = useRoute()
const authStore = useAuthStore()

const isLoading = ref(true)
const errorMessage = ref("")

onMounted(async () => {
    const orderId = route.query.orderId as string
    const sessionId = route.query.sessionId as string
    const sessionResult = route.query.sessionResult as string

    if (!orderId || !sessionId || !sessionResult) {
        if (!orderId) errorMessage.value = "Missing required parameters (Order ID)"
        if (!sessionId) errorMessage.value = "Missing required parameters (Session ID)"
        if (!sessionResult) errorMessage.value = "Missing required parameters (Session Res)"
        isLoading.value = false
        return
    }

    try {
        await authStore.updateOrderAdyenSession(parseInt(orderId, 10), sessionId, sessionResult)
    } catch (err) {
        if (err instanceof ApiError) {
            errorMessage.value = err.message
        } else {
            errorMessage.value = "Failed to process payment confirmation"
        }
    } finally {
        isLoading.value = false
    }
})
</script>

<template>
    <AppLayout>
        <section class="success-page">
            <div v-if="isLoading" class="loading-state">
                <div class="loading-icon">...</div>
                <h1>Processing Payment</h1>
                <p class="message">Please wait while we confirm your payment...</p>
            </div>

            <div v-else-if="errorMessage" class="error-state">
                <div class="error-icon">!</div>
                <h1>Something went wrong</h1>
                <p class="message">{{ errorMessage }}</p>
                <p class="sub-message">
                    If you completed payment, please contact support with your order details.
                </p>
                <div class="actions">
                    <RouterLink to="/account">
                        <Button variant="primary">View Account</Button>
                    </RouterLink>
                    <RouterLink to="/catalogue">
                        <Button variant="ghost">Continue Shopping</Button>
                    </RouterLink>
                </div>
            </div>

            <div v-else class="success-state">
                <div class="success-icon">✓</div>
                <h1>Payment Successful!</h1>
                <p class="message">
                    Thank you for your purchase. Your order has been received and is being
                    processed.
                </p>
                <p class="sub-message">You will receive an email confirmation shortly.</p>
                <div class="actions">
                    <RouterLink to="/catalogue">
                        <Button variant="primary">Continue Shopping</Button>
                    </RouterLink>
                    <RouterLink to="/account">
                        <Button variant="ghost">View Account</Button>
                    </RouterLink>
                </div>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.success-page {
    text-align: center;
    padding: 64px 24px;
    max-width: 500px;
    margin: 0 auto;
}

.success-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(126, 200, 126, 0.15);
    border: 2px solid #7ec87e;
    color: #7ec87e;
    font-size: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
}

.loading-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(245, 140, 70, 0.15);
    border: 2px solid #f58a46;
    color: #f58a46;
    font-size: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
    animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
    0%,
    100% {
        opacity: 1;
    }
    50% {
        opacity: 0.5;
    }
}

.error-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(243, 139, 139, 0.15);
    border: 2px solid #f38b8b;
    color: #f38b8b;
    font-size: 40px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
}

h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0 0 16px;
}

.message {
    font-size: 16px;
    color: var(--text);
    margin: 0 0 8px;
    line-height: 1.6;
}

.sub-message {
    font-size: 14px;
    color: var(--muted);
    margin: 0 0 32px;
}

.actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
    align-items: center;
}

.actions a {
    text-decoration: none;
}
</style>
