<script setup lang="ts">
interface Props {
    message: string
    variant?: "success" | "error"
}

withDefaults(defineProps<Props>(), {
    variant: "success",
})

const emit = defineEmits<{
    close: []
}>()

function handleClose() {
    emit("close")
}
</script>

<template>
    <div :class="['toast', variant]">
        <span class="toast-message">{{ message }}</span>
        <button class="toast-close" aria-label="Dismiss" @click="handleClose">×</button>
    </div>
</template>

<style scoped>
.toast {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 12px 16px;
    font-size: 14px;
}

.toast.success {
    background: rgba(139, 243, 139, 0.1);
    border: 1px solid rgba(139, 243, 139, 0.3);
    color: #8bf38b;
}

.toast.error {
    background: rgba(243, 139, 139, 0.1);
    border: 1px solid rgba(243, 139, 139, 0.3);
    color: #f38b8b;
}

.toast-message {
    flex: 1;
}

.toast-close {
    background: none;
    border: none;
    cursor: pointer;
    color: currentColor;
    font-size: 18px;
    line-height: 1;
    padding: 0;
    opacity: 0.7;
    display: flex;
    align-items: center;
}

.toast-close:hover {
    opacity: 1;
}
</style>
