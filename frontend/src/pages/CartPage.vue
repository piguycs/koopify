<script setup lang="ts">
import { RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import { useCartStore } from "@/stores/cart"
import router from "@/router"

const cartStore = useCartStore()

function formatPrice(cents: number): string {
    return `€${(cents / 100).toFixed(2)}`
}

function getDiscountedPrice(priceCents: number, discountPercent: number | null): number {
    if (!discountPercent) return priceCents
    return Math.round(priceCents * (1 - discountPercent / 100))
}

function incrementQuantity(itemId: number, currentQty: number) {
    if (currentQty < 10) {
        cartStore.updateQuantity(itemId, currentQty + 1)
    }
}

function decrementQuantity(itemId: number, currentQty: number) {
    if (currentQty > 1) {
        cartStore.updateQuantity(itemId, currentQty - 1)
    }
}

function gotoCatalogue() {
    router.push({ name: "catalogue" })
}

function handleCheckout() {
    alert("Checkout not yet implemented - backend integration pending")
}
</script>

<template>
    <AppLayout>
        <section class="page-header">
            <h1>Shopping Cart</h1>
        </section>

        <section v-if="cartStore.items.length === 0" class="empty-state">
            <h2>Your cart is empty</h2>
            <p>Browse our catalogue and add some items to your cart.</p>
            <RouterLink to="/catalogue">
                <Button variant="primary">Browse Catalogue</Button>
            </RouterLink>
        </section>

        <section v-else class="cart-content">
            <div class="cart-items">
                <article v-for="item in cartStore.items" :key="item.id" class="cart-item">
                    <div class="item-image">
                        <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" />
                        <div v-else class="item-placeholder"></div>
                    </div>

                    <div class="item-details">
                        <RouterLink :to="`/product/${item.slug}`" class="item-name">
                            {{ item.name }}
                        </RouterLink>
                        <div class="item-price">
                            <span v-if="item.discountPercent" class="original-price">
                                {{ formatPrice(item.priceEurCents) }}
                            </span>
                            <span
                                class="current-price"
                                :class="{ discounted: item.discountPercent }"
                            >
                                {{
                                    formatPrice(
                                        getDiscountedPrice(
                                            item.priceEurCents,
                                            item.discountPercent ?? null,
                                        ),
                                    )
                                }}
                            </span>
                        </div>
                    </div>

                    <div class="item-quantity">
                        <button
                            class="qty-btn"
                            @click="decrementQuantity(item.id, item.quantity)"
                            :disabled="item.quantity <= 1"
                        >
                            -
                        </button>
                        <span class="qty-value">{{ item.quantity }}</span>
                        <button
                            class="qty-btn"
                            @click="incrementQuantity(item.id, item.quantity)"
                            :disabled="item.quantity >= 10"
                        >
                            +
                        </button>
                    </div>

                    <div class="item-subtotal">
                        {{
                            formatPrice(
                                getDiscountedPrice(
                                    item.priceEurCents,
                                    item.discountPercent ?? null,
                                ) * item.quantity,
                            )
                        }}
                    </div>

                    <button
                        class="remove-btn"
                        @click="cartStore.removeItem(item.id)"
                        title="Remove item"
                    >
                        ×
                    </button>
                </article>
            </div>

            <aside class="cart-summary">
                <h2>Order Summary</h2>
                <div class="summary-row">
                    <span>Subtotal ({{ cartStore.totalItems }} items)</span>
                    <span>{{ formatPrice(cartStore.totalPrice) }}</span>
                </div>
                <div class="summary-row total">
                    <span>Total</span>
                    <span>{{ formatPrice(cartStore.totalPrice) }}</span>
                </div>
                <div class="summary-actions">
                    <Button variant="primary" @click="handleCheckout"> Proceed to Checkout </Button>
                    <Button variant="ghost" @click="gotoCatalogue">Continue Shopping</Button>
                    <Button variant="link" @click="cartStore.clearCart"> Clear Cart </Button>
                </div>
            </aside>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-header h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0 0 8px;
}

.empty-state {
    text-align: center;
    padding: 64px 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
}

.empty-state h2 {
    margin: 0;
    font-size: 24px;
}

.empty-state p {
    margin: 0 0 16px;
    color: var(--muted);
}

.cart-content {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 32px;
    align-items: start;
}

.cart-items {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.cart-item {
    display: grid;
    grid-template-columns: 80px 1fr auto auto auto;
    gap: 16px;
    align-items: center;
    padding: 16px;
    background: var(--panel);
    border: 1px solid var(--border);
}

.item-image {
    width: 80px;
    height: 80px;
    background: var(--panel-dark);
    border: 1px solid var(--border);
    overflow: hidden;
}

.item-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.item-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, var(--panel-dark) 0%, rgba(245, 140, 70, 0.1) 100%);
}

.item-details {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.item-name {
    font-weight: 600;
    color: var(--text);
    text-decoration: none;
}

.item-name:hover {
    text-decoration: underline;
}

.item-price {
    display: flex;
    align-items: center;
    gap: 8px;
}

.original-price {
    font-size: 13px;
    color: var(--muted);
    text-decoration: line-through;
}

.current-price {
    font-weight: 600;
}

.current-price.discounted {
    color: var(--accent);
}

.item-quantity {
    display: flex;
    align-items: center;
    border: 1px solid var(--border);
}

.qty-btn {
    width: 32px;
    height: 32px;
    background: transparent;
    border: none;
    color: var(--text);
    font-size: 16px;
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
    width: 32px;
    text-align: center;
    font-weight: 600;
    border-left: 1px solid var(--border);
    border-right: 1px solid var(--border);
    line-height: 32px;
}

.item-subtotal {
    font-weight: 700;
    min-width: 70px;
    text-align: right;
}

.remove-btn {
    width: 32px;
    height: 32px;
    background: transparent;
    border: 1px solid transparent;
    color: var(--muted);
    font-size: 20px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
}

.remove-btn:hover {
    color: #f38b8b;
    border-color: rgba(243, 139, 139, 0.4);
    background: rgba(243, 139, 139, 0.1);
}

.cart-summary {
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 24px;
    position: sticky;
    top: 24px;
}

.cart-summary h2 {
    font-family: "Rajdhani", sans-serif;
    font-size: 20px;
    margin: 0 0 20px;
}

.summary-row {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    font-size: 14px;
}

.summary-row.total {
    border-top: 1px solid var(--border);
    margin-top: 8px;
    padding-top: 16px;
    font-weight: 700;
    font-size: 18px;
}

.summary-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 24px;
}

.summary-actions a {
    text-decoration: none;
}

@media (max-width: 800px) {
    .cart-content {
        grid-template-columns: 1fr;
    }

    .cart-summary {
        position: static;
    }

    .cart-item {
        grid-template-columns: 60px 1fr;
        gap: 12px;
    }

    .item-quantity,
    .item-subtotal {
        grid-column: 2;
    }

    .remove-btn {
        position: absolute;
        top: 12px;
        right: 12px;
    }

    .cart-item {
        position: relative;
    }
}

@media (max-width: 500px) {
    .item-image {
        width: 60px;
        height: 60px;
    }
}
</style>
