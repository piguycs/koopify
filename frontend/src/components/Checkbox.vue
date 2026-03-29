<script setup lang="ts">
import { useId } from "vue"

interface Props {
    modelValue: boolean
    label?: string
    disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    disabled: false,
})

const emit = defineEmits<{
    "update:modelValue": [value: boolean]
}>()

const inputId = useId()

function handleChange(event: Event) {
    const target = event.target as HTMLInputElement
    emit("update:modelValue", target.checked)
}
</script>

<template>
    <label :for="inputId" :class="['checkbox', { disabled }]">
        <input
            :id="inputId"
            type="checkbox"
            :checked="modelValue"
            :disabled="disabled"
            @change="handleChange"
        />
        <span class="checkbox-box">
            <span v-if="modelValue" class="checkbox-check">✓</span>
        </span>
        <span v-if="label" class="checkbox-label">{{ label }}</span>
    </label>
</template>

<style scoped>
.checkbox {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
}

.checkbox.disabled {
    opacity: 0.6;
    cursor: default;
}

.checkbox input {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
}

.checkbox-box {
    width: 20px;
    height: 20px;
    background: var(--panel-dark);
    border: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
        background-color 0.15s ease,
        border-color 0.15s ease;
    flex-shrink: 0;
}

.checkbox-check {
    font-size: 13px;
    color: rgba(245, 140, 70, 0.9);
    font-weight: bold;
    line-height: 1;
}

.checkbox input:checked + .checkbox-box {
    background: var(--panel-dark);
    border-color: rgba(245, 140, 70, 0.6);
}

.checkbox input:focus + .checkbox-box {
    outline: 2px solid rgba(245, 140, 70, 0.3);
    outline-offset: 1px;
}

.checkbox.disabled .checkbox-box {
    opacity: 0.5;
}

.checkbox-label {
    font-size: 14px;
    color: var(--text);
}
</style>
