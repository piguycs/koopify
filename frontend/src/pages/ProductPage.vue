<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute, useRouter, RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import { useProductStore, type ProductResponse } from "@/stores/products"
import { useCartStore } from "@/stores/cart"
import { useAuthStore } from "@/stores/auth"

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()
const cartStore = useCartStore()
const authStore = useAuthStore()

const product = ref<ProductResponse | null>(null)
const isLoading = ref(true)
const errorMessage = ref("")
const quantity = ref(1)

const slug = route.params.slug as string

const currentCartQty = computed(() =>
    product.value ? cartStore.getItemQuantity(product.value.id) : 0,
)
const maxAddableQty = computed(() => Math.max(0, 10 - currentCartQty.value))
const isAtCartMax = computed(() => currentCartQty.value >= 10)

async function loadProduct() {
    isLoading.value = true
    errorMessage.value = ""
    try {
        product.value = await productStore.getProductBySlug(slug)
    } catch {
        errorMessage.value = "Product not found"
    } finally {
        isLoading.value = false
    }
}

function formatPrice(cents: number): string {
    return `€${(cents / 100).toFixed(2)}`
}

function getDiscountedPrice(priceCents: number, discountPercent: number | null): number {
    if (!discountPercent) return priceCents
    return Math.round(priceCents * (1 - discountPercent / 100))
}

function addToCart() {
    if (!product.value || !product.value.inStock) return
    const qtyToAdd = Math.min(quantity.value, maxAddableQty.value)
    if (qtyToAdd <= 0) return
    cartStore.addItem(product.value, qtyToAdd)
}

function gotoModifyProduct() {
    const id = product.value?.id
    let path = "/admin/products/" + id

    let backUrl: string | undefined = "/product/" + product.value?.slug
    if (id) router.push({ path, state: { backUrl } })
}

function incrementQuantity() {
    if (product.value && quantity.value < maxAddableQty.value) {
        quantity.value++
    }
}

function decrementQuantity() {
    if (quantity.value > 1) {
        quantity.value--
    }
}

// attempting to make functions pure, after writhing 10000 lines of code with impure functions
// not using a strict language gave me the freedom which I only now realise I actually hate
function isInStock(product: ProductResponse): boolean {
    return product.inStock && product.inventoryCount > 0
}

function inventoryCount(product: ProductResponse): number {
    if (!product.inStock) return 0
    return product.inventoryCount
}

function showBtnSection(product: ProductResponse): boolean {
    return (authStore.currentUser?.admin && !isInStock(product)) || isInStock(product)
}

onMounted(() => {
    loadProduct()
})
</script>

<template>
    <AppLayout>
        <section class="back-link">
            <RouterLink to="/catalogue" class="back-btn"> &larr; Back to Catalogue </RouterLink>
        </section>

        <section v-if="isLoading" class="loading">
            <p>Loading product...</p>
        </section>

        <section v-else-if="errorMessage" class="error">
            <p>{{ errorMessage }}</p>
        </section>

        <section v-else-if="product" class="product-detail">
            <div class="product-image-section">
                <div class="product-image-container">
                    <img
                        v-if="product.imageUrl"
                        :src="product.imageUrl"
                        :alt="product.name"
                        class="product-image"
                    />
                    <div v-else class="product-placeholder"></div>
                    <div v-if="product.discountPercent" class="discount-badge">
                        -{{ product.discountPercent }}%
                    </div>
                </div>
            </div>

            <div class="product-info-section">
                <div class="product-meta">
                    <span v-for="cat in product.categories" :key="cat.id" class="category-tag">
                        {{ cat.name }}
                    </span>
                </div>

                <h1 class="product-name">{{ product.name }}</h1>

                <div class="product-pricing">
                    <span v-if="product.discountPercent" class="original-price">
                        {{ formatPrice(product.priceEurCents) }}
                    </span>
                    <span class="current-price" :class="{ discounted: product.discountPercent }">
                        {{
                            formatPrice(
                                getDiscountedPrice(
                                    product.priceEurCents,
                                    product.discountPercent ?? null,
                                ),
                            )
                        }}
                    </span>
                </div>

                <p class="product-description">{{ product.description }}</p>

                <div class="stock-status">
                    <span
                        :class="[
                            'stock-indicator',
                            isInStock(product) ? 'in-stock' : 'out-of-stock',
                        ]"
                    ></span>
                    {{ isInStock(product) ? "In Stock" : "Out of Stock" }}
                </div>

                <div class="add-to-cart-section" v-if="showBtnSection(product)">
                    <div v-if="isInStock(product) && !isAtCartMax" class="quantity-selector">
                        <button
                            class="qty-btn"
                            @click="decrementQuantity"
                            :disabled="quantity <= 1"
                        >
                            -
                        </button>
                        <span class="qty-value">{{ quantity }}</span>
                        <button
                            class="qty-btn"
                            @click="incrementQuantity"
                            :disabled="quantity >= maxAddableQty || isAtCartMax"
                        >
                            +
                        </button>
                    </div>

                    <Button
                        v-if="isInStock(product)"
                        variant="primary"
                        @click="addToCart"
                        :disabled="isAtCartMax"
                    >
                        {{ isAtCartMax ? "Cart Full" : "Add to Cart" }}
                    </Button>
                    <Button
                        v-if="authStore.currentUser?.admin"
                        variant="ghost"
                        @click="gotoModifyProduct"
                    >
                        Modify this product
                    </Button>
                </div>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.back-link {
    margin-bottom: 24px;
}

.back-btn {
    color: var(--muted);
    text-decoration: none;
    font-size: 14px;
    transition: color 0.2s ease;
}

.back-btn:hover {
    color: var(--text);
}

.loading,
.error {
    text-align: center;
    padding: 48px 0;
    color: var(--muted);
}

.error {
    color: #f38b8b;
}

.product-detail {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 48px;
    align-items: start;
}

.product-image-section {
    position: sticky;
    top: 24px;
}

.product-image-container {
    position: relative;
    aspect-ratio: 1;
    background: var(--panel);
    border: 1px solid var(--border);
}

.product-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.product-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, var(--panel-dark) 0%, rgba(245, 140, 70, 0.1) 100%);
}

.discount-badge {
    position: absolute;
    top: 16px;
    right: 16px;
    background: var(--accent);
    color: var(--bg);
    padding: 6px 14px;
    border-radius: 0;
    font-size: 14px;
    font-weight: 700;
}

.product-info-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.product-meta {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

.category-tag {
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--muted);
    padding: 4px 10px;
    border: 1px solid var(--border);
}

.product-name {
    font-family: "Rajdhani", sans-serif;
    font-size: 36px;
    margin: 0;
    line-height: 1.1;
}

.product-pricing {
    display: flex;
    align-items: center;
    gap: 12px;
}

.original-price {
    font-size: 18px;
    color: var(--muted);
    text-decoration: line-through;
}

.current-price {
    font-size: 28px;
    font-weight: 700;
}

.current-price.discounted {
    color: var(--accent);
}

.product-description {
    font-size: 15px;
    line-height: 1.7;
    color: var(--text);
    margin: 0;
    white-space: pre-line;
}

.stock-status {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
}

.stock-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
}

.stock-indicator.in-stock {
    background: #7ec87e;
}

.stock-indicator.out-of-stock {
    background: #f38b8b;
}

.stock-count {
    color: var(--muted);
    font-weight: 400;
}

.add-to-cart-section {
    display: flex;
    gap: 16px;
    align-items: center;
    margin-top: 16px;
    padding-top: 24px;
    border-top: 1px solid var(--border);
}

.quantity-selector {
    display: flex;
    align-items: center;
    border: 1px solid var(--border);
}

.qty-btn {
    width: 38px;
    height: 38px;
    background: transparent;
    border: none;
    color: var(--text);
    font-size: 18px;
    cursor: pointer;
    transition: background 0.2s ease;
}

.qty-btn:hover:not(:disabled) {
    background: rgba(245, 140, 70, 0.1);
}

.qty-btn:disabled {
    opacity: 0.4;
    cursor: default;
}

.qty-value {
    width: 40px;
    text-align: center;
    font-weight: 600;
    border-left: 1px solid var(--border);
    border-right: 1px solid var(--border);
    line-height: 38px;
}

.cart-qty-info {
    font-size: 13px;
    color: var(--muted);
    margin-right: 8px;
}

@media (max-width: 800px) {
    .product-detail {
        grid-template-columns: 1fr;
        gap: 24px;
    }

    .product-image-section {
        position: static;
    }

    .product-name {
        font-size: 28px;
    }

    .add-to-cart-section {
        flex-direction: column;
        align-items: stretch;
    }

    .quantity-selector {
        justify-content: center;
    }
}
</style>
