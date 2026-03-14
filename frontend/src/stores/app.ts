import { defineStore } from "pinia"

export const useAppStore = defineStore("app", {
    state: () => ({
        cartCount: 2,
    }),
    actions: {
        incrementCart(amount = 1) {
            this.cartCount += amount
        },
        resetCart() {
            this.cartCount = 0
        },
    },
})
