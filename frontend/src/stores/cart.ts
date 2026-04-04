import { defineStore } from "pinia"
import type { ProductResponse } from "@/stores/products"

export type CartItemData = {
    id: number
    name: string
    slug: string
    imageUrl?: string | null
    priceEurCents: number
    discountPercent?: number | null
    quantity: number
}

const STORAGE_KEY = "koopify_cart"

function loadFromStorage(): CartItemData[] {
    try {
        const stored = localStorage.getItem(STORAGE_KEY)
        return stored ? JSON.parse(stored) : []
    } catch {
        return []
    }
}

export const useCartStore = defineStore("cart", {
    state: () => ({
        items: loadFromStorage() as CartItemData[],
        pulseTrigger: 0,
    }),
    getters: {
        totalItems: state => state.items.reduce((sum, item) => sum + item.quantity, 0),
        totalPrice: state =>
            state.items.reduce((sum, item) => {
                const price = item.discountPercent
                    ? Math.round(item.priceEurCents * (1 - item.discountPercent / 100))
                    : item.priceEurCents
                return sum + price * item.quantity
            }, 0),
    },
    actions: {
        addItem(product: ProductResponse, quantity: number = 1) {
            const existing = this.items.find(item => item.id === product.id)
            if (existing) {
                existing.quantity += quantity
            } else {
                this.items.push({
                    id: product.id,
                    name: product.name,
                    slug: product.slug,
                    imageUrl: product.imageUrl,
                    priceEurCents: product.priceEurCents,
                    discountPercent: product.discountPercent,
                    quantity,
                })
            }
            this.persistToStorage()
            this.pulseTrigger++
        },
        removeItem(productId: number) {
            const index = this.items.findIndex(item => item.id === productId)
            if (index !== -1) {
                this.items.splice(index, 1)
                this.persistToStorage()
            }
        },
        updateQuantity(productId: number, quantity: number) {
            const item = this.items.find(item => item.id === productId)
            if (item) {
                if (quantity <= 0) {
                    this.removeItem(productId)
                } else {
                    item.quantity = quantity
                    this.persistToStorage()
                }
            }
        },
        clearCart() {
            this.items = []
            this.persistToStorage()
        },
        getItemQuantity(productId: number): number {
            const item = this.items.find(item => item.id === productId)
            return item?.quantity ?? 0
        },
        isAtMax(productId: number): boolean {
            return this.getItemQuantity(productId) >= 10
        },
        persistToStorage() {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(this.items))
        },
    },
})
