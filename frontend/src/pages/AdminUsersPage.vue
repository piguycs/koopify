<script setup lang="ts">
import AppLayout from "@/layouts/AppLayout.vue"
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <p class="eyebrow">Admin</p>
            <h1>Manage Users</h1>
            <p class="subtle">Stub page for user administration.</p>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0;
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
</style>
