<script setup lang="ts">
import { RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
</script>

<template>
    <AppLayout>
        <section class="success-page">
            <!-- Kind of a hacky way to do this with the style, as I didnt want to bring in any icon libs -->
            <div class="success-icon">✓</div>
            <h1>Payment Successful!</h1>
            <p class="message">
                Thank you for your purchase. Your order has been received and is being processed.
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
