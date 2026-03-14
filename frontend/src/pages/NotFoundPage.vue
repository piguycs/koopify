<script setup lang="ts">
import { RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <h1>Page not found</h1>
            <p>The page you are looking for does not exist.</p>
            <RouterLink class="ghost" to="/">Back to home</RouterLink>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0;
}

.page-section p {
    margin: 0;
    color: var(--muted);
}

.ghost {
    border: 1px solid var(--border-strong);
    border-radius: 0;
    padding: 10px 16px;
    width: fit-content;
    text-decoration: none;
    color: var(--text);
}
</style>
