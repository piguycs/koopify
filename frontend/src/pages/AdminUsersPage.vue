<!-- TODO: add pagination -->

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useAuthStore } from "@/stores/auth"
import AppLayout from "@/layouts/AppLayout.vue"
import ModalDialog from "@/components/ModalDialog.vue"
import { ApiError } from "@/api/client"
import type { UserResponse } from "@/stores/auth"

enum AdminModalState {
    Disabled,
    Promoting,
    Demoting,
    Deleting,
}

const authStore = useAuthStore()

const users = ref<UserResponse[]>([])
const isLoading = ref(false)
const errorMessage = ref("")
const successMessage = ref("")

// these are for filtering/sorting the rows
const searchQuery = ref("")
const sortKey = ref<keyof UserResponse | null>(null)
const sortDirection = ref<'asc' | 'desc'>('asc')

// Admin toggle modal state
const adminModalState = ref(AdminModalState.Disabled)
const isDemoting = computed(() => adminModalState.value == AdminModalState.Demoting)
const isModalOpen = computed(() => adminModalState.value != AdminModalState.Disabled)
const userToToggle = ref<UserResponse | null>(null)

function closeAdminModal() {
    adminModalState.value = AdminModalState.Disabled
    userToToggle.value = null
}

async function confirmAdminToggle() {
    if (!userToToggle.value) return
    
    const user = userToToggle.value
    closeAdminModal()
    await performToggleAdmin(user)
}

// Delete modal state
const userToDelete = ref<UserResponse | null>(null)
const isDeleteModalOpen = computed(() => adminModalState.value === AdminModalState.Deleting)

function openDeleteModal(user: UserResponse) {
    // Prevent self-deletion
    if (user.id === authStore.currentUser?.id) {
        errorMessageToast(null, "You cannot delete yourself")
        return
    }
    
    userToDelete.value = user
    adminModalState.value = AdminModalState.Deleting
}

function closeDeleteModal() {
    adminModalState.value = AdminModalState.Disabled
    userToDelete.value = null
}

async function confirmDelete() {
    if (!userToDelete.value) return
    
    const user = userToDelete.value
    closeDeleteModal()
    
    try {
        const updated = await authStore.requestUserDeletionAdmin(user.id)
        
        // Update the user in the list
        const index = users.value.findIndex(u => u.id === updated.id)
        if (index !== -1) {
            users.value[index] = updated
        }
        
        successMessageToast(`User ${user.displayName} scheduled for deletion`)
    } catch (err) {
        errorMessageToast(err, "Failed to schedule user deletion")
    }
}

type SortableKey = 'id' | 'displayName' | 'email' | 'admin' | 'deletionScheduledAt'

function clearSearch() {
    searchQuery.value = ""
}

const filteredUsers = computed(() => {
    if (!searchQuery.value.trim()) return users.value
    
    const query = searchQuery.value.toLowerCase().trim()
    return users.value.filter(user => 
        user.displayName.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query)
    )
})

function toggleSort(key: SortableKey) {
    if (sortKey.value === key) {
        // Cycle: asc -> desc -> null (default)
        if (sortDirection.value === 'asc') {
            sortDirection.value = 'desc'
        } else {
            sortKey.value = null
        }
    } else {
        // New column, default to ascending
        sortKey.value = key
        sortDirection.value = 'asc'
    }
}

function getSortIndicator(key: SortableKey): string {
    if (sortKey.value !== key) return ''
    return sortDirection.value === 'asc' ? ' ▲' : ' ▼'
}

const sortedUsers = computed(() => {
    if (!sortKey.value) return filteredUsers.value

    const key = sortKey.value
    const dir = sortDirection.value

    return [...filteredUsers.value].sort((a, b) => {
        let valA = a[key]
        let valB = b[key]

        // Handle null/undefined values
        if (valA === null || valA === undefined) valA = ''
        if (valB === null || valB === undefined) valB = ''

        // Convert to comparable values
        if (typeof valA === 'string' && typeof valB === 'string') {
            valA = valA.toLowerCase()
            valB = valB.toLowerCase()
        }

        if (valA < valB) return dir === 'asc' ? -1 : 1
        if (valA > valB) return dir === 'asc' ? 1 : -1
        return 0
    })
})

function clearMessages() {
    errorMessage.value = ""
    successMessage.value = ""
}

// TODO: perhaps a proper toast someday, not today though
function successMessageToast(message: string) {
    successMessage.value = message
}

// TODO: perhaps a proper toast someday, not today though
function errorMessageToast(err: unknown, message: string) {
    if (err instanceof ApiError) {
        errorMessage.value = err.message
    } else {
        errorMessage.value = message
    }
}

async function loadUsers() {
    isLoading.value = true
    clearMessages()

    try {
        users.value = await authStore.listUsers()
    } catch (err) {
        errorMessageToast(err, "Failed to load users")
    } finally {
        isLoading.value = false
    }
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
        timeZoneName: 'short'
    })
}

async function performToggleAdmin(user: UserResponse) {
    clearMessages()

    try {
        const updated = await authStore.updateUserAdmin(
            user.id,
            !user.admin
        )

        const index = users.value.findIndex(u => u.id === updated.id)
        if (index !== -1) {
            users.value[index] = updated
        }

        successMessageToast(`User ${updated.admin ? "promoted to" : "demoted from"} admin`)
    } catch (err) {
        errorMessageToast(err, "Failed to update admin status")
    }
}

function toggleAdmin(user: UserResponse) {
    // Prevent self-demotion
    if (user.admin && user.id === authStore.currentUser?.id) {
        errorMessageToast(null, "You cannot demote yourself")
        return
    }
    
    userToToggle.value = user
    if (adminModalState.value == AdminModalState.Disabled) {
        if (user.admin) {
            adminModalState.value = AdminModalState.Demoting
        } else {
            adminModalState.value = AdminModalState.Promoting
        }
    }
}

async function triggerReset(user: UserResponse) {
    clearMessages()

    try {
        const resp = await authStore.triggerPasswordReset(user.id)
        successMessageToast(resp.message)
    } catch (err) {
        errorMessageToast(err, "Failed to trigger password reset")
    }
}

onMounted(() => {
    loadUsers()
})
</script>

<template>
    <AppLayout>
        <section class="page-section">
            <p class="eyebrow">Admin</p>
            <h1>Manage Users</h1>
            <p class="subtle">View and manage user accounts.</p>
        </section>

        <div class="admin-content">
            <div v-if="successMessage" class="toast success">{{ successMessage }}</div>
            <div v-if="errorMessage" class="toast error">{{ errorMessage }}</div>

            <div v-if="isLoading" class="loading">Loading users...</div>
            
            <div v-else class="table-wrapper">
                <div class="table-controls">
                    <div class="search-box">
                        <input
                            v-model="searchQuery"
                            type="text"
                            placeholder="Search by name or email..."
                            class="search-input"
                        />
                        <button 
                            v-if="searchQuery" 
                            class="clear-btn"
                            @click="clearSearch"
                        >
                            Clear
                        </button>
                    </div>
                    <div class="results-count">
                        {{ sortedUsers.length }} user{{ sortedUsers.length === 1 ? '' : 's' }}
                        <span v-if="searchQuery">found</span>
                    </div>
                </div>
                
                <div class="table-container">
                <table class="users-table">
                    <thead>
                        <tr>
                            <th class="col-id sortable" @click="toggleSort('id')">
                                ID{{ getSortIndicator('id') }}
                            </th>
                            <th class="col-name sortable" @click="toggleSort('displayName')">
                                Display Name{{ getSortIndicator('displayName') }}
                            </th>
                            <th class="col-email sortable" @click="toggleSort('email')">
                                Email{{ getSortIndicator('email') }}
                            </th>
                            <th class="col-admin sortable" @click="toggleSort('admin')">
                                Admin{{ getSortIndicator('admin') }}
                            </th>
                            <th class="col-deletion sortable" @click="toggleSort('deletionScheduledAt')">
                                Deletion{{ getSortIndicator('deletionScheduledAt') }}
                            </th>
                            <th class="col-actions">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in sortedUsers" :key="user.id">
                            <td class="col-id">{{ user.id }}</td>
                            <td class="col-name">
                                <div class="truncated" :title="user.displayName">
                                    {{ user.displayName }}
                                </div>
                            </td>
                            <td class="col-email">
                                <div class="truncated" :title="user.email">
                                    {{ user.email }}
                                </div>
                            </td>
                            <td class="col-admin">
                                <span v-if="user.admin" class="admin-badge">ADMIN</span>
                                <span v-else class="dash">—</span>
                            </td>
                            <td class="col-deletion">
                                <span v-if="user.deletionScheduledAt">
                                    {{ formatDate(user.deletionScheduledAt) }}
                                </span>
                                <span v-else class="dash">—</span>
                            </td>
                            <td class="col-actions">
                                <div class="action-buttons">
                                    <button 
                                        class="action-btn"
                                        :class="{ danger: user.admin }"
                                        @click="toggleAdmin(user)"
                                    >
                                        {{ user.admin ? "Demote" : "Promote" }}
                                    </button>
                                    <button 
                                        class="action-btn"
                                        @click="triggerReset(user)"
                                    >
                                        Reset
                                    </button>
                                    <button 
                                        class="action-btn danger"
                                        @click="openDeleteModal(user)"
                                    >
                                        Delete
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            </div>
        </div>

        <ModalDialog
            :open="isModalOpen"
            :title="isDemoting ? 'Confirm Demote' : 'Confirm Promotion'"
            :description="isDemoting 
                ? 'Are you sure you want to demote this user from admin? They will lose access to admin features immediately.' 
                : 'Are you sure you want to promote this user to admin? They will gain access to all admin features.'"
            @close="closeAdminModal"
        >
            <template #actions>
                <button class="action-btn ghost" @click="closeAdminModal">
                    Cancel
                </button>
                <button 
                    class="action-btn" 
                    :class="isDemoting ? 'danger' : 'primary'"
                    @click="confirmAdminToggle"
                >
                    {{ isDemoting ? 'Demote' : 'Promote' }}
                </button>
            </template>
        </ModalDialog>

        <ModalDialog
            :open="isDeleteModalOpen"
            title="Confirm Deletion"
            description="Are you sure you want to schedule this user for deletion? They will have 24 hours to cancel the deletion by logging in."
            @close="closeDeleteModal"
        >
            <template #actions>
                <button class="action-btn ghost" @click="closeDeleteModal">
                    Cancel
                </button>
                <button class="action-btn danger" @click="confirmDelete">
                    Delete
                </button>
            </template>
        </ModalDialog>
    </AppLayout>
</template>

<style scoped>
.page-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.page-section h1 {
    font-family: "Rajdhani", sans-serif;
    font-size: 32px;
    margin: 0;
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

.admin-content {
    margin-top: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.table-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.table-controls {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}

.search-box {
    display: flex;
    align-items: center;
    gap: 8px;
}

.search-input {
    background: var(--panel);
    border: 1px solid var(--border);
    padding: 8px 12px;
    font-family: inherit;
    font-size: 14px;
    color: var(--text);
    min-width: 240px;
}

.search-input::placeholder {
    color: var(--muted);
}

.search-input:focus {
    outline: none;
    border-color: rgba(245, 140, 70, 0.6);
}

.clear-btn {
    border: 1px solid var(--border);
    background: transparent;
    color: var(--muted);
    padding: 8px 12px;
    font-family: inherit;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
}

.clear-btn:hover {
    border-color: var(--border-strong);
    color: var(--text);
}

.results-count {
    color: var(--muted);
    font-size: 14px;
}

.toast {
    padding: 12px 16px;
    font-size: 14px;
}

.toast.success {
    background: rgba(139, 243, 139, 0.1);
    border: 1px solid rgba(139, 243, 139, 0.3);
    color: #8bf38b;
}

.toast.error {
    background: rgba(243, 139, 139, 0.1);
    border: 1px solid rgba(243, 139, 139, 0.3);
    color: #f38b8b;
}

.loading {
    padding: 40px;
    text-align: center;
    color: var(--muted);
}

.table-container {
    background: var(--panel);
    border: 1px solid var(--border);
    overflow-x: auto;
}

.users-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
}

.users-table th,
.users-table td {
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid var(--border);
}

.users-table th {
    background: var(--panel-dark);
    font-weight: 600;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--muted);
    white-space: nowrap;
}

.users-table th.sortable {
    cursor: pointer;
    user-select: none;
}

.users-table th.sortable:hover {
    color: var(--text);
    background: rgba(245, 140, 70, 0.1);
}

.users-table tbody tr:hover {
    background: rgba(245, 140, 70, 0.05);
}

.users-table tbody tr:last-child td {
    border-bottom: none;
}

/* Column widths */
.col-id {
    width: 60px;
}

.col-name {
    min-width: 140px;
    max-width: 200px;
}

.col-email {
    min-width: 180px;
    max-width: 240px;
}

.col-admin {
    width: 80px;
    text-align: center;
}

.col-deletion {
    width: 100px;
    white-space: nowrap;
}

.col-actions {
    width: 160px;
}

/* Truncated text with tooltip */
.truncated {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    cursor: help;
}

/* Admin badge */
.admin-badge {
    border: 1px solid rgba(245, 140, 70, 0.6);
    color: var(--text);
    padding: 2px 8px;
    font-size: 10px;
    letter-spacing: 2px;
}

.dash {
    color: var(--muted);
}

/* Action buttons */
.action-buttons {
    display: flex;
    gap: 8px;
}

.action-btn {
    border: 1px solid var(--border-strong);
    background: transparent;
    color: var(--text);
    padding: 6px 10px;
    font-family: inherit;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
    min-width: 70px;
    text-align: center;
}

.action-btn:hover {
    background: rgba(245, 140, 70, 0.1);
}

.action-btn.danger {
    border-color: rgba(243, 139, 139, 0.6);
    color: #f38b8b;
}

.action-btn.danger:hover {
    background: rgba(243, 139, 139, 0.1);
}

.action-btn.ghost {
    border-color: var(--border);
    color: var(--muted);
}

.action-btn.ghost:hover {
    border-color: var(--border-strong);
    color: var(--text);
    background: transparent;
}

.action-btn.primary {
    border-color: rgba(139, 243, 139, 0.6);
    color: #8bf38b;
}

.action-btn.primary:hover {
    background: rgba(139, 243, 139, 0.1);
}

/* Responsive */
@media (max-width: 768px) {
    .col-deletion,
    .col-id {
        display: none;
    }
}
</style>
