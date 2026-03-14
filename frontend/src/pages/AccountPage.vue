<script setup lang="ts">
import { computed } from "vue"
import { useAuthStore } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import EditableField from "@/components/EditableField.vue"
import ActionCard from "@/components/ActionCard.vue"

const authStore = useAuthStore()

const userDisplayName = computed(() => authStore.currentUser?.displayName)
const userEmail = computed(() => authStore.currentUser?.email)

const saveDisplayName = (value: string) => {
    console.info("WIP: display name updated", value)
}

const saveEmail = (value: string) => {
    console.info("WIP: email updated", value)
}

</script>

<template>
    <AppLayout>
        <section class="account">
            <div class="account-card">
                <div class="account-header">
                    <p class="eyebrow">Account</p>
                    <h1>View account</h1>
                    <p class="subtle">Manage your profile details.</p>
                </div>

                <div class="account-content">
                    <EditableField
                        label="Display name"
                        :value="userDisplayName"
                        @save="saveDisplayName"
                    />
                    <EditableField
                        label="Email"
                        :value="userEmail"
                        @save="saveEmail"
                    />
                    <ActionCard
                        title="Reset password"
                        description="Start the password reset flow for this account."
                        button-text="Reset password"
                        @click="console.info('WIP: password reset')"
                    />
                    <ActionCard
                        title="Delete account"
                        description="Permanently remove this account and all related data."
                        button-text="Delete account"
                        danger
                        @click="console.info('WIP: delete account')"
                    />
                </div>
            </div>
        </section>
    </AppLayout>
</template>

<style scoped>
.account {
    display: flex;
    justify-content: center;
}

.account-card {
    width: min(520px, 100%);
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 28px;
}

.account-header h1 {
    font-family: "Rajdhani", sans-serif;
    margin: 8px 0 6px;
    font-size: 30px;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 3px;
    font-size: 11px;
    font-weight: 600;
    opacity: 0.7;
    margin: 0;
}

.subtle {
    color: var(--muted);
    margin: 0;
}

.account-content {
    margin-top: 22px;
    display: grid;
    gap: 16px;
}
</style>
