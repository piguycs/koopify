import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import HomePage from "@/pages/HomePage.vue"
import CataloguePage from "@/pages/CataloguePage.vue"
import ProductPage from "@/pages/ProductPage.vue"
import SignInPage from "@/pages/SignInPage.vue"
import SignUpPage from "@/pages/SignUpPage.vue"
import AccountPage from "@/pages/AccountPage.vue"
import NotFoundPage from "@/pages/NotFoundPage.vue"
import ForbiddenPage from "@/pages/ForbiddenPage.vue"
import AdminProductsPage from "@/pages/AdminProductsPage.vue"
import AdminProductEditPage from "@/pages/AdminProductEditPage.vue"
import AdminUsersPage from "@/pages/AdminUsersPage.vue"
import AdminOrdersPage from "@/pages/AdminOrdersPage.vue"
import CheckoutComplete from "@/pages/CheckoutComplete.vue"
import CartPage from "@/pages/CartPage.vue"

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: "/", name: "home", component: HomePage },
        { path: "/product/:slug", name: "product", component: ProductPage },
        { path: "/catalogue", name: "catalogue", component: CataloguePage },

        { path: "/sign-in", name: "sign-in", component: SignInPage },
        { path: "/sign-up", name: "sign-up", component: SignUpPage },

        // authenticated route
        { path: "/account", name: "account", component: AccountPage, meta: { requiresAuth: true } },

        {
            path: "/checkout-complete",
            name: "checkout-complete",
            component: CheckoutComplete,
            meta: { requiresAuth: true },
        },

        { path: "/cart", name: "cart", component: CartPage },

        {
            path: "/admin/products",
            name: "admin-products",
            component: AdminProductsPage,
            meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
            path: "/admin/products/new",
            name: "admin-products-new",
            component: AdminProductEditPage,
            meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
            path: "/admin/products/:id",
            name: "admin-products-edit",
            component: AdminProductEditPage,
            meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
            path: "/admin/users",
            name: "admin-users",
            component: AdminUsersPage,
            meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
            path: "/admin/orders",
            name: "admin-orders",
            component: AdminOrdersPage,
            meta: { requiresAuth: true, requiresAdmin: true },
        },

        { path: "/forbidden", name: "forbidden", component: ForbiddenPage },
        { path: "/:pathMatch(.*)*", name: "not-found", component: NotFoundPage },
    ],
})

router.beforeEach(async to => {
    const authStore = useAuthStore()
    if (authStore.token && !authStore.currentUser) {
        try {
            await authStore.fetchCurrentUser()
        } catch {
            authStore.signOut()
        }
    }

    if (to.meta.requiresAdmin && !authStore.currentUser?.admin) {
        return { name: "not-found" }
    }

    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        return { name: "forbidden", query: { redirect: to.fullPath } }
    }
    return true
})

export default router
