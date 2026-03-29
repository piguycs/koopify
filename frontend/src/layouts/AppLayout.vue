<script setup lang="ts">
import { onMounted } from "vue"
import { RouterLink, useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { useCartStore } from "@/stores/cart"
import Button from "@/components/Button.vue"

const authStore = useAuthStore()
const cartStore = useCartStore()
const router = useRouter()

const handleSignOut = () => {
    authStore.signOut()
    router.push("/")
}

onMounted(async () => {
    if (authStore.token && !authStore.currentUser) {
        try {
            await authStore.fetchCurrentUser()
        } catch {
            authStore.signOut()
        }
    }
})
</script>

<template>
    <div class="page">
        <header class="nav">
            <RouterLink class="logo" to="/">Koopify</RouterLink>
            <nav class="nav-links">
                <RouterLink to="/">Home</RouterLink>
                <RouterLink to="/catalogue">Catalogue</RouterLink>
            </nav>
            <div class="nav-actions">
                <template v-if="!authStore.isAuthenticated">
                    <RouterLink class="ghost" to="/sign-in">Sign in</RouterLink>
                </template>
                <template v-else>
                    <div class="account-menu">
                        <RouterLink class="ghost" to="/account">View Account</RouterLink>
                        <div class="account-dropdown">
                            <Button variant="ghost" type="button" @click="handleSignOut">
                                Sign out
                            </Button>
                        </div>
                    </div>
                </template>
                <RouterLink to="/cart">
                    <Button variant="primary" type="button">
                        Cart ({{ cartStore.totalItems }})
                    </Button>
                </RouterLink>
            </div>
        </header>

        <main class="content">
            <slot />
        </main>

        <footer class="footer">
            <div>
                <div class="logo">Koopify</div>
                <p>A fictional store to sell fictional videogame items</p>
            </div>
            <div class="footer-links">
                <!--
                <RouterLink to="/catalog">Inventory</RouterLink>
                <RouterLink to="/trade">Support</RouterLink>
                <RouterLink to="/">Contact</RouterLink>
                -->
            </div>
        </footer>
    </div>
</template>

<style scoped>
.page {
    min-height: 100vh;
    background: var(--bg);
    display: flex;
    flex-direction: column;
}

.nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 24px 8vw;
    border-bottom: 1px solid var(--border);
    background: rgba(13, 12, 11, 0.92);
}

.logo {
    font-family: "Rajdhani", sans-serif;
    font-weight: 700;
    letter-spacing: 2px;
    font-size: 20px;
}

.nav-links {
    display: flex;
    gap: 18px;
}

.nav-links a {
    text-decoration: none;
    font-weight: 500;
    font-size: 14px;
    opacity: 0.8;
}

.nav-links a.router-link-active {
    opacity: 1;
}

.nav-actions {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
}

.account-menu {
    position: relative;
    display: inline-flex;
}

.account-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 8px 6px;
    display: none;
    min-width: 100%;
    z-index: 10;
}

.account-dropdown .ghost {
    width: 100%;
    justify-content: center;
}

.account-menu:hover .account-dropdown,
.account-menu:focus-within .account-dropdown {
    display: block;
}

.ghost,
.primary {
    border: 1px solid transparent;
    border-radius: 0;
    padding: 10px 16px;
    font-family: inherit;
    font-weight: 600;
    font-size: 13px;
    cursor: pointer;
    transition:
        transform 0.2s ease,
        box-shadow 0.2s ease;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 38px;
}

.ghost {
    background: transparent;
    border-color: var(--border-strong);
    color: var(--text);
}

.primary {
    background: #2a1c16;
    color: var(--text);
    border-color: rgba(245, 140, 70, 0.6);
}

.ghost:hover,
.primary:hover {
    transform: translateY(-1px);
}

.content {
    padding: 36px 8vw 72px;
    display: flex;
    flex-direction: column;
    gap: 48px;
    flex: 1;
}

.footer {
    border-top: 1px solid var(--border);
    color: var(--muted);
    padding: 28px 8vw;
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 16px;
}

.footer-links {
    display: flex;
    gap: 16px;
}

.footer a {
    text-decoration: none;
    font-size: 13px;
}

@media (max-width: 900px) {
    .nav {
        flex-direction: column;
        gap: 14px;
    }

    .nav-links {
        flex-wrap: wrap;
        justify-content: center;
    }

    .nav-actions {
        width: 100%;
        justify-content: center;
    }
}

@media (max-width: 600px) {
    .content {
        padding: 24px 6vw 56px;
    }
}
</style>
