<script setup lang="ts">
interface Props {
    modelValue?: string
    type?: "text" | "email" | "password" | "number" | "search"
    placeholder?: string
    required?: boolean
    autocomplete?: string
    disabled?: boolean
    minWidth?: string
}

const props = withDefaults(defineProps<Props>(), {
    type: "text",
})

const emit = defineEmits<{
    "update:modelValue": [value: string],
    keydown: [KeyboardEvent]
}>()

function handleInput(event: Event) {
    const target = event.target as HTMLInputElement
    emit("update:modelValue", target.value)
}
</script>

<template>
    <input
        :value="modelValue"
        :type="type"
        :placeholder="placeholder"
        :required="required"
        :autocomplete="autocomplete"
        :disabled="disabled"
        :style="minWidth ? { minWidth } : undefined"
        @input="handleInput"
        @keydown="emit('keydown', $event)"
    />
</template>

<style scoped>
input {
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 10px 12px;
    font-family: inherit;
    font-size: 14px;
    color: var(--text);
    height: 38px;
    box-sizing: border-box;
    outline: none;
}

input::placeholder {
    color: var(--muted);
}

input:focus {
    border-color: rgba(245, 140, 70, 0.6);
}

input:disabled {
    opacity: 0.6;
    cursor: default;
}
</style>
