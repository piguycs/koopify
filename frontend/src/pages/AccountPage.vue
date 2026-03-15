<script setup lang="ts">
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import AccountCard from "@/components/AccountCard.vue"
import EditableField from "@/components/EditableField.vue"
import ActionCard from "@/components/ActionCard.vue"
import ModalDialog from "@/components/ModalDialog.vue"
import Button from "@/components/Button.vue"
import { ApiError } from "@/api/client"

const authStore = useAuthStore()
const router = useRouter()

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

async function saveField(payload: { displayName?: string; email?: string }, success: string) {
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
    value !== currentUser.value?.email && saveField({ email: value }, "Email updated.")

const deletionScheduledAt = computed(() => {
    const value = authStore.currentUser?.deletionScheduledAt
    return value ? new Date(value) : null
})

const deletionScheduledLabel = computed(() =>
    deletionScheduledAt.value
        ? deletionScheduledAt.value.toLocaleString(undefined, {
              hour12: false,
              timeZoneName: "short",
          })
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
        statusMessage.value =
            "Deletion requested. You can cancel until the scheduled deletion time."
        showDeletionModal.value = false
    } catch (err) {
        apiError(err, "Account deletion update failed")
    } finally {
        isDeleting.value = false
    }
}

const isAdmin = computed(() => Boolean(currentUser.value?.admin))

const goToAdminInventory = () => {
    router.push("/admin/inventory")
}

const goToAdminUsers = () => {
    router.push("/admin/users")
}
</script>

<template>
    <AppLayout>
        <section class="account">
            <div class="account-stack">
                <AccountCard
                    v-if="isAdmin"
                    eyebrow="Admin"
                    title="Admin actions"
                    description="Manage inventory and user access."
                    badge="ADMIN"
                >
                    <ActionCard
                        title="Manage Inventory"
                        description="Add, edit, and update stock."
                        button-text="Open inventory"
                        @click="goToAdminInventory"
                    />
                    <ActionCard
                        title="Manage Users"
                        description="Promote users and handle resets."
                        button-text="Open users"
                        @click="goToAdminUsers"
                    />
                </AccountCard>
                <AccountCard
                    eyebrow="Account"
                    title="View account"
                    description="Manage your profile details."
                >
                    <!-- TODO: move this to a toast component -->
                    <p v-if="errorMessage" class="toast error">{{ errorMessage }}</p>
                    <p v-if="statusMessage" class="toast success">{{ statusMessage }}</p>
                    <EditableField
                        label="Display name"
                        :value="userDisplayName"
                        @save="saveDisplayName"
                    />
                    <EditableField label="Email" :value="userEmail" @save="saveEmail" />
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
                </AccountCard>
            </div>
        </section>
        <ModalDialog
            :open="showDeletionModal"
            title="Confirm account deletion"
            description="Requesting deletion starts a 24 hour window. You can sign in during that time to cancel. After the scheduled time, logins are blocked."
            @close="closeDeletionModal"
        >
            <template #actions>
                <Button variant="ghost" type="button" @click="closeDeletionModal">Cancel</Button>
                <Button
                    variant="danger"
                    type="button"
                    @click="confirmDeletionRequest"
                    :disabled="isDeleting"
                >
                    Request deletion
                </Button>
            </template>
        </ModalDialog>
    </AppLayout>
</template>

<style scoped>
.account {
    display: flex;
    justify-content: center;
}

.account-stack {
    width: min(520px, 100%);
    display: grid;
    gap: 18px;
}

.toast {
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
