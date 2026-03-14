<script setup lang="ts">
import { RouterLink } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"

const featuredProducts = [
    {
        name: "Ember Foil Pack",
        price: "120 Coins",
        note: "Six cards, one foil guaranteed.",
        tag: "Hot",
    },
    {
        name: "Arcane Deck Box",
        price: "260 Coins",
        note: "Rigid shell, magnetic latch.",
        tag: "Gear",
    },
    {
        name: "Guild Token Set",
        price: "90 Coins",
        note: "24 acrylic counters.",
        tag: "Utility",
    },
    {
        name: "Boss Fight Playmat",
        price: "180 Coins",
        note: "Stitched edges, matte finish.",
        tag: "Mat",
    },
]
</script>

<template>
    <AppLayout>
        <section class="hero">
            <div class="hero-copy">
                <p class="eyebrow">koopify</p>
                <h1>Do IRL trading in this store</h1>
                <p class="lead">A fictional store to sell fictional videogame items.</p>
                <div class="hero-actions">
                    <RouterLink class="primary" to="/catalog">Shop now</RouterLink>
                </div>
            </div>
        </section>

        <section class="section">
            <div class="section-header">
                <h3>Featured stock</h3>
                <p class="subtle">Small rotation, big pulls.</p>
            </div>
            <div class="product-grid">
                <article v-for="product in featuredProducts" :key="product.name" class="product">
                    <div class="product-tag">{{ product.tag }}</div>
                    <div class="product-media"></div>
                    <div class="product-info">
                        <h4>{{ product.name }}</h4>
                        <p>{{ product.note }}</p>
                        <div class="product-footer">
                            <span class="price">{{ product.price }}</span>
                            <button class="primary small">Add</button>
                        </div>
                    </div>
                </article>
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

.ghost,
.primary {
    border: 1px solid transparent;
    border-radius: 0;
    padding: 10px 16px;
    font-family: inherit;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s ease;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
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

.primary.small {
    padding: 8px 14px;
    font-size: 12px;
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
