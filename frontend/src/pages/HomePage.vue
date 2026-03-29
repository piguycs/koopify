<script setup lang="ts">
import AppLayout from "@/layouts/AppLayout.vue"
import Button from "@/components/Button.vue"
import router from "@/router";

function gotoShop() {
    router.push({name: "catalogue"})
}
</script>

<template>
    <AppLayout>
        <section class="hero">
            <div class="hero-copy">
                <p class="eyebrow">koopify</p>
                <h1>Buy various gaming related products here</h1>
                <p class="lead">A fictional store to sell items</p>
                <Button variant="primary" @click="gotoShop">Shop now</Button>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.hero {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 28px;
    align-items: center;
}

.hero-copy h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: clamp(30px, 4vw, 44px);
    margin: 12px 0 12px;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 3px;
    font-size: 11px;
    font-weight: 600;
    opacity: 0.7;
}

.lead {
    font-size: 15px;
    line-height: 1.6;
    color: var(--muted);
}

.hero-actions {
    display: flex;
    gap: 12px;
    margin-top: 20px;
    flex-wrap: wrap;
}

.section {
    display: flex;
    flex-direction: column;
    gap: 18px;
}

.section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}

.section-header h3 {
    font-family: "Rajdhani", sans-serif;
    font-size: 24px;
    margin: 0;
}

.subtle {
    color: var(--muted);
    margin: 0;
    font-size: 13px;
}

.product-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 16px;
}

.product {
    background: var(--panel);
    border-radius: 0;
    padding: 16px;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 14px;
    border: 1px solid var(--border);
}

.product-tag {
    position: absolute;
    top: 14px;
    left: 14px;
    background: var(--accent);
    color: var(--bg);
    padding: 3px 10px;
    border-radius: 0;
    font-size: 10px;
    font-weight: 700;
}

.product-media {
    height: 140px;
    border-radius: 0;
    background: var(--panel-dark);
    border: 1px solid rgba(245, 140, 70, 0.2);
}

.product-info h4 {
    margin: 0 0 6px;
    font-size: 16px;
}

.product-info p {
    margin: 0;
    color: var(--muted);
    font-size: 13px;
}

.product-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 10px;
}

.price {
    font-weight: 700;
}

.ghost:hover,
.primary:hover {
    transform: translateY(-1px);
}

@media (max-width: 600px) {
    .hero-actions {
        flex-direction: column;
        align-items: stretch;
    }
}
</style>
