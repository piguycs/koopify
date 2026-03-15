<script setup lang="ts">
import Button from "@/components/Button.vue"

interface Props {
    title: string
    description: string
    buttonText: string
    danger?: boolean
    disabled?: boolean
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
        <Button
            :variant="danger ? 'danger' : 'ghost'"
            type="button"
            :disabled="disabled"
            @click="emit('click')"
        >
            {{ buttonText }}
        </Button>
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
</style>
