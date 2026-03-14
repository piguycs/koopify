<script setup lang="ts">
import { computed, ref } from "vue"
import { useAuthStore } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import EditableField from "@/components/EditableField.vue"
import ActionCard from "@/components/ActionCard.vue"
import ModalDialog from "@/components/ModalDialog.vue"
import { ApiError } from "@/api/client"

const authStore = useAuthStore()

const currentUser = computed(() => authStore.currentUser)
const userDisplayName = computed(() => currentUser.value?.displayName)
const userEmail = computed(() => currentUser.value?.email)

const isSaving = ref(false)
const isDeleting = ref(false)
const errorMessage = ref("")
const statusMessage = ref("")
const showDeletionModal = ref(false)

function closeDeletionModal() {
    showDeletionModal.value = false
}

function clearValues() {
    errorMessage.value = ""
    statusMessage.value = ""
}

function apiError(err: unknown, message: string) {
    if (err instanceof ApiError) {
        errorMessage.value = err.message
    } else {
        errorMessage.value = message
    }
}

async function saveField(payload: { displayName?: string; email?: string} , success: string) {
    if (isSaving.value) {
        return
    }

    clearValues()
    isSaving.value = true

    try {
        await authStore.updateCurrentUser(payload)
        statusMessage.value = success
    } catch (err) {
        apiError(err, "Update failed")
    } finally {
        isSaving.value = false
    }
}

const saveDisplayName = (value: string) =>
    value !== currentUser.value?.displayName &&
    saveField({ displayName: value }, "Display name updated.")

const saveEmail = (value: string) =>
    value !== currentUser.value?.email &&
    saveField({ email: value }, "Email updated.")

const deletionScheduledAt = computed(() => {
    const value = authStore.currentUser?.deletionScheduledAt
    return value ? new Date(value) : null
})

const deletionScheduledLabel = computed(() =>
    deletionScheduledAt.value
        ? deletionScheduledAt.value.toLocaleString(undefined, { timeZoneName: "short" })
        : null,
)

const deletionDescription = computed(() => {
    if (!deletionScheduledAt.value) {
        return "Request account deletion. You can cancel until the scheduled deletion time."
    }

    const scheduled = deletionScheduledLabel.value
    return scheduled ? `Deletion scheduled for ${scheduled}.` : "Deletion scheduled."
})

const deletionActionLabel = computed(() =>
    deletionScheduledAt.value ? "Cancel deletion" : "Request deletion",
)

async function handleDeletionAction() {
    if (isDeleting.value) {
        return
    }

    clearValues()
    isDeleting.value = true

    try {
        if (deletionScheduledAt.value) {
            await authStore.cancelDeletion()
            statusMessage.value = "Deletion request canceled."
        } else {
            showDeletionModal.value = true
        }
    } catch (err) {
        apiError(err, "Account deletion update failed")
    } finally {
        isDeleting.value = false
    }
}

async function confirmDeletionRequest() {
    if (isDeleting.value) {
        return
    }

    clearValues()
    isDeleting.value = true

    try {
        await authStore.requestDeletion()
        statusMessage.value = "Deletion requested. You can cancel until the scheduled deletion time."
        showDeletionModal.value = false
    } catch (err) {
        apiError(err, "Account deletion update failed")
    } finally {
        isDeleting.value = false
    }
}

</script>

<template>
    <AppLayout>
        <section class="account">
            <div class="account-card">
                <div class="account-header">
                    <p class="eyebrow">Account</p>
                    <h1>View account</h1>
                    <p class="subtle">Manage your profile details.</p>
                </div>

                <div class="account-content">
                    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
                    <p v-if="statusMessage" class="status">{{ statusMessage }}</p>
                    <EditableField
                        label="Display name"
                        :value="userDisplayName"
                        @save="saveDisplayName"
                    />
                    <EditableField
                        label="Email"
                        :value="userEmail"
                        @save="saveEmail"
                    />
                    <ActionCard
                        title="Reset password"
                        description="WIP: Password reset flow is not wired yet."
                        button-text="Reset password"
                        @click="console.info('WIP: password reset')"
                    />
                    <ActionCard
                        title="Delete account"
                        :description="deletionDescription"
                        :button-text="deletionActionLabel"
                        danger
                        :disabled="isDeleting"
                        @click="handleDeletionAction"
                    />
                </div>
            </div>
        </section>
        <ModalDialog
            :open="showDeletionModal"
            title="Confirm account deletion"
            description="Requesting deletion starts a 24 hour window. You can sign in during that time to cancel. After the scheduled time, logins are blocked."
            @close="closeDeletionModal"
        >
            <template #actions>
                <button class="ghost" type="button" @click="closeDeletionModal">Cancel</button>
                <button class="danger" type="button" @click="confirmDeletionRequest" :disabled="isDeleting">
                    Request deletion
                </button>
            </template>
        </ModalDialog>
    </AppLayout>
</template>

<style scoped>
.account {
    display: flex;
    justify-content: center;
}

.account-card {
    width: min(520px, 100%);
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 28px;
}

.account-header h1 {
    font-family: "Rajdhani", sans-serif;
    margin: 8px 0 6px;
    font-size: 30px;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 3px;
    font-size: 11px;
    font-weight: 600;
    opacity: 0.7;
    margin: 0;
}

.subtle {
    color: var(--muted);
    margin: 0;
}

.account-content {
    margin-top: 22px;
    display: grid;
    gap: 16px;
}

.error {
    margin: 0;
    font-size: 13px;
    color: #f38b8b;
}

.status {
    margin: 0;
    font-size: 13px;
    color: var(--muted);
}

.ghost {
    border: 1px solid var(--border-strong);
    background: transparent;
    color: var(--text);
    padding: 10px 14px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
}

.danger {
    border: 1px solid rgba(243, 139, 139, 0.6);
    background: transparent;
    color: #f38b8b;
    padding: 10px 14px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
}

.ghost:disabled,
.danger:disabled {
    opacity: 0.6;
    cursor: default;
}
</style>
