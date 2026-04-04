<script setup lang="ts">
import { computed, ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore, type OrderResponse, type UpdateUserPayload } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import AccountCard from "@/components/AccountCard.vue"
import EditableField from "@/components/EditableField.vue"
import ActionCard from "@/components/ActionCard.vue"
import ModalDialog from "@/components/ModalDialog.vue"
import Button from "@/components/Button.vue"
import Toast from "@/components/Toast.vue"
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

const orders = ref<OrderResponse[]>([])
const isLoadingOrders = ref(false)
const ordersError = ref("")
const expandedOrderId = ref<number | null>(null)

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

async function saveField(payload: UpdateUserPayload, success: string) {
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

const saveDisplayName = (value: string) => saveField({ displayName: value }, "Display name updated.")
const saveEmail = (value: string) => saveField({ email: value }, "Email updated.")
const savePassword = (value: string) => saveField({ password: value }, "Password updated")

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
    if (isDeleting.value) return
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

function goToAdminProducts() {
    router.push({ name: "admin-products" })
}

function goToAdminUsers() {
    router.push({ name: "admin-users" })
}

function goToAdminOrders() {
    router.push({ name: "admin-orders" })
}

function toggleExpand(orderId: number) {
    if (expandedOrderId.value === orderId) {
        expandedOrderId.value = null
    } else {
        expandedOrderId.value = orderId
    }
}

function formatPrice(cents: number): string {
    return `€${(cents / 100).toFixed(2)}`
}

function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    })
}

function formatStatus(status: string): string {
    return status.charAt(0).toUpperCase() + status.slice(1)
}

function canCompleteOrder(order: OrderResponse): boolean {
    return order.status === 'pending' && order.adyenPaymentLink != null;
}

function completeOrder(order: OrderResponse) {
    if (!order.adyenPaymentLink) return;
    window.location.href = order.adyenPaymentLink
}

async function loadOrders() {
    isLoadingOrders.value = true
    ordersError.value = ""
    try {
        orders.value = await authStore.listOrders()
    } catch (err) {
        if (err instanceof ApiError) {
            ordersError.value = err.message
        } else {
            ordersError.value = "Failed to load orders"
        }
    } finally {
        isLoadingOrders.value = false
    }
}

onMounted(loadOrders)

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
                        title="Manage Products"
                        description="Add, edit, and manage product listings."
                        button-text="Open products"
                        @click="goToAdminProducts"
                    />
                    <ActionCard
                        title="Manage Users"
                        description="Promote users and handle resets."
                        button-text="Open users"
                        @click="goToAdminUsers"
                    />
                    <ActionCard
                        title="Manage Orders"
                        description="View all orders across the platform."
                        button-text="Open orders"
                        @click="goToAdminOrders"
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
                    <EditableField label="Password" sensitive @save="savePassword" />
                    <ActionCard
                        title="Delete account"
                        :description="deletionDescription"
                        :button-text="deletionActionLabel"
                        danger
                        :disabled="isDeleting"
                        @click="handleDeletionAction"
                    />
                </AccountCard>

                <AccountCard
                    eyebrow="Orders"
                    title="Your orders"
                    description="View your order history."
                >
                    <Toast
                        v-if="ordersError"
                        :message="ordersError"
                        variant="error"
                        @close="ordersError = ''"
                    />

                    <div v-if="isLoadingOrders" class="loading">Loading orders</div>

                    <div v-else-if="orders.length === 0" class="empty-orders">
                        <p>No orders yet</p>
                    </div>

                    <div v-else class="orders-list">
                        <div v-for="order in orders" :key="order.id" class="order-item">
                            <div class="order-header" @click="toggleExpand(order.id)">
                                <div class="order-info">
                                    <span class="order-id">#{{ order.id }}</span>
                                    <span :class="['status-badge', order.status]">
                                        {{ formatStatus(order.status) }}
                                    </span>
                                    <Button
                                        size="tiny"
                                        variant="danger"
                                        v-if="canCompleteOrder(order)"
                                        @click="completeOrder(order)"
                                    >
                                        COMPLETE THE PAYMENT
                                    </Button>
                                </div>
                                <div class="order-meta">
                                    <span class="order-total">{{
                                        formatPrice(order.totalEurCents)
                                    }}</span>
                                    <span class="expand-icon">
                                        {{ expandedOrderId === order.id ? "−" : "+" }}
                                    </span>
                                </div>
                            </div>
                            <p class="order-date">{{ formatDate(order.createdAt) }}</p>

                            <div v-if="expandedOrderId === order.id" class="order-details">
                                <table class="items-table">
                                    <thead>
                                        <tr>
                                            <th>Product</th>
                                            <th>Qty</th>
                                            <th>Price</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr v-for="item in order.items" :key="item.id">
                                            <td>{{ item.productName }}</td>
                                            <td>{{ item.quantity }}</td>
                                            <td>
                                                {{
                                                    formatPrice(item.unitPriceCents * item.quantity)
                                                }}
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
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

.loading {
    text-align: center;
    padding: 24px 0;
    color: var(--muted);
}

.empty-orders {
    text-align: center;
    padding: 24px 0;
    color: var(--muted);
    font-size: 14px;
}

.orders-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.order-item {
    background: var(--bg);
    border: 1px solid var(--border);
}

.order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 14px;
    cursor: pointer;
    transition: background 0.15s ease;
}

.order-header:hover {
    background: rgba(245, 140, 70, 0.05);
}

.order-info {
    display: flex;
    align-items: center;
    gap: 10px;
}

.order-id {
    font-weight: 600;
    font-size: 14px;
}

.order-meta {
    display: flex;
    align-items: center;
    gap: 12px;
}

.order-total {
    font-weight: 600;
    font-size: 14px;
}

.expand-icon {
    color: var(--accent);
    font-weight: 700;
    font-size: 16px;
}

.order-date {
    padding: 0 14px 12px;
    margin: 0;
    font-size: 12px;
    color: var(--muted);
}

.order-details {
    padding: 12px 14px;
    border-top: 1px solid var(--border);
    background: rgba(245, 140, 70, 0.03);
}

.items-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
}

.items-table th,
.items-table td {
    padding: 8px 10px;
    text-align: left;
    border-bottom: 1px solid var(--border);
}

.items-table th {
    font-weight: 600;
    color: var(--muted);
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.items-table tbody tr:last-child td {
    border-bottom: none;
}

.items-table td:last-child {
    text-align: right;
}

.status-badge {
    display: inline-block;
    padding: 3px 8px;
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    border-radius: 2px;
}

.status-badge.pending {
    background: rgba(245, 140, 70, 0.2);
    color: #f58a46;
    border: 1px solid rgba(245, 140, 70, 0.4);
}

.status-badge.completed {
    background: rgba(126, 200, 126, 0.2);
    color: #7ec87e;
    border: 1px solid rgba(126, 200, 126, 0.4);
}

.status-badge.failed,
.status-badge.cancelled {
    background: rgba(243, 139, 139, 0.2);
    color: #f38b8b;
    border: 1px solid rgba(243, 139, 139, 0.4);
}
</style>
