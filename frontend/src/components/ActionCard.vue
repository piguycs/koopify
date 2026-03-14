<script setup lang="ts">
interface Props {
    title: string
    description: string
    buttonText: string
    danger?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
    click: []
}>()
</script>

<template>
    <div class="action-card">
        <h2>{{ title }}</h2>
        <p>{{ description }}</p>
        <button 
            class="ghost" 
            :class="{ danger }" 
            type="button"
            @click="emit('click')"
        >
            {{ buttonText }}
        </button>
    </div>
</template>

<style scoped>
.action-card {
    background: var(--panel-dark);
    border: 1px solid var(--border);
    padding: 16px;
}

.action-card h2 {
    margin: 0 0 6px;
    font-size: 18px;
    font-family: "Rajdhani", sans-serif;
}

.action-card p {
    margin: 0 0 14px;
    color: var(--muted);
    font-size: 13px;
}

.action-card .ghost {
    border: 1px solid var(--border-strong);
    background: transparent;
    color: var(--text);
    padding: 10px 14px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
}

.action-card .ghost.danger {
    border-color: rgba(243, 139, 139, 0.6);
    color: #f38b8b;
}
</style>
