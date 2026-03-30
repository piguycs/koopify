<script setup lang="ts">
interface Props {
    variant?: "primary" | "ghost" | "danger" | "link"
    size?: "default" | "small" | "tiny"
    disabled?: boolean
    loading?: boolean
    type?: "button" | "submit" | "reset"
}

const props = withDefaults(defineProps<Props>(), {
    variant: "primary",
    size: "default",
    type: "button",
})

const emit = defineEmits<{
    click: [MouseEvent]
}>()

function handleClick(event: MouseEvent) {
    if (props.disabled || props.loading) return
    emit("click", event)
}
</script>

<template>
    <button
        :class="[variant, size, { disabled: disabled || loading }]"
        :type="type"
        :disabled="disabled || loading"
        @click="handleClick"
    >
        <span v-if="loading" class="spinner"></span>
        <slot />
    </button>
</template>

<style scoped>
button {
    border: 1px solid transparent;
    border-radius: 0;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
    transition:
        transform 0.2s ease,
        box-shadow 0.2s ease;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

button.default {
    padding: 10px 16px;
    font-size: 13px;
    min-height: 38px;
}

button.small {
    padding: 6px 10px;
    font-size: 12px;
    min-height: 30px;
}

button.tiny {
    display: inline-block;
    padding: 3px 8px;
    font-size: 10px;
    /* border-radius: 2px; */
}

button.primary {
    background: #2a1c16;
    color: var(--text);
    border-color: rgba(245, 140, 70, 0.6);
}

button.ghost {
    background: transparent;
    border-color: var(--border-strong);
    color: var(--text);
}

button.danger {
    background: transparent;
    border-color: rgba(243, 139, 139, 0.6);
    color: #f38b8b;
}

button.link {
    background: transparent;
    border: none;
    color: var(--muted);
    padding: 0;
    min-height: auto;
}

button:hover:not(.disabled):not(.link) {
    transform: translateY(-1px);
}

button.primary:hover:not(.disabled) {
    background: rgba(245, 140, 70, 0.1);
}

button.ghost:hover:not(.disabled) {
    background: rgba(245, 140, 70, 0.1);
}

button.danger:hover:not(.disabled) {
    background: rgba(243, 139, 139, 0.1);
}

button.link:hover:not(.disabled) {
    color: var(--text);
    text-decoration: underline;
}

button.disabled {
    opacity: 0.6;
    cursor: default;
}

.spinner {
    width: 14px;
    height: 14px;
    border: 2px solid transparent;
    border-top-color: currentColor;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
</style>
