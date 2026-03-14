<script setup lang="ts">
interface Props {
    eyebrow: string
    title: string
    description: string
    badge?: string
}

defineProps<Props>()
</script>

<template>
    <div class="account-card">
        <div class="account-header">
            <p class="eyebrow">{{ eyebrow }}</p>
            <div class="title-row">
                <h2>{{ title }}</h2>
                <span v-if="badge" class="badge">{{ badge }}</span>
            </div>
            <p class="subtle">{{ description }}</p>
        </div>
        <div class="account-body">
            <slot />
        </div>
    </div>
</template>

<style scoped>
.account-card {
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 28px;
}

.account-header h2 {
    font-family: "Rajdhani", sans-serif;
    margin: 8px 0 6px;
    font-size: 26px;
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

.title-row {
    display: flex;
    align-items: center;
    gap: 12px;
}

.badge {
    border: 1px solid rgba(245, 140, 70, 0.6);
    color: var(--text);
    padding: 2px 8px;
    font-size: 10px;
    letter-spacing: 2px;
}

.account-body {
    margin-top: 22px;
    display: grid;
    gap: 16px;
}
</style>
