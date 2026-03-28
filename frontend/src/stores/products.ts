import { defineStore } from "pinia"
import { apiClient } from "@/api/client"
import { useAuthStore } from "@/stores/auth"

export type CategoryResponse = {
    id: number
    name: string
    slug: string
}

export type ProductResponse = {
    id: number
    name: string
    slug: string
    description: string
    imageUrl?: string | null
    priceEurCents: number
    discountPercent?: number | null
    inventoryCount: number
    inStock: boolean
    isActive: boolean
    createdAt: string
    updatedAt: string
    categories: CategoryResponse[]
}

export type CreateProductPayload = {
    name: string
    slug: string
    description: string
    imageUrl?: string | null
    priceEurCents: number
    discountPercent?: number | null
    inventoryCount: number
    inStock: boolean
    isActive: boolean
    categoryIds: number[]
}

export type UpdateProductPayload = CreateProductPayload

export type CreateCategoryPayload = {
    name: string
    slug: string
}

function authToken(): string {
    return useAuthStore().token
}

// NOTE: it is probably not efficient to use a store for this, as these hold no state. But I am just
// copying the auth module and dont feel like creating brand new code for this. Even <leader>gq to
// auto-warp comments does not work for typescript in neovim, it is geniuenly harrowing to write
// code in a language whose LSP is so underdeveloped
export const useProductStore = defineStore("products", {
    actions: {
        // Public
        async listActiveProducts(categorySlug?: string): Promise<ProductResponse[]> {
            const query = categorySlug ? `?category=${encodeURIComponent(categorySlug)}` : ""
            return apiClient.get<ProductResponse[]>(`/public_api/v1/products${query}`)
        },

        async getProductBySlug(slug: string): Promise<ProductResponse> {
            return apiClient.get<ProductResponse>(`/public_api/v1/products/${encodeURIComponent(slug)}`)
        },

        async listCategories(): Promise<CategoryResponse[]> {
            return apiClient.get<CategoryResponse[]>("/public_api/v1/categories")
        },

        // Admin
        async adminListAllProducts(): Promise<ProductResponse[]> {
            return apiClient.get<ProductResponse[]>("/api/v1/products", {
                authToken: authToken(),
            })
        },

        async adminGetProduct(id: number): Promise<ProductResponse> {
            return apiClient.get<ProductResponse>(`/api/v1/products/${id}`, {
                authToken: authToken(),
            })
        },

        async adminCreateProduct(payload: CreateProductPayload): Promise<ProductResponse> {
            return apiClient.post<ProductResponse>("/api/v1/products", payload, {
                authToken: authToken(),
            })
        },

        async adminUpdateProduct(id: number, payload: UpdateProductPayload): Promise<ProductResponse> {
            return apiClient.put<ProductResponse>(`/api/v1/products/${id}`, payload, {
                authToken: authToken(),
            })
        },

        async adminDeleteProduct(id: number): Promise<void> {
            return apiClient.delete<void>(`/api/v1/products/${id}`, {
                authToken: authToken(),
            })
        },

        async adminCreateCategory(payload: CreateCategoryPayload): Promise<CategoryResponse> {
            return apiClient.post<CategoryResponse>("/api/v1/categories", payload, {
                authToken: authToken(),
            })
        },
    },
})
