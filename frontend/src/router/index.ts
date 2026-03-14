import { createRouter, createWebHistory } from "vue-router"
import HomePage from "@/pages/HomePage.vue"
import CatalogPage from "@/pages/CatalogPage.vue"
import NotFoundPage from "@/pages/NotFoundPage.vue"

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: "/", name: "home", component: HomePage },
        { path: "/catalog", name: "catalog", component: CatalogPage },
        { path: "/:pathMatch(.*)*", name: "not-found", component: NotFoundPage },
    ],
})

export default router
