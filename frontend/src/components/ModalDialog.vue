<script setup lang="ts">
interface Props {
    open: boolean
    title: string
    description?: string
}

defineProps<Props>()

const emit = defineEmits<{
    close: []
}>()

const onBackdropClick = () => {
    emit("close")
}
</script>

<template>
    <div v-if="open" class="modal" role="dialog" aria-modal="true">
        <div class="modal-backdrop" @click="onBackdropClick"></div>
        <div class="modal-card">
            <header class="modal-header">
                <h2>{{ title }}</h2>
            </header>
            <p v-if="description" class="modal-description">{{ description }}</p>
            <div class="modal-actions">
                <slot name="actions"></slot>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal {
    position: fixed;
    inset: 0;
    display: grid;
    place-items: center;
    z-index: 50;
}

.modal-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(9, 9, 9, 0.7);
}

.modal-card {
    position: relative;
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 24px;
    width: min(420px, 90vw);
    box-shadow: 0 24px 60px rgba(0, 0, 0, 0.4);
}

.modal-header h2 {
    margin: 0;
    font-family: "Rajdhani", sans-serif;
    font-size: 24px;
}

.modal-description {
    margin: 10px 0 0;
    color: var(--muted);
    font-size: 14px;
    line-height: 1.5;
}

.modal-actions {
    display: flex;
    gap: 12px;
    margin-top: 18px;
}
</style>
