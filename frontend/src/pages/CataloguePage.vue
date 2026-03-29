<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import { useProductStore, type CategoryResponse, type ProductResponse } from "@/stores/products"
import { useAppStore } from "@/stores/app"

const productStore = useProductStore()
const appStore = useAppStore()

const DEFAULT_PER_PAGE = 16

const products = ref<ProductResponse[]>([])
const categories = ref<CategoryResponse[]>([])
const totalProducts = ref(0)
const currentStart = ref(0)
const currentEnd = ref(DEFAULT_PER_PAGE)
const isLoading = ref(false)
const errorMessage = ref("")

const selectedCategory = ref<string | null>(null)
const searchQuery = ref("")

const filteredProducts = computed(() => {
    if (!searchQuery.value.trim()) {
        return products.value
    }
    const query = searchQuery.value.toLowerCase()
    return products.value.filter(
        p =>
            p.name.toLowerCase().includes(query) ||
            p.description.toLowerCase().includes(query)
    )
})

const hasNextPage = computed(() => currentEnd.value < totalProducts.value)
const hasPrevPage = computed(() => currentStart.value > 0)

const currentPage = computed(() => Math.floor(currentStart.value / DEFAULT_PER_PAGE) + 1)
const totalPages = computed(() => Math.ceil(totalProducts.value / DEFAULT_PER_PAGE))

async function loadProducts() {
    isLoading.value = true
    errorMessage.value = ""
    try {
        const result = await productStore.listActiveProducts(
            currentStart.value,
            currentEnd.value,
            selectedCategory.value || undefined
        )
        products.value = result.products
        totalProducts.value = result.totalProducts
    } catch (err) {
        errorMessage.value = "Failed to load products"
    } finally {
        isLoading.value = false
    }
}

async function loadCategories() {
    try {
        categories.value = await productStore.listCategories()
    } catch {
        console.error("Failed to load categories")
    }
}

function selectCategory(slug: string | null) {
    selectedCategory.value = slug
    currentStart.value = 0
    currentEnd.value = DEFAULT_PER_PAGE
    loadProducts()
}

function nextPage() {
    if (hasNextPage.value) {
        currentStart.value = currentEnd.value
        currentEnd.value = Math.min(currentEnd.value + DEFAULT_PER_PAGE, totalProducts.value)
        loadProducts()
    }
}

function prevPage() {
    if (hasPrevPage.value) {
        currentEnd.value = currentStart.value
        currentStart.value = Math.max(0, currentStart.value - DEFAULT_PER_PAGE)
        loadProducts()
    }
}

function formatPrice(cents: number): string {
    return `€${(cents / 100).toFixed(2)}`
}

function getDiscountedPrice(priceCents: number, discountPercent: number | null): number {
    if (!discountPercent) return priceCents
    return Math.round(priceCents * (1 - discountPercent / 100))
}

function addToCart(product: ProductResponse) {
    appStore.incrementCart()
}

onMounted(() => {
    loadProducts()
    loadCategories()
})

watch(selectedCategory, () => {
    loadProducts()
})
</script>

<template>
    <AppLayout>
        <section class="page-header">
            <div class="header-content">
                <h1>Catalogue</h1>
                <p>Browse packs, singles, and curated drops.</p>
            </div>
        </section>

        <section class="filters">
            <div class="filters-container">
                <div class="search-box">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Search products..."
                        class="search-input"
                    />
                </div>
                <div class="category-filters">
                    <button
                        :class="['category-btn', { active: selectedCategory === null }]"
                        @click="selectCategory(null)"
                    >
                        All
                    </button>
                    <button
                        v-for="cat in categories"
                        :key="cat.id"
                        :class="['category-btn', { active: selectedCategory === cat.slug }]"
                        @click="selectCategory(cat.slug)"
                    >
                        {{ cat.name }}
                    </button>
                </div>
            </div>
        </section>

        <section v-if="isLoading" class="loading">
            <p>Loading products...</p>
        </section>

        <section v-else-if="errorMessage" class="error">
            <p>{{ errorMessage }}</p>
        </section>

        <section v-else-if="filteredProducts.length === 0" class="empty-state">
            <p>No products found.</p>
        </section>

        <section v-else class="product-grid-section">
            <div class="product-grid">
                <article
                    v-for="product in filteredProducts"
                    :key="product.id"
                    class="product-card"
                >
                    <div v-if="product.discountPercent" class="discount-badge">
                        -{{ product.discountPercent }}%
                    </div>
                    <div class="product-media">
                        <img
                            v-if="product.imageUrl"
                            :src="product.imageUrl"
                            :alt="product.name"
                            class="product-image"
                        />
                        <div v-else class="product-placeholder"></div>
                    </div>
                    <div class="product-info">
                        <div class="product-tags">
                            <span
                                v-for="cat in product.categories"
                                :key="cat.id"
                                class="product-tag"
                            >
                                {{ cat.name }}
                            </span>
                        </div>
                        <h4 class="product-name">{{ product.name }}</h4>
                        <div class="product-pricing">
                            <span
                                v-if="product.discountPercent"
                                class="original-price"
                            >
                                {{ formatPrice(product.priceEurCents) }}
                            </span>
                            <span
                                class="price"
                                :class="{ discounted: product.discountPercent }"
                            >
                                {{
                                    formatPrice(
                                        getDiscountedPrice(
                                            product.priceEurCents,
                                            product.discountPercent ?? null
                                        )
                                    )
                                }}
                            </span>
                        </div>
                        <div class="product-actions">
                            <RouterLink
                                :to="`/product/${product.slug}`"
                                class="view-details-btn"
                            >
                                View Details
                            </RouterLink>
                            <Button
                                variant="primary"
                                size="small"
                                :disabled="!product.inStock"
                                @click="addToCart(product)"
                            >
                                Add
                            </Button>
                        </div>
                    </div>
                </article>
            </div>

            <div class="pagination">
                <Button
                    variant="ghost"
                    :disabled="!hasPrevPage"
                    @click="prevPage"
                >
                    Previous
                </Button>
                <span class="page-indicator">
                    Page {{ currentPage }} of {{ totalPages }}
                </span>
                <Button
                    variant="ghost"
                    :disabled="!hasNextPage"
                    @click="nextPage"
                >
                    Next
                </Button>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 24px;
    flex-wrap: wrap;
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

.search-input {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 0;
    padding: 10px 14px;
    font-family: inherit;
    font-size: 14px;
    color: var(--text);
    width: 240px;
}

.search-input::placeholder {
    color: var(--muted);
}

.search-input:focus {
    outline: none;
    border-color: var(--accent);
}

.filters {
    margin-top: 8px;
}

.filters-container {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
}

.search-box {
    flex-shrink: 0;
}

.category-filters {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

.category-btn {
    background: transparent;
    border: 1px solid var(--border);
    border-radius: 0;
    padding: 6px 14px;
    font-family: inherit;
    font-size: 13px;
    font-weight: 500;
    color: var(--muted);
    cursor: pointer;
    transition: all 0.2s ease;
}

.category-btn:hover {
    border-color: var(--accent);
    color: var(--text);
}

.category-btn.active {
    background: var(--accent);
    border-color: var(--accent);
    color: var(--bg);
}

.loading,
.error,
.empty-state {
    text-align: center;
    padding: 48px 0;
    color: var(--muted);
}

.error {
    color: #f38b8b;
}

.product-grid-section {
    display: flex;
    flex-direction: column;
    gap: 32px;
}

.product-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 16px;
}

.product-card {
    background: var(--panel);
    border: 1px solid var(--border);
    border-radius: 0;
    padding: 16px;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.discount-badge {
    position: absolute;
    top: 14px;
    left: 14px;
    background: var(--accent);
    color: var(--bg);
    padding: 3px 10px;
    border-radius: 0;
    font-size: 10px;
    font-weight: 700;
}

.product-media {
    height: 140px;
    border-radius: 0;
    background: var(--panel-dark);
    border: 1px solid rgba(245, 140, 70, 0.2);
    overflow: hidden;
}

.product-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.product-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(
        135deg,
        var(--panel-dark) 0%,
        rgba(245, 140, 70, 0.1) 100%
    );
}

.product-info {
    display: flex;
    flex-direction: column;
    gap: 8px;
    flex: 1;
}

.product-tags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
}

.product-tag {
    font-size: 10px;
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.product-tag:not(:last-child)::after {
    content: "/";
    margin-left: 6px;
    opacity: 0.5;
}

.product-name {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
}

.product-pricing {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: auto;
}

.original-price {
    font-size: 13px;
    color: var(--muted);
    text-decoration: line-through;
}

.price {
    font-weight: 700;
    font-size: 15px;
}

.price.discounted {
    color: var(--accent);
}

.product-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
}

.view-details-btn {
    flex: 1;
    background: transparent;
    border: 1px solid var(--border-strong);
    border-radius: 0;
    padding: 6px 10px;
    font-family: inherit;
    font-size: 12px;
    font-weight: 600;
    color: var(--text);
    text-decoration: none;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s ease;
}

.view-details-btn:hover {
    transform: translateY(-1px);
    background: rgba(245, 140, 70, 0.1);
}

.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 16px;
    padding: 16px 0;
}

.page-indicator {
    font-size: 14px;
    color: var(--muted);
}

@media (max-width: 700px) {
    .page-header {
        flex-direction: column;
    }

    .filters-container {
        flex-direction: column;
        align-items: flex-start;
    }
    
    .search-input {
        width: 100%;
    }

    .product-grid {
        grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    }
}
</style>
