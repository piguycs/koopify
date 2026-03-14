import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import HomePage from "@/pages/HomePage.vue"
import CatalogPage from "@/pages/CatalogPage.vue"
import SignInPage from "@/pages/SignInPage.vue"
import SignUpPage from "@/pages/SignUpPage.vue"
import AccountPage from "@/pages/AccountPage.vue"
import NotFoundPage from "@/pages/NotFoundPage.vue"
import ForbiddenPage from "@/pages/ForbiddenPage.vue"

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: "/", name: "home", component: HomePage },
        { path: "/catalog", name: "catalog", component: CatalogPage },

        { path: "/sign-in", name: "sign-in", component: SignInPage },
        { path: "/sign-up", name: "sign-up", component: SignUpPage },

        // authenticated route
        { path: "/account", name: "account", component: AccountPage, meta: { requiresAuth: true } },

        { path: "/forbidden", name: "forbidden", component: ForbiddenPage },
        { path: "/:pathMatch(.*)*", name: "not-found", component: NotFoundPage },
    ],
})

router.beforeEach((to) => {
    const authStore = useAuthStore()
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        return { name: "forbidden", query: { redirect: to.fullPath } }
    }
    return true
})

export default router
