<script setup lang="ts">
import { ref, watch } from "vue"

interface Props {
  label: string
  value?: string | null
  type?: string
}

const props = withDefaults(defineProps<Props>(), {
    value: null,
    type: "text",
})

const emit = defineEmits<{
    save: [value: string],
}>()

const editing = ref(false)
const localValue = ref<string | null>(null)

watch(
    () => props.value,
    (val) => {
        if (!editing.value) localValue.value = val
    },
    { immediate: true },
)

const startEdit = () => {
    localValue.value = props.value
    editing.value = true
}

const cancelEdit = () => {
    editing.value = false
    localValue.value = props.value
}

const save = () => {
    if (localValue.value) {
        emit("save", localValue.value.trim())
    }

    editing.value = false
}
</script>

<template>
    <div class="detail">
        <span class="label">{{ label }}</span>
        <div class="detail-value">
            <span v-if="!editing" class="value">
                {{ value || "—" }}
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
                <button class="link" type="button" @click="startEdit">
                    Edit
                </button>
            </div>
            <div v-else class="detail-actions">
                <button class="ghost" type="button" @click="save">
                    Save
                </button>
                <button class="ghost" type="button" @click="cancelEdit">
                    Cancel
                </button>
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

.ghost {
    border: 1px solid var(--border-strong);
    background: transparent;
    color: var(--text);
    padding: 8px 12px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
}
.link {
    background: transparent;
    border: none;
    color: var(--muted);
    padding: 0;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
}

.link:hover {
    color: var(--text);
    text-decoration: underline;
}
</style>
