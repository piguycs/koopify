<script setup lang="ts">
import AppLayout from "@/layouts/AppLayout.vue"
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <h1>Catalogue</h1>
            <p>Browse packs, singles, and curated drops.</p>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0 0 12px;
}

.page-section p {
    margin: 0;
    color: var(--muted);
}
</style>
