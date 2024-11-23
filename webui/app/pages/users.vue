<script lang="ts" setup>
import type {User} from '~/types'

const defaultColumns = [{
  key: 'id',
  label: '#'
}, {
  key: 'name',
  label: 'Name',
  sortable: true
}, {
  key: 'email',
  label: 'Email',
  sortable: true
}, {
  key: 'location',
  label: 'Location'
}, {
  key: 'status',
  label: 'Status'
}]

const q = ref('')
const selected = ref<User[]>([])
const selectedColumns = ref(defaultColumns)
const selectedStatuses = ref([])
const selectedLocations = ref([])
const sort = ref({column: 'id', direction: 'asc' as const})
const input = ref<{ input: HTMLInputElement }>()
const isNewUserModalOpen = ref(false)

const columns = computed(() => defaultColumns.filter(column => selectedColumns.value.includes(column)))

const query = computed(() => ({
  q: q.value,
  statuses: selectedStatuses.value,
  locations: selectedLocations.value,
  sort: sort.value.column,
  order: sort.value.direction
}))

const {data: users, pending} = await useFetch<User[]>('/api/users', {query, default: () => []})

const defaultLocations = users.value.reduce((acc, user) => {
  if (!acc.includes(user.location)) {
    acc.push(user.location)
  }
  return acc
}, [] as string[])

const defaultStatuses = users.value.reduce((acc, user) => {
  if (!acc.includes(user.status)) {
    acc.push(user.status)
  }
  return acc
}, [] as string[])

function onSelect(row: User) {
  const index = selected.value.findIndex(item => item.id === row.id)
  if (index === -1) {
    selected.value.push(row)
  } else {
    selected.value.splice(index, 1)
  }
}

defineShortcuts({
  '/': () => {
    input.value?.input?.focus()
  }
})
</script>

<template>
  <UDashboardPage>
    <UDashboardPanel grow>
      <UDashboardNavbar
        :badge="users.length"
        title="Users"
      >
        <template #right>
          <UInput
            ref="input"
            v-model="q"
            autocomplete="off"
            class="hidden lg:block"
            icon="i-heroicons-funnel"
            placeholder="Filter users..."
            @keydown.esc="$event.target.blur()"
          >
            <template #trailing>
              <UKbd value="/"/>
            </template>
          </UInput>

          <UButton
            color="gray"
            label="New user"
            trailing-icon="i-heroicons-plus"
            @click="isNewUserModalOpen = true"
          />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <template #left>
          <USelectMenu
            v-model="selectedStatuses"
            :options="defaultStatuses"
            :ui-menu="{ option: { base: 'capitalize' } }"
            icon="i-heroicons-check-circle"
            multiple
            placeholder="Status"
          />
          <USelectMenu
            v-model="selectedLocations"
            :options="defaultLocations"
            icon="i-heroicons-map-pin"
            multiple
            placeholder="Location"
          />
        </template>

        <template #right>
          <USelectMenu
            v-model="selectedColumns"
            :options="defaultColumns"
            class="hidden lg:block"
            icon="i-heroicons-adjustments-horizontal-solid"
            multiple
          >
            <template #label>
              Display
            </template>
          </USelectMenu>
        </template>
      </UDashboardToolbar>

      <UDashboardModal
        v-model="isNewUserModalOpen"
        :ui="{ width: 'sm:max-w-md' }"
        description="Add a new user to your database"
        title="New user"
      >
        <!-- ~/components/users/UsersForm.vue -->
        <UsersForm @close="isNewUserModalOpen = false"/>
      </UDashboardModal>

      <UTable
        v-model="selected"
        v-model:sort="sort"
        :columns="columns"
        :loading="pending"
        :rows="users"
        :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
        class="w-full"
        sort-mode="manual"
        @select="onSelect"
      >
        <template #name-data="{ row }">
          <div class="flex items-center gap-3">
            <UAvatar
              :alt="row.name"
              size="xs"
              v-bind="row.avatar"
            />

            <span class="text-gray-900 dark:text-white font-medium">{{ row.name }}</span>
          </div>
        </template>

        <template #status-data="{ row }">
          <UBadge
            :color="row.status === 'subscribed' ? 'green' : row.status === 'bounced' ? 'orange' : 'red'"
            :label="row.status"
            class="capitalize"
            variant="subtle"
          />
        </template>
      </UTable>
    </UDashboardPanel>
  </UDashboardPage>
</template>
