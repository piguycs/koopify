<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useProductStore } from "@/stores/products"
import type { CategoryResponse } from "@/stores/products"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import Checkbox from "@/components/Checkbox.vue"
import Input from "@/components/Input.vue"
import Toast from "@/components/Toast.vue"
import { ApiError } from "@/api/client"

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()

const isNew = computed(() => route.name === "admin-products-new")
const productId = computed(() => {
    const id = route.params.id
    return typeof id === "string" ? parseInt(id, 10) : null
})

// Form fields
const name = ref("")
const slug = ref("")
const description = ref("")
const imageUrl = ref("")
const priceEurCents = ref(0)
const discountPercent = ref<number | null>(null)
const inventoryCount = ref(0)
const inStock = ref(false)
const isActive = ref(true)
const selectedCategories = ref<CategoryResponse[]>([])

let backUrl: string | null = null

// Category dropdown
const allCategories = ref<CategoryResponse[]>([])
const categoryInput = ref("")
const isCreatingCategory = ref(false)

// Page state
const isLoading = ref(false)
const isSaving = ref(false)
const errorMessage = ref("")
const successMessage = ref("")

// Slug is auto-derived from name unless manually edited
const slugManuallyEdited = ref(false)

function slugify(value: string): string {
    return value
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/^-+|-+$/g, "")
}

watch(name, newName => {
    if (!slugManuallyEdited.value) {
        slug.value = slugify(newName)
    }
})

function handleSlugInput(value: string) {
    slug.value = value
    slugManuallyEdited.value = value !== slugify(name.value)
}

// Category dropdown: options are allCategories not already selected,
// filtered by what the user has typed
const categoryDropdownOptions = computed(() => {
    const query = categoryInput.value.toLowerCase().trim()
    const selectedIds = new Set(selectedCategories.value.map(c => c.id))
    return allCategories.value.filter(
        c => !selectedIds.has(c.id) && c.name.toLowerCase().includes(query),
    )
})

// Whether the typed input matches no existing category (used to offer creation)
const canCreateCategory = computed(() => {
    const query = categoryInput.value.trim()
    if (!query) return false
    return !allCategories.value.some(c => c.name.toLowerCase() === query.toLowerCase())
})

function addCategory(cat: CategoryResponse) {
    if (!selectedCategories.value.find(c => c.id === cat.id)) {
        selectedCategories.value.push(cat)
    }
    categoryInput.value = ""
}

function removeCategory(cat: CategoryResponse) {
    selectedCategories.value = selectedCategories.value.filter(c => c.id !== cat.id)
}

async function createAndAddCategory() {
    const name = categoryInput.value.trim()
    if (!name) return

    isCreatingCategory.value = true
    try {
        const newCat = await productStore.adminCreateCategory({
            name,
            slug: slugify(name),
        })
        allCategories.value.push(newCat)
        addCategory(newCat)
    } catch (err) {
        errorMessageToast(err, "Failed to create category")
    } finally {
        isCreatingCategory.value = false
    }
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
    errorMessage.value = ""
}

function validateForm(): string | null {
    if (!name.value.trim()) return "Name is required"
    if (!slug.value.trim()) return "Slug is required"
    if (!description.value.trim()) return "Description is required"
    if (priceEurCents.value < 0) return "Price cannot be negative"
    if (inventoryCount.value < 0) return "Inventory count cannot be negative"
    if (
        discountPercent.value !== null &&
        (discountPercent.value < 0 || discountPercent.value > 100)
    ) {
        return "Discount must be between 0 and 100"
    }
    return null
}

async function save() {
    const validationError = validateForm()
    if (validationError) {
        errorMessage.value = validationError
        return
    }

    isSaving.value = true
    errorMessage.value = ""

    const payload = {
        name: name.value.trim(),
        slug: slug.value.trim(),
        description: description.value.trim(),
        imageUrl: imageUrl.value.trim() || null,
        priceEurCents: priceEurCents.value,
        discountPercent: discountPercent.value,
        inventoryCount: inventoryCount.value,
        inStock: inStock.value,
        isActive: isActive.value,
        categoryIds: selectedCategories.value.map(c => c.id),
    }

    try {
        if (isNew.value) {
            await productStore.adminCreateProduct(payload)
            // router.push({ name: "admin-products", state: { toast: "Product created" } })
            goBackWithToast("Product created")
        } else if (productId.value !== null) {
            await productStore.adminUpdateProduct(productId.value, payload)
            // router.push({ name: "admin-products", state: { toast: "Product saved" } })
            goBackWithToast("Product saved")
        }
    } catch (err) {
        errorMessageToast(err, "Failed to save product")
    } finally {
        isSaving.value = false
    }
}

/// TOAST WILL BE IGNORED IF THE `backUrl` DOES NOT SUPPORT TOAST AND IT IS SET
function goBackWithToast(toast: string) {
    if (backUrl) {
        router.push({ path: backUrl, state: { toast } })
    } else {
        router.push({ name: "admin-products", state: { toast } })
    }
}

function goBack() {
    router.back()
    // router.push({ name: "admin-products" })
}

async function loadProduct() {
    if (isNew.value || productId.value === null) return

    isLoading.value = true
    try {
        const product = await productStore.adminGetProduct(productId.value)
        name.value = product.name
        slug.value = product.slug
        description.value = product.description
        imageUrl.value = product.imageUrl ?? ""
        priceEurCents.value = product.priceEurCents
        discountPercent.value = product.discountPercent ?? null
        inventoryCount.value = product.inventoryCount
        inStock.value = product.inStock
        isActive.value = product.isActive
        selectedCategories.value = [...product.categories]
        slugManuallyEdited.value = true
    } catch (err) {
        errorMessageToast(err, "Failed to load product")
    } finally {
        isLoading.value = false
    }
}

async function loadCategories() {
    try {
        allCategories.value = await productStore.listCategories()
    } catch {
        // Non-fatal: category dropdown just won't have suggestions
    }
}

onMounted(async () => {
    const backUrlData = (history.state as Record<string, unknown>)?.backUrl
    if (typeof backUrlData === "string" && backUrlData) backUrl = backUrlData

    await Promise.all([loadCategories(), loadProduct()])
})
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <p class="eyebrow">Admin / Products</p>
            <h1>{{ isNew ? "New Product" : "Edit Product" }}</h1>
        </section>

        <div v-if="isLoading" class="loading">Loading product...</div>

        <form v-else class="product-form" @submit.prevent="save">
            <Toast
                v-if="successMessage"
                :message="successMessage"
                variant="success"
                @close="successMessage = ''"
            />
            <Toast
                v-if="errorMessage"
                :message="errorMessage"
                variant="error"
                @close="errorMessage = ''"
            />

            <div class="form-section">
                <h2 class="section-title">Basic Info</h2>

                <div class="field-row">
                    <Input v-model="name" label="Name" placeholder="Product name" required />
                    <div class="input-wrapper">
                        <label class="input-label">Slug</label>
                        <input
                            class="default"
                            :value="slug"
                            placeholder="product-slug"
                            @input="handleSlugInput(($event.target as HTMLInputElement).value)"
                        />
                    </div>
                </div>

                <div class="input-wrapper full-width">
                    <label class="input-label">Description</label>
                    <textarea v-model="description" placeholder="Product description" rows="5" />
                </div>

                <div class="field-row">
                    <Input v-model="imageUrl" label="Image URL" placeholder="https://..." />
                </div>
            </div>

            <div class="form-section">
                <h2 class="section-title">Pricing & Inventory</h2>

                <div class="field-row">
                    <div class="input-wrapper">
                        <label class="input-label">Price (cents)</label>
                        <input
                            v-model.number="priceEurCents"
                            class="default"
                            type="number"
                            min="0"
                            placeholder="e.g. 1999 for €19.99"
                        />
                    </div>
                    <div class="input-wrapper">
                        <label class="input-label">Discount (%)</label>
                        <input
                            v-model.number="discountPercent"
                            class="default"
                            type="number"
                            min="0"
                            max="100"
                            placeholder="0–100 or leave empty"
                        />
                    </div>
                    <div class="input-wrapper">
                        <label class="input-label">Inventory Count</label>
                        <input
                            v-model.number="inventoryCount"
                            class="default"
                            type="number"
                            min="0"
                        />
                    </div>
                </div>

                <div class="field-row checkboxes">
                    <Checkbox v-model="inStock" label="In Stock" />
                    <Checkbox v-model="isActive" label="Active (visible on storefront)" />
                </div>
            </div>

            <div class="form-section">
                <h2 class="section-title">Categories</h2>

                <div class="category-chips">
                    <span v-for="cat in selectedCategories" :key="cat.id" class="chip">
                        {{ cat.name }}
                        <button type="button" class="chip-remove" @click="removeCategory(cat)">
                            ×
                        </button>
                    </span>
                    <span v-if="selectedCategories.length === 0" class="no-categories">
                        No categories selected
                    </span>
                </div>

                <div class="category-add-row">
                    <div class="category-dropdown-wrapper">
                        <input
                            v-model="categoryInput"
                            class="default category-input"
                            placeholder="Search or type new category..."
                        />
                        <ul
                            v-if="categoryInput && categoryDropdownOptions.length > 0"
                            class="dropdown"
                        >
                            <li
                                v-for="cat in categoryDropdownOptions"
                                :key="cat.id"
                                class="dropdown-item"
                                @mousedown.prevent="addCategory(cat)"
                            >
                                {{ cat.name }}
                            </li>
                        </ul>
                    </div>
                    <Button
                        v-if="canCreateCategory"
                        variant="ghost"
                        :disabled="isCreatingCategory"
                        @click="createAndAddCategory"
                    >
                        {{ isCreatingCategory ? "Creating..." : `Add "${categoryInput.trim()}"` }}
                    </Button>
                </div>
            </div>

            <div class="form-actions">
                <Button variant="ghost" @click="goBack">Cancel</Button>
                <Button variant="primary" type="submit" :disabled="isSaving">
                    {{ isSaving ? "Saving..." : isNew ? "Create Product" : "Save Changes" }}
                </Button>
            </div>
        </form>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 32px;
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

.loading {
    padding: 40px;
    text-align: center;
    color: var(--muted);
}

.product-form {
    display: flex;
    flex-direction: column;
    gap: 32px;
    max-width: 760px;
}

.form-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.section-title {
    font-family: "Rajdhani", sans-serif;
    font-size: 18px;
    font-weight: 600;
    margin: 0;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border);
}

.field-row {
    display: flex;
    gap: 16px;
    flex-wrap: wrap;
}

.field-row > * {
    flex: 1;
    min-width: 180px;
}

.full-width {
    width: 100%;
}

/* Reuse Input.vue's internal styles for raw inputs/textareas */
.input-wrapper {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.input-label {
    font-size: 13px;
    color: var(--muted);
}

input.default,
textarea {
    border: 1px solid var(--border);
    padding: 10px 12px;
    font-family: inherit;
    font-size: 14px;
    color: var(--text);
    background: var(--panel);
    box-sizing: border-box;
    outline: none;
    width: 100%;
}

input.default {
    height: 38px;
}

textarea {
    resize: vertical;
    min-height: 120px;
}

input.default:focus,
textarea:focus {
    border-color: rgba(245, 140, 70, 0.6);
}

.checkboxes {
    gap: 24px;
}

/* Category chips */
.category-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    min-height: 32px;
    align-items: center;
}

.chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    border: 1px solid rgba(245, 140, 70, 0.4);
    padding: 3px 10px;
    font-size: 13px;
}

.chip-remove {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--muted);
    font-size: 16px;
    padding: 0;
    line-height: 1;
    display: flex;
    align-items: center;
}

.chip-remove:hover {
    color: #f38b8b;
}

.no-categories {
    font-size: 13px;
    color: var(--muted);
}

.category-add-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
}

.category-dropdown-wrapper {
    position: relative;
    flex: 1;
    max-width: 360px;
}

.category-input {
    width: 100%;
}

.dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    background: var(--panel-dark);
    border: 1px solid var(--border);
    list-style: none;
    margin: 0;
    padding: 0;
    z-index: 10;
    max-height: 200px;
    overflow-y: auto;
}

.dropdown-item {
    padding: 8px 12px;
    font-size: 13px;
    cursor: pointer;
}

.dropdown-item:hover {
    background: rgba(245, 140, 70, 0.1);
}

/* Actions */
.form-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid var(--border);
}
</style>
