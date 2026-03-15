<script setup lang="ts">
import { computed, useId } from "vue"

interface Props {
    modelValue?: string
    type?: "text" | "email" | "password" | "number" | "search"
    placeholder?: string
    required?: boolean
    autocomplete?: string
    disabled?: boolean
    variant?: "default" | "dark"
    minWidth?: string
    label?: string
    error?: string
}

const props = withDefaults(defineProps<Props>(), {
    type: "text",
    variant: "default",
})

const emit = defineEmits<{
    "update:modelValue": [value: string],
    keydown: [KeyboardEvent]
}>()

const inputId = useId()

const ariaDescribedBy = computed(() => {
    if (props.error) return `${inputId}-error`
    return undefined
})

const isInvalid = computed(() => !!props.error)

function handleInput(event: Event) {
    const target = event.target as HTMLInputElement
    emit("update:modelValue", target.value)
}
</script>

<template>
    <div class="input-wrapper">
        <label v-if="label" :for="inputId" class="input-label">
            {{ label }}
            <span v-if="required" class="required" aria-hidden="true">*</span>
        </label>
        <input
            :id="inputId"
            :class="[variant]"
            :value="modelValue"
            :type="type"
            :placeholder="placeholder"
            :required="required"
            :autocomplete="autocomplete"
            :disabled="disabled"
            :aria-required="required || undefined"
            :aria-invalid="isInvalid || undefined"
            :aria-describedby="ariaDescribedBy"
            :style="minWidth ? { minWidth } : undefined"
            @input="handleInput"
            @keydown="emit('keydown', $event)"
        />
        <span v-if="error" :id="`${inputId}-error`" class="input-error" role="alert">
            {{ error }}
        </span>
    </div>
</template>

<style scoped>
.input-wrapper {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.input-label {
    font-size: 13px;
    color: var(--muted);
    display: flex;
    align-items: center;
    gap: 4px;
}

.input-label .required {
    color: #f38b8b;
}

input {
    border: 1px solid var(--border);
    padding: 10px 12px;
    font-family: inherit;
    font-size: 14px;
    color: var(--text);
    height: 38px;
    box-sizing: border-box;
    outline: none;
}

input.default {
    background: var(--panel);
}

input.dark {
    background: var(--panel-dark);
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

.input-error {
    font-size: 13px;
    color: #f38b8b;
}
</style>
