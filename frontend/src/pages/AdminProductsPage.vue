<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useProductStore } from "@/stores/products"
import type { ProductResponse } from "@/stores/products"
import AppLayout from "@/layouts/AppLayout.vue"
import ModalDialog from "@/components/ModalDialog.vue"
import Button from "@/components/Button.vue"
import Input from "@/components/Input.vue"
import Toast from "@/components/Toast.vue"
import { ApiError } from "@/api/client"

type SortableKey = "id" | "name" | "priceEurCents" | "inventoryCount" | "isActive"

const router = useRouter()
const productStore = useProductStore()

const products = ref<ProductResponse[]>([])
const isLoading = ref(false)
const errorMessage = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const searchQuery = ref("")
const sortKey = ref<SortableKey | null>(null)
const sortDirection = ref<"asc" | "desc">("asc")

const productToDelete = ref<ProductResponse | null>(null)
const isDeleteModalOpen = computed(() => productToDelete.value !== null)

function formatPrice(cents: number): string {
    return `€${(cents / 100).toFixed(2)}`
}

function errorMessageToast(err: unknown, message: string) {
    if (err instanceof ApiError) {
        errorMessage.value = err.message
    } else {
        errorMessage.value = message
    }
}

function successMessageToast(message: string) {
    successMessage.value = message
}

async function loadProducts() {
    isLoading.value = true
    
    try {
        products.value = await productStore.adminListAllProducts()
    } catch (err) {
        errorMessageToast(err, "Failed to load products")
    } finally {
        isLoading.value = false
    }
}

function navigateToNew() {
    router.push({ name: "admin-products-new" })
}

function navigateToEdit(product: ProductResponse) {
    router.push({ name: "admin-products-edit", params: { id: product.id } })
}

function openDeleteModal(product: ProductResponse) {
    productToDelete.value = product
}

function closeDeleteModal() {
    productToDelete.value = null
}

async function confirmDelete() {
    if (!productToDelete.value) return
    const product = productToDelete.value
    closeDeleteModal()

    try {
        await productStore.adminDeleteProduct(product.id)
        products.value = products.value.filter(p => p.id !== product.id)
        successMessageToast(`Product "${product.name}" deleted`)
    } catch (err) {
        errorMessageToast(err, "Failed to delete product")
    }
}

function clearSearch() {
    searchQuery.value = ""
}

function toggleSort(key: SortableKey) {
    if (sortKey.value === key) {
        if (sortDirection.value === "asc") {
            sortDirection.value = "desc"
        } else {
            sortKey.value = null
        }
    } else {
        sortKey.value = key
        sortDirection.value = "asc"
    }
}

function getSortIndicator(key: SortableKey): string {
    if (sortKey.value !== key) return ""
    return sortDirection.value === "asc" ? " ▲" : " ▼"
}

const filteredProducts = computed(() => {
    if (!searchQuery.value.trim()) return products.value
    const query = searchQuery.value.toLowerCase().trim()
    return products.value.filter(
        p =>
            p.name.toLowerCase().includes(query) ||
            p.slug.toLowerCase().includes(query),
    )
})

const sortedProducts = computed(() => {
    if (!sortKey.value) return filteredProducts.value
    const key = sortKey.value
    const dir = sortDirection.value

    return [...filteredProducts.value].sort((a, b) => {
        let valA: string | number | boolean = a[key]
        let valB: string | number | boolean = b[key]

        if (typeof valA === "string" && typeof valB === "string") {
            valA = valA.toLowerCase()
            valB = valB.toLowerCase()
        }

        if (valA < valB) return dir === "asc" ? -1 : 1
        if (valA > valB) return dir === "asc" ? 1 : -1
        return 0
    })
})

onMounted(() => {
    const toast = (history.state as Record<string, unknown>)?.toast
    if (typeof toast === "string") {
        successMessageToast(toast)
    }
    loadProducts()
})
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <p class="eyebrow">Admin</p>
            <h1>Manage Products</h1>
            <p class="subtle">View, create and manage products.</p>
        </section>

        <div class="admin-content">
            <Toast v-if="successMessage" :message="successMessage" variant="success" @close="successMessage = null" />
            <Toast v-if="errorMessage" :message="errorMessage" variant="error" @close="errorMessage = null" />

            <div v-if="isLoading" class="loading">Loading products</div>

            <div v-else class="table-wrapper">
                <div class="table-controls">
                    <div class="search-box">
                        <Input
                            v-model="searchQuery"
                            type="search"
                            placeholder="Search by name or slug"
                            min-width="240px"
                        />
                        <Button v-if="searchQuery" variant="ghost" @click="clearSearch">
                            Clear
                        </Button>
                    </div>
                    <div class="controls-right">
                        <span class="results-count">
                            {{ sortedProducts.length }}
                            product{{ sortedProducts.length === 1 ? "" : "s" }}
                            <span v-if="searchQuery">found</span>
                        </span>
                        <Button variant="primary" @click="navigateToNew">New Product</Button>
                    </div>
                </div>

                <div class="table-container">
                    <table class="products-table">
                        <thead>
                            <tr>
                                <th class="col-id sortable" @click="toggleSort('id')">
                                    ID{{ getSortIndicator("id") }}
                                </th>
                                <th class="col-name sortable" @click="toggleSort('name')">
                                    Name{{ getSortIndicator("name") }}
                                </th>
                                <th class="col-categories">Categories</th>
                                <th
                                    class="col-price sortable"
                                    @click="toggleSort('priceEurCents')"
                                >
                                    Price{{ getSortIndicator("priceEurCents") }}
                                </th>
                                <th
                                    class="col-stock sortable"
                                    @click="toggleSort('inventoryCount')"
                                >
                                    Stock{{ getSortIndicator("inventoryCount") }}
                                </th>
                                <th class="col-active sortable" @click="toggleSort('isActive')">
                                    Active{{ getSortIndicator("isActive") }}
                                </th>
                                <th class="col-actions">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="product in sortedProducts" :key="product.id">
                                <td class="col-id">{{ product.id }}</td>
                                <td class="col-name">
                                    <div class="truncated" :title="product.name">
                                        {{ product.name }}
                                    </div>
                                    <div class="slug">{{ product.slug }}</div>
                                </td>
                                <td class="col-categories">
                                    <span
                                        v-for="cat in product.categories"
                                        :key="cat.id"
                                        class="category-chip"
                                    >
                                        {{ cat.name }}
                                    </span>
                                    <span v-if="product.categories.length === 0" class="dash">
                                        —
                                    </span>
                                </td>
                                <td class="col-price">
                                    {{ formatPrice(product.priceEurCents) }}
                                </td>
                                <td class="col-stock">{{ product.inventoryCount }}</td>
                                <td class="col-active">
                                    <span v-if="product.isActive" class="badge active">Yes</span>
                                    <span v-else class="badge inactive">No</span>
                                </td>
                                <td class="col-actions">
                                    <div class="action-buttons">
                                        <Button
                                            variant="ghost"
                                            size="small"
                                            @click="navigateToEdit(product)"
                                        >
                                            Edit
                                        </Button>
                                        <Button
                                            variant="danger"
                                            size="small"
                                            @click="openDeleteModal(product)"
                                        >
                                            Delete
                                        </Button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <ModalDialog
            :open="isDeleteModalOpen"
            title="Delete Product"
            :description="`Are you sure you want to delete &quot;${productToDelete?.name}&quot;? This cannot be undone.`"
            @close="closeDeleteModal"
        >
            <template #actions>
                <Button variant="ghost" size="small" @click="closeDeleteModal">Cancel</Button>
                <Button variant="danger" size="small" @click="confirmDelete">Delete</Button>
            </template>
        </ModalDialog>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0;
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

.admin-content {
    margin-top: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.table-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.table-controls {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}

.search-box {
    display: flex;
    align-items: center;
    gap: 8px;
}

.controls-right {
    display: flex;
    align-items: center;
    gap: 16px;
}

.results-count {
    color: var(--muted);
    font-size: 14px;
}

.loading {
    padding: 40px;
    text-align: center;
    color: var(--muted);
}

.table-container {
    background: var(--panel);
    border: 1px solid var(--border);
    overflow-x: auto;
}

.products-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
}

.products-table th,
.products-table td {
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid var(--border);
}

.products-table th {
    background: var(--panel-dark);
    font-weight: 600;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--muted);
    white-space: nowrap;
}

.products-table th.sortable {
    cursor: pointer;
    user-select: none;
}

.products-table th.sortable:hover {
    color: var(--text);
    background: rgba(245, 140, 70, 0.1);
}

.products-table tbody tr:hover {
    background: rgba(245, 140, 70, 0.05);
}

.products-table tbody tr:last-child td {
    border-bottom: none;
}

.col-id {
    width: 60px;
}

.col-name {
    min-width: 160px;
    max-width: 220px;
}

.col-categories {
    min-width: 120px;
}

.col-price {
    width: 90px;
    white-space: nowrap;
}

.col-stock {
    width: 70px;
    text-align: center;
}

.col-active {
    width: 70px;
    text-align: center;
}

.col-actions {
    width: 140px;
}

.truncated {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    cursor: help;
}

.slug {
    font-size: 11px;
    color: var(--muted);
    margin-top: 2px;
}

.category-chip {
    display: inline-block;
    border: 1px solid var(--border);
    padding: 1px 6px;
    font-size: 11px;
    margin-right: 4px;
    white-space: nowrap;
}

.badge {
    padding: 2px 8px;
    font-size: 10px;
    letter-spacing: 2px;
}

.badge.active {
    border: 1px solid rgba(139, 243, 139, 0.4);
    color: #8bf38b;
}

.badge.inactive {
    border: 1px solid var(--border);
    color: var(--muted);
}

.dash {
    color: var(--muted);
}

.action-buttons {
    display: flex;
    gap: 8px;
}

@media (max-width: 768px) {
    .col-id,
    .col-categories {
        display: none;
    }
}
</style>
