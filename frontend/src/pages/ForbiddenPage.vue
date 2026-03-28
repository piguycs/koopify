<script setup lang="ts">
import { computed } from "vue"
import { RouterLink, useRoute } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"

const route = useRoute()

// NOTE: we need to be sure we do not redirect the user to unauthorised pages
// unauthorised pages should result in a 404
const redirectTarget = computed(() => {
    const redirect = route.query.redirect
    return typeof redirect === "string" && redirect.length > 0 ? redirect : "/catalogue"
})
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <p class="eyebrow">Access blocked</p>
            <h1>Authentication required</h1>
            <p class="subtle">Sign in to continue to the requested page.</p>
            <div class="actions">
                <RouterLink
                    class="primary"
                    :to="{ path: '/sign-in', query: { redirect: redirectTarget } }"
                >
                    Sign in
                </RouterLink>
                <RouterLink class="ghost" to="/">Back to home</RouterLink>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 3px;
    font-size: 11px;
    font-weight: 600;
    opacity: 0.7;
    margin: 0;
}

.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0;
}

.subtle {
    margin: 0;
    color: var(--muted);
}

.actions {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 8px;
}

.primary {
    border: 1px solid transparent;
    border-radius: 0;
    padding: 12px 16px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s ease;
    background: #2a1c16;
    color: var(--text);
    border-color: rgba(245, 140, 70, 0.6);
    text-decoration: none;
}

.primary:hover {
    transform: translateY(-1px);
}

.ghost {
    border: 1px solid var(--border-strong);
    border-radius: 0;
    padding: 10px 16px;
    width: fit-content;
    text-decoration: none;
    color: var(--text);
}
</style>
