<script setup lang="ts">
import { ref, watch } from "vue"
import Button from "@/components/Button.vue"

interface Props {
    label: string
    value?: string | null
    type?: string
    sensitive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    value: null,
    type: "text",
    sensitive: false,
})

const emit = defineEmits<{
    save: [value: string]
}>()

const HIDDEN_VALUE = "********"

const editing = ref(false)
const localValue = ref<string | null>(null)
const initialValue = ref<string | null>(null)

watch(
    () => props.value,
    val => {
        if (!editing.value) localValue.value = val
    },
    { immediate: true },
)

const displayValue = () => (props.sensitive ? HIDDEN_VALUE : props.value || "—")

const startEdit = () => {
    initialValue.value = displayValue()
    localValue.value = props.sensitive ? null : props.value
    editing.value = true
}

const cancelEdit = () => {
    editing.value = false
    localValue.value = props.value
}

const save = () => {
    const currentValue = localValue.value?.trim() ?? ""
    if (currentValue && currentValue !== initialValue.value) {
        emit("save", currentValue)
    }

    editing.value = false
}
</script>

<template>
    <div class="detail">
        <span class="label">{{ label }}</span>
        <div class="detail-value">
            <span v-if="!editing" class="value">
                {{ displayValue() }}
            </span>
            <input
                v-else
                v-model="localValue"
                :type="type"
                class="field-input"
                @keydown.enter="save"
                @keydown.esc="cancelEdit"
            />
            <div v-if="!editing" class="detail-actions">
                <Button variant="link" size="small" type="button" @click="startEdit"> Edit </Button>
            </div>
            <div v-else class="detail-actions">
                <Button variant="ghost" size="small" type="button" @click="save"> Save </Button>
                <Button variant="ghost" size="small" type="button" @click="cancelEdit">
                    Cancel
                </Button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.detail {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 12px 14px;
    background: var(--panel-dark);
    border: 1px solid var(--border);
}

.detail-value {
    display: flex;
    align-items: center;
    gap: 12px;
}

.label {
    color: var(--muted);
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
}

.value {
    font-weight: 600;
}

.detail-actions {
    display: inline-flex;
    gap: 8px;
    align-items: center;
}

.field-input {
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 8px 10px;
    color: var(--text);
    font-family: inherit;
    width: 100%;
}
</style>
