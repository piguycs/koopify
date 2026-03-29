<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useAuthStore, type OrderResponse } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import Toast from "@/components/Toast.vue"
import { ApiError } from "@/api/client"

const authStore = useAuthStore()

const orders = ref<OrderResponse[]>([])
const isLoading = ref(false)
const pollingOrderId = ref<number | null>(null)
const errorMessage = ref("")

const expandedOrderId = ref<number | null>(null)

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

async function loadOrders() {
    isLoading.value = true
    errorMessage.value = ""
    try {
        orders.value = await authStore.listAllOrders()
    } catch (err) {
        // no content, means not really an error
        if (err instanceof ApiError && err.status == 204) {
            // we dont set any message actually, not an error condition
            // errorMessage.value = err.message
        } else if (err instanceof ApiError) {
            errorMessage.value = JSON.stringify(err)
        } else {
            errorMessage.value = "Failed to load orders"
        }
    } finally {
        isLoading.value = false
    }
}

async function refreshOrderStatus(orderId: number) {
    if (pollingOrderId.value === orderId) return
    pollingOrderId.value = orderId
    try {
        const updatedOrder = await authStore.pollOrder(orderId)
        const index = orders.value.findIndex((o) => o.id === orderId)
        if (index !== -1) {
            orders.value[index] = updatedOrder
        }
    } catch (err) {
        if (err instanceof ApiError) {
            errorMessage.value = err.message
        } else {
            errorMessage.value = "Failed to refresh order status"
        }
    } finally {
        pollingOrderId.value = null
    }
}

onMounted(() => {
    loadOrders()
})
</script>

<template>
    <AppLayout>
        <section class="page-header">
            <div class="header-content">
                <h1>Orders</h1>
                <p>View all orders across the platform</p>
            </div>
        </section>

        <section class="orders-section">
            <Toast
                v-if="errorMessage"
                :message="errorMessage"
                variant="error"
                @close="errorMessage = ''"
            />

            <div v-if="isLoading" class="loading">Loading orders</div>

            <div v-else-if="orders.length === 0" class="empty-state">
                <p>No orders found</p>
            </div>

            <div v-else class="table-container">
                <table class="orders-table">
                    <thead>
                        <tr>
                            <th class="col-id">ID</th>
                            <th class="col-user">User ID</th>
                            <th class="col-status">Status</th>
                            <th class="col-total">Total</th>
                            <th class="col-date">Date</th>
                            <th class="col-items">Items</th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-for="order in orders" :key="order.id">
                            <tr
                                :class="{ expanded: expandedOrderId === order.id }"
                                @click="toggleExpand(order.id)"
                            >
                                <td class="col-id">#{{ order.id }}</td>
                                <td class="col-user">{{ order.userId }}</td>
                                <td class="col-status">
                                    <div class="status-cell">
                                        <span :class="['status-badge', order.status]">
                                            {{ formatStatus(order.status) }}
                                        </span>
                                        <button
                                            class="refresh-btn"
                                            :disabled="pollingOrderId === order.id"
                                            :title="'Refresh status'"
                                            @click.stop="refreshOrderStatus(order.id)"
                                        >
                                            <span v-if="pollingOrderId === order.id">...</span>
                                            <span v-else>↻</span>
                                        </button>
                                    </div>
                                </td>
                                <td class="col-total">{{ formatPrice(order.totalEurCents) }}</td>
                                <td class="col-date">{{ formatDate(order.createdAt) }}</td>
                                <td class="col-items">
                                    <span class="items-count"
                                        >{{ order.items.length }} item(s)</span
                                    >
                                    <span class="expand-indicator">
                                        {{ expandedOrderId === order.id ? "−" : "+" }}
                                    </span>
                                </td>
                            </tr>
                            <tr v-if="expandedOrderId === order.id" class="order-items-row">
                                <td colspan="6">
                                    <div class="order-items-detail">
                                        <h4>Order Items</h4>
                                        <table class="items-table">
                                            <thead>
                                                <tr>
                                                    <th>Product</th>
                                                    <th>Quantity</th>
                                                    <th>Unit Price</th>
                                                    <th>Subtotal</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                <tr v-for="item in order.items" :key="item.id">
                                                    <td>{{ item.productName }}</td>
                                                    <td>{{ item.quantity }}</td>
                                                    <td>{{ formatPrice(item.unitPriceCents) }}</td>
                                                    <td>
                                                        {{
                                                            formatPrice(
                                                                item.unitPriceCents * item.quantity,
                                                            )
                                                        }}
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                        <p v-if="order.adyenReference" class="adyen-ref">
                                            Adyen Reference: {{ order.adyenReference }}
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-header {
    margin-bottom: 32px;
}

.header-content h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0 0 8px;
}

.header-content p {
    margin: 0;
    color: var(--muted);
}

.loading,
.empty-state {
    text-align: center;
    padding: 48px 0;
    color: var(--muted);
}

.table-container {
    overflow-x: auto;
}

.orders-table {
    width: 100%;
    border-collapse: collapse;
    background: var(--panel);
    border: 1px solid var(--border);
}

.orders-table th,
.orders-table td {
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid var(--border);
}

.orders-table th {
    background: rgba(245, 140, 70, 0.05);
    font-weight: 600;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--muted);
}

.orders-table tbody tr {
    cursor: pointer;
    transition: background 0.15s ease;
}

.orders-table tbody tr:hover {
    background: rgba(245, 140, 70, 0.05);
}

.orders-table tbody tr.expanded {
    background: rgba(245, 140, 70, 0.1);
}

.col-id {
    width: 80px;
}

.col-user {
    width: 100px;
}

.col-status {
    width: 120px;
}

.col-total {
    width: 100px;
}

.col-date {
    width: 160px;
}

.col-items {
    width: 140px;
}

.status-badge {
    display: inline-block;
    padding: 4px 10px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    border-radius: 2px;
}

.status-cell {
    display: flex;
    align-items: center;
    gap: 8px;
}

.refresh-btn {
    background: none;
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 4px 8px;
    cursor: pointer;
    color: var(--muted);
    font-size: 14px;
    line-height: 1;
    transition: all 0.15s ease;
}

.refresh-btn:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
}

.refresh-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
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

.items-count {
    color: var(--muted);
    font-size: 13px;
}

.expand-indicator {
    margin-left: 8px;
    color: var(--accent);
    font-weight: 700;
}

.order-items-row td {
    background: rgba(245, 140, 70, 0.03);
    padding: 0;
}

.order-items-detail {
    padding: 20px 24px;
    border-top: 1px solid var(--border);
}

.order-items-detail h4 {
    margin: 0 0 12px;
    font-size: 14px;
    font-weight: 600;
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.items-table {
    width: 100%;
    max-width: 600px;
    border-collapse: collapse;
    background: var(--bg);
    border: 1px solid var(--border);
}

.items-table th,
.items-table td {
    padding: 10px 14px;
    text-align: left;
    border-bottom: 1px solid var(--border);
}

.items-table th {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--muted);
    background: rgba(245, 140, 70, 0.03);
}

.items-table td {
    font-size: 14px;
}

.items-table tbody tr:last-child td {
    border-bottom: none;
}

.adyen-ref {
    margin: 16px 0 0;
    font-size: 12px;
    color: var(--muted);
    font-family: monospace;
}

@media (max-width: 768px) {
    .col-date {
        display: none;
    }
}

@media (max-width: 600px) {
    .col-user {
        display: none;
    }
}
</style>
